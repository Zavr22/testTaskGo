package handler

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"testTask/cmd/models"
)

func (h *Handler) CreateUser(c echo.Context) error {
	return nil
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
	return nil
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
