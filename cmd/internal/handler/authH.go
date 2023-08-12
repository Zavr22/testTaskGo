package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"testTask/cmd/models"
)

func (h *Handler) SignUp(c echo.Context) error {
	user := models.SignUpInput{}
	errBind := c.Bind(&user)
	if errBind != nil {
		logrus.WithFields(logrus.Fields{
			"user": user,
		}).Errorf("Bind json, %s", errBind)
		return echo.NewHTTPError(http.StatusInternalServerError, models.CommonResponse{Message: "data not correct"})
	}
	userID, err := h.authS.SignUp(c.Request().Context(), &user)
	if err != nil {
		logrus.Errorf("CREATE USER request, %s", err)
		return echo.NewHTTPError(http.StatusBadRequest, models.CommonResponse{Message: "cannot create user"})
	}
	return c.JSON(http.StatusOK, userID)
}

func (h *Handler) SignIn(c echo.Context) error {
	user := models.SignInInput{}
	errBind := c.Bind(&user)
	if errBind != nil {
		logrus.WithFields(logrus.Fields{
			"user": user,
		}).Errorf("Bind json, %s", errBind)
		return echo.NewHTTPError(http.StatusInternalServerError, models.CommonResponse{Message: "data not correct"})
	}
	token, err := h.authS.SignIn(c.Request().Context(), &user)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.CommonResponse{Message: "fail to sign in"})
	}

	return c.JSON(http.StatusOK, token)
}
