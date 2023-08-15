package handler

import (
	"github.com/Zavr22/testTaskGo/internal/models"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
)

// SignUp used to sign up using serv userS
//
// @Summary SignUp
// @Tags auth
// @Description create account
// @ID create-account
// @Accept  json
// @Produce  json
// @Param input body models.SignUpInput true "account info"
// @Success 200 {object} models.CreateUserResponse
// @Failure 400 {object} models.CommonResponse "cannot create user"
// @Failure 500 {object} models.CommonResponse "data not correct"
// @Router /auth/sign_up [post]
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
	return c.JSON(http.StatusOK, models.CreateUserResponse{UserID: userID})
}

// SignIn used to sign in using serv userS
//
// @Summary SignIn
// @Tags auth
// @Description login
// @ID login
// @Accept  json
// @Produce  json
// @Param input body models.SignInInput true "enter username and password"
// @Success 200 {object} models.CommonResponse
// @Failure 401 {object} models.CommonResponse
// @Failure 500 {object} models.CommonResponse
// @Router /auth/sign_in [post]
func (h *Handler) SignIn(c echo.Context) error {
	user := models.SignInInput{}
	errBind := c.Bind(&user)
	if errBind != nil {
		logrus.WithFields(logrus.Fields{
			"user": user,
		}).Errorf("Bind json, %s", errBind)
		return echo.NewHTTPError(http.StatusInternalServerError, models.CommonResponse{Message: "data not correct"})
	}
	err := h.authS.SignIn(c.Request().Context(), &user)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.CommonResponse{Message: "fail to sign in"})
	}

	return c.JSON(http.StatusOK, models.CommonResponse{Message: "you successfully signed in"})
}
