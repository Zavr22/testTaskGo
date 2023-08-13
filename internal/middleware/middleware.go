package middleware

import (
	"encoding/base64"
	"github.com/Zavr22/testTaskGo/internal/utils"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

func BasicAuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if c.Path() == "/auth/sign_up" || c.Path() == "/auth/sign_in" || c.Path() == "/swagger/*" {
				return next(c)
			}
			req := c.Request()
			headers := req.Header
			atHeader := headers.Get("Authorization")
			if atHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "no auth header")
			}
			credentials := strings.Split(atHeader, ":")
			if len(credentials) != 2 {
				return echo.NewHTTPError(http.StatusBadRequest, "invalid auth header")
			}
			decodedUserID, err := base64.StdEncoding.DecodeString(credentials[0])
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "error decoding userID")
			}

			decodedPassword, err := base64.StdEncoding.DecodeString(credentials[1])
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "error decoding password")
			}
			decodedCredentials := []string{string(decodedUserID), string(decodedPassword)}

			valid, err := utils.IsUserValid(decodedCredentials)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
			}
			if !valid {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid username or password")
			}
			return next(c)
		}
	}
}

func AdminMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			req := c.Request()
			headers := req.Header
			atHeader := headers.Get("Authorization")
			credentials := strings.Split(atHeader, ":")
			decodedUserID, err := base64.StdEncoding.DecodeString(credentials[0])
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "error decoding userID")
			}
			decodedPassword, err := base64.StdEncoding.DecodeString(credentials[1])
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "error decoding password")
			}
			decodedCredentials := []string{string(decodedUserID), string(decodedPassword)}
			admin, err := utils.IsAdmin(decodedCredentials)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
			}
			if !admin {
				return echo.NewHTTPError(http.StatusForbidden, "you are not an admin")
			}
			return next(c)
		}
	}
}
