package service

import (
	"backend-mobAppRest/internal/model"
	"backend-mobAppRest/internal/repository"
	"backend-mobAppRest/internal/repository/customRepositoryError"
	"backend-mobAppRest/internal/service/customServiceError"
	"backend-mobAppRest/pkg/tg"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"strconv"
	"time"
)

type authService struct {
	SignKey       []byte
	AuthSignature string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration

	AuthRepository  repository.AuthRepository
	TokenRepository repository.TokenRepository
	UserRepository  repository.UserRepository
}

func (a *authService) SignUp(username, email, password string) (*model.Tokens, error) {
	signingPassword, err := signedPassword(password, a.AuthSignature)
	if err != nil {
		go tg.SendError(err.Error(), "/api/auth/signup")
		slog.Debug("error signing password", err, password, a.AuthSignature)
		return nil, customServiceError.ErrUnknown
	}

	hashedPassword, err := hashPassword(signingPassword)
	if err != nil {
		go tg.SendError(err.Error(), "/api/auth/signup")
		slog.Debug("error hashing password", err, signingPassword)
		return nil, customServiceError.ErrUnknown
	}

	newUser := &model.User{
		Name:     username,
		Email:    email,
		Password: hashedPassword,
	}

	id, err := a.AuthRepository.CreateUser(newUser)
	if err != nil {
		go tg.SendError(err.Error(), "/api/auth/signup")
		slog.Debug("error create user", err, newUser)
		return nil, customServiceError.ErrUnknown
	}

	userIDStr := strconv.Itoa(id)

	tokensModel, err := generateTokens(userIDStr, a.AccessExpiry, a.RefreshExpiry, a.SignKey)
	if err != nil {
		go tg.SendError(err.Error(), "/api/auth/signup")
		slog.Debug("error generate tokens", err, userIDStr, a.AccessExpiry, a.RefreshExpiry, a.SignKey)
		return nil, customServiceError.ErrUnknown
	}

	err = a.TokenRepository.CreateTokens(tokensModel.AccessToken, tokensModel.RefreshToken)
	if err != nil {
		go tg.SendError(err.Error(), "/api/auth/signup")
		slog.Debug("error create tokens in database", err, tokensModel.AccessToken, tokensModel.RefreshToken)
		return nil, customServiceError.ErrUnknown
	}

	return &model.Tokens{
		AccessToken:  tokensModel.AccessToken.Token,
		RefreshToken: tokensModel.RefreshToken.Token,
	}, nil
}

func (a *authService) SignIn(email, password string) (*model.Tokens, error) {
	signingPassword, err := signedPassword(password, a.AuthSignature)
	if err != nil {
		go tg.SendError(err.Error(), "/api/auth/signin")
		slog.Debug("error signing password", err, password, a.AuthSignature)
		return nil, customServiceError.ErrUnknown
	}

	user, err := a.UserRepository.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, customRepositoryError.ErrInvalidEmailOrPassword) {
			return nil, customServiceError.ErrInvalidEmailOrPassword
		}
		go tg.SendError(err.Error(), "/api/auth/signin")
		slog.Debug("error get user by email in database", err, email)
		return nil, customServiceError.ErrUnknown
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(signingPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, customServiceError.ErrInvalidEmailOrPassword
		}
		go tg.SendError(err.Error(), "/api/auth/signin")
		slog.Debug("error compare hash and password", err, password, a.AuthSignature)
		return nil, customServiceError.ErrUnknown
	}

	userID := user.ID
	userIDStr := strconv.Itoa(userID)

	tokensModel, err := generateTokens(userIDStr, a.AccessExpiry, a.RefreshExpiry, a.SignKey)
	if err != nil {
		go tg.SendError(err.Error(), "/api/auth/signin")
		slog.Debug("generate tokens", err, signingPassword)
		return nil, customServiceError.ErrUnknown
	}

	err = a.TokenRepository.CreateTokens(tokensModel.AccessToken, tokensModel.RefreshToken)
	if err != nil {
		go tg.SendError(err.Error(), "/api/auth/signin")
		slog.Debug("error create tokens in database", err, tokensModel.AccessToken, tokensModel.RefreshToken)
		return nil, customServiceError.ErrUnknown
	}

	return &model.Tokens{
		AccessToken:  tokensModel.AccessToken.Token,
		RefreshToken: tokensModel.RefreshToken.Token,
	}, nil

}

func newAuthService(SignKey []byte, accessExpiry time.Duration, refreshExpiry time.Duration, authSignature string,
	authRepository repository.AuthRepository, tokenRepository repository.TokenRepository, userRepository repository.UserRepository) *authService {
	return &authService{
		SignKey:         SignKey,
		AccessExpiry:    accessExpiry,
		RefreshExpiry:   refreshExpiry,
		AuthSignature:   authSignature,
		AuthRepository:  authRepository,
		TokenRepository: tokenRepository,
		UserRepository:  userRepository,
	}
}

func generateTokens(userID string, accessExpiry, refreshExpiry time.Duration, signKey []byte) (*model.TokensModel, error) {
	accessTokenExp := time.Now().Add(accessExpiry)
	accessTokenIss := time.Now()

	refreshTokenExp := time.Now().Add(refreshExpiry)
	refreshTokenIss := time.Now()

	accessClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(accessTokenExp),
		IssuedAt:  jwt.NewNumericDate(accessTokenIss),
	})

	accessToken, err := accessClaims.SignedString(signKey)
	if err != nil {
		slog.Debug("error signing access token", err, signKey)
		return nil, err
	}

	refreshClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(refreshTokenExp),
		IssuedAt:  jwt.NewNumericDate(refreshTokenIss),
	})

	refreshToken, err := refreshClaims.SignedString(signKey)
	if err != nil {
		slog.Debug("error signing refresh token", err, signKey)
		return nil, err
	}

	modelAccessToken := &model.AccessToken{
		Issuer:       "null",
		Audience:     "null",
		Subject:      userID,
		ExpirationAt: accessTokenExp,
		IssuedAt:     accessTokenIss,
		Token:        accessToken,
	}

	modelRefreshToken := &model.RefreshToken{
		Issuer:       "null",
		Audience:     "null",
		Subject:      userID,
		ExpirationAt: refreshTokenExp,
		IssuedAt:     refreshTokenIss,
		Token:        refreshToken,
	}

	return &model.TokensModel{
		AccessToken:  modelAccessToken,
		RefreshToken: modelRefreshToken,
	}, nil
}

func signedPassword(password string, signature string) (string, error) {
	h := hmac.New(sha256.New, []byte(signature))
	_, err := h.Write([]byte(password))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
