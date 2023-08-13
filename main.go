package main

import (
	_ "github.com/Zavr22/testTaskGo/docs"
	"github.com/Zavr22/testTaskGo/internal/handler"
	middleware2 "github.com/Zavr22/testTaskGo/internal/middleware"
	repository2 "github.com/Zavr22/testTaskGo/internal/repository"
	service2 "github.com/Zavr22/testTaskGo/internal/service"
	"github.com/Zavr22/testTaskGo/internal/utils"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"
	"os"
)

// @title TestTask Server
// @version 1.0
// @description API Server for Test Task

// @host localhost:9000
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	e := echo.New()
	e.GET("/swagger/*", echoSwagger.WrapHandler)
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
	e.Use(middleware2.BasicAuthMiddleware())
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	utils.SetRedisClient(rdb)
	userRepo := repository2.NewUserRepo(rdb)
	authRepo := repository2.NewAuthRepo(rdb)

	userServ := service2.NewUserService(userRepo)
	authServ := service2.NewAuthService(authRepo)

	profileHandler := handler.NewHandler(userServ, authServ)
	profileHandler.InitRoutes(e)
}
