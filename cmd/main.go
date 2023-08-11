package main

import (
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"os"
	"testTask/cmd/internal/handler"
	"testTask/cmd/internal/repository"
	"testTask/cmd/internal/service"
)

func main() {
	e := echo.New()
	logger := logrus.New()
	logger.Out = os.Stdout
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			logrus.WithFields(logrus.Fields{
				"URI":    values.URI,
				"status": values.Status,
			}).Info("request")
			return nil
		},
	}))
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	userRepo := repository.NewUserRepo(rdb)
	authRepo := repository.NewAuthRepo(rdb)

	userServ := service.NewUserService(userRepo)
	authServ := service.NewAuthService(authRepo)

	profileHandler := handler.NewHandler(userServ, authServ)
	profileHandler.InitRoutes(e)
}
