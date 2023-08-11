package middleware

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func BasicAuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if c.Path() == "/auth/sign_up" || c.Path() == "/auth/sign_in" {
				return next(c)
			}
			req := c.Request()
			headers := req.Header
			atHeader := headers.Get("Authorization")
			if atHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "no auth header")
			}

			return next(c)
		}
	}
}
