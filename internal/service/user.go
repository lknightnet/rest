package service

import (
	"backend-mobAppRest/internal/model"
	"backend-mobAppRest/internal/repository"
	"backend-mobAppRest/internal/repository/customRepositoryError"
	"backend-mobAppRest/internal/service/customServiceError"
	"backend-mobAppRest/pkg/tg"
	"errors"
	"log/slog"
)

type userService struct {
	UserRepository repository.UserRepository
}

func (u *userService) ChangeInformation(token string, user *model.User) error {
	err := u.UserRepository.ChangeInformation(token, user)
	if err != nil {
		if errors.Is(err, customRepositoryError.ErrUserNotFound) {
			return customServiceError.ErrUnknown
		}
		go tg.SendError(err.Error(), "/api/user/change")
		slog.Debug("error get user by access token in database", err, token)
		return customServiceError.ErrUnknown
	}
	return nil
}

func (u *userService) GetUserByAccessToken(token string) (*model.ViewUser, error) {
	user, err := u.UserRepository.GetUserByToken(token)
	if err != nil {
		if errors.Is(err, customRepositoryError.ErrUserNotFound) {
			return nil, customServiceError.ErrUnknown
		}
		go tg.SendError(err.Error(), "/api/user/get")
		slog.Debug("error get user by access token in database", err, token)
		return nil, customServiceError.ErrUnknown
	}

	viewUser := &model.ViewUser{
		Name:    user.Name,
		Phone:   user.Phone,
		Bonuses: user.Bonuses,
		Address: user.Address,
	}

	return viewUser, nil
}

func newUserService(userRepository repository.UserRepository) *userService {
	return &userService{UserRepository: userRepository}
}
