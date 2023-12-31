package main

import (
	"context"
	_ "github.com/Zavr22/testTaskGo/docs"
	"github.com/Zavr22/testTaskGo/internal/handler"
	repository2 "github.com/Zavr22/testTaskGo/internal/repository"
	service2 "github.com/Zavr22/testTaskGo/internal/service"
	"github.com/Zavr22/testTaskGo/internal/utils"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"
	"os"
	"os/signal"
	"syscall"
)

// @title TestTask Server
// @version 1.0
// @description API Server for Test Task

// @host localhost:9000
// @BasePath /

// @securityDefinitions.basic BasicAuth
// @description Basic authentication username and password
// @type basic
func main() {
	e := echo.New()

	// SWAGGER
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// LOGGERS
	logger := logrus.New()
	logger.Out = os.Stdout

	// MIDDLEWARE
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

	// REDIS
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	// SKIPPED HANDLERS BY MIDDLEWARE
	skipperFunc := func(c echo.Context) bool {
		if c.Path() == "/auth/sign_up" {
			return true
		}
		if c.Path() == "/auth/sign_in" {
			return true
		}
		if c.Path() == "/swagger/*" {
			return true
		}
		return false
	}

	// MIDDLEWARE
	config := middleware.BasicAuthConfig{
		Skipper: skipperFunc,
		Validator: func(username, password string, c echo.Context) (bool, error) {
			val, err := utils.IsUserValid(username, password)
			if err != nil {
				logger.Println(err)
				return false, err
			}
			if val == true {
				return true, nil
			} else {
				return false, nil
			}
		},
	}
	e.Use(middleware.BasicAuthWithConfig(config))

	// init utils, services, repos
	utils.SetRedisClient(rdb)
	userRepo := repository2.NewUserRepo(rdb)
	authRepo := repository2.NewAuthRepo(rdb)

	userServ := service2.NewUserService(userRepo)
	authServ := service2.NewAuthService(authRepo)

	profileHandler := handler.NewHandler(userServ, authServ)
	profileHandler.InitRoutes(e)

	// Graceful shutdown
	logrus.Print("App Started")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("App Shutting Down")

	if err := e.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}

	if err := rdb.Close(); err != nil {
		logrus.Errorf("error occured on db connection close: %s", err.Error())
	}
}
