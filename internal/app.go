package internal

import (
	"backend-mobAppRest/config"
	"backend-mobAppRest/internal/controller"
	"backend-mobAppRest/internal/model"
	"backend-mobAppRest/internal/repository"
	"backend-mobAppRest/internal/service"
	"backend-mobAppRest/pkg/database"
	"backend-mobAppRest/pkg/server"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Config) {
	connection, err := database.NewConnection(cfg.Database.URI, 0, 5, &model.User{}, &model.Tokens{},
		&model.Product{}, &model.Cart{}, &model.AccessToken{}, &model.RefreshToken{}, &model.Category{}, &model.UserAddress{},
		&model.City{})
	if err != nil {
		panic(err)
	}

	repositories := repository.NewRepositories(connection)

	deps := &service.DependenciesService{
		SignKey:           []byte(cfg.JWT.SignKey),
		AuthSignature:     cfg.Auth.AuthSignature,
		AccessExpiry:      cfg.JWT.AccessExpiry,
		RefreshExpiry:     cfg.JWT.RefreshExpiry,
		AuthRepository:    repositories.AuthRepository,
		TokenRepository:   repositories.TokenRepository,
		UserRepository:    repositories.UserRepository,
		CatalogRepository: repositories.CatalogRepository,
		CartRepository:    repositories.CartRepository,
	}

	services := service.NewService(deps)

	route := gin.New()
	controller.RouteAPI(route, services)

	srv := server.NewServer(route, server.Port(cfg.Http.Port), server.ReadTimeout(cfg.Http.ReadTimeout),
		server.WriteTimeout(cfg.Http.WriteTimeout), server.ShutdownTimeout(cfg.Http.ShutdownTimeout))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Println("Run: " + s.String())
	case err := <-srv.Notify():
		log.Println(err, "Run: signal.Notify")
	}

	err = srv.Shutdown()
	if err != nil {
		log.Println(err, "Run: server shutdown")
	}
}
