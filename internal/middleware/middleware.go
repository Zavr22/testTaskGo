package middleware

import (
	"github.com/Zavr22/testTaskGo/internal/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func AdminMiddleware() echo.MiddlewareFunc {
	return middleware.BasicAuth(func(username string, password string, context echo.Context) (bool, error) {
		admin, err := utils.IsAdmin(username, password)
		if admin {
			return true, nil
		}
		return false, err
	})
}
