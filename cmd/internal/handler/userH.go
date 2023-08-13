package handler

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"testTask/cmd/models"
)

func (h *Handler) CreateUser(c echo.Context) error {
	var reqBody models.UserProfile
	errBind := c.Bind(&reqBody)
	if errBind != nil {
		logrus.WithFields(logrus.Fields{
			"reqBody": reqBody,
		}).Errorf("Bind json, %s", errBind)
		return echo.NewHTTPError(http.StatusInternalServerError, models.CommonResponse{Message: "data not correct"})
	}
	userID, err := h.userS.CreateUser(c.Request().Context(), reqBody.Email, reqBody.Username, reqBody.Password, reqBody.Admin)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"userID": userID,
		}).Errorf("error while creating user, %s", err)
		return echo.NewHTTPError(http.StatusBadRequest, models.CommonResponse{Message: "error while creating user"})
	}
	return c.JSON(http.StatusOK, userID)
}

func (h *Handler) GetUsers(c echo.Context) error {
	users, err := h.userS.GetAllUsers(c.Request().Context())
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"users": users,
		}).Errorf("Get users failed, %s", err)
		return echo.NewHTTPError(http.StatusBadRequest, "get users failed")
	}
	return c.JSON(http.StatusOK, users)
}

func (h *Handler) GetUserByID(c echo.Context) error {
	c.Param("id")
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"userID": userID,
		}).Errorf("wrong id format, %s", err)
		return echo.NewHTTPError(http.StatusBadRequest, "wrong id format")
	}
	user, err := h.userS.GetUser(c.Request().Context(), userID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"user": user,
		}).Errorf("cannot get user, %s", err)
		return echo.NewHTTPError(http.StatusBadRequest, "cannot get user")
	}
	return c.JSON(http.StatusOK, user)
}

func (h *Handler) UpdateUser(c echo.Context) error {
	var reqBody models.UpdateProfileInput
	errBind := c.Bind(&reqBody)
	if errBind != nil {
		logrus.WithFields(logrus.Fields{
			"reqBody": reqBody,
		}).Errorf("Bind json, %s", errBind)
		return echo.NewHTTPError(http.StatusInternalServerError, models.CommonResponse{Message: "data not correct"})
	}
	c.Param("id")
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"userID": userID,
		}).Errorf("wrong id format, %s", err)
		return echo.NewHTTPError(http.StatusBadRequest, "wrong id format")
	}
	errUpdate := h.userS.UpdateProfile(c.Request().Context(), userID, reqBody)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"userID": userID,
		}).Errorf("error while creating user, %s", errUpdate)
		return echo.NewHTTPError(http.StatusBadRequest, models.CommonResponse{Message: "error while creating user"})
	}
	return c.JSON(http.StatusOK, models.CommonResponse{Message: "user updated successfully"})
}

func (h *Handler) DeleteUser(c echo.Context) error {
	c.Param("id")
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"userID": userID,
		}).Errorf("wrong id format, %s", err)
		return echo.NewHTTPError(http.StatusBadRequest, "wrong id format")
	}
	errDelete := h.userS.DeleteProfile(c.Request().Context(), userID)
	if err != nil {
		logrus.Errorf("error while delete user, %s", errDelete)
	}
	return c.JSON(http.StatusOK, models.CommonResponse{Message: "user deleted successfully"})
}
