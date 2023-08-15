package handler

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Zavr22/testTaskGo/internal/models"
	"github.com/Zavr22/testTaskGo/internal/service"
	mock_service "github.com/Zavr22/testTaskGo/internal/service/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Тест для хендлера CreateUser
func TestHandler_CreateUser(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockUser, user *models.SignUpInput)

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            models.SignUpInput
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "ok",
			inputBody: `{"email": "testemail@g.c", "username": "username",  "password": "qwerty", "admin" : true}`,
			inputUser: models.SignUpInput{
				Username: "username",
				Email:    "testemail@g.c",
				Password: "qwerty",
				Admin:    true,
			},
			mockBehavior: func(r *mock_service.MockUser, user *models.SignUpInput) {
				ctx := context.Background()
				r.EXPECT().CreateUser(ctx, user).Return(id, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: fmt.Sprintf(`{"user_id":"%s"}`, id),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockUser(c)
			test.mockBehavior(repo, &test.inputUser)

			userService := service.NewUserService(repo)
			handler := Handler{userS: userService}

			// Init Endpoint
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/create_user", bytes.NewBufferString(test.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			// Perform Request
			err := handler.CreateUser(ctx)

			// Assertions
			assert.NoError(t, err)
			assert.JSONEq(t, test.expectedResponseBody, rec.Body.String())
		})
	}
}

// Тест для хендлера GetUsers
func TestHandler_GetUsers(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockUser)
	tests := []struct {
		name                 string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "ok",
			mockBehavior: func(r *mock_service.MockUser) {
				ctx := context.Background()
				r.EXPECT().GetAllUsers(ctx).Return([]*models.UserResponse{}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: "[]",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			mockUserService := mock_service.NewMockUser(c)
			test.mockBehavior(mockUserService)

			userService := service.NewUserService(mockUserService)
			handler := Handler{userS: userService}

			// Init Endpoint
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			// Perform Request
			err := handler.GetUsers(ctx)

			// Assertions
			require.NoError(t, err)
			require.Equal(t, test.expectedStatusCode, rec.Code)
			require.Equal(t, strings.TrimSpace(test.expectedResponseBody), strings.TrimSpace(rec.Body.String()))
		})
	}
}

// Тест для хендлера GetUserByID
func TestHandler_GetUserByID(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockUser, id uuid.UUID)

	tests := []struct {
		name                 string
		requestID            string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "ok",
			mockBehavior: func(r *mock_service.MockUser, id uuid.UUID) {
				ctx := context.Background()
				r.EXPECT().GetUser(ctx, id).Return(models.UserResponse{}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"ID":"00000000-0000-0000-0000-000000000000","user_name":"","email":""}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			mockUserService := mock_service.NewMockUser(c)
			test.mockBehavior(mockUserService, id)

			handler := &Handler{userS: mockUserService}

			// Init Endpoint
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/users/%s", id), nil)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.SetParamNames("id")
			ctx.SetParamValues(id.String())

			// Perform Request
			err := handler.GetUserByID(ctx)

			// Assertions
			require.NoError(t, err)
			require.Equal(t, test.expectedStatusCode, rec.Code)
			require.Equal(t, strings.TrimSpace(test.expectedResponseBody), strings.TrimSpace(rec.Body.String()))
		})
	}
}

func TestHandler_UpdateUser(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockUser, id uuid.UUID, input models.UpdateProfileInput)

	tests := []struct {
		name                 string
		requestID            string
		requestBody          string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "ok",
			requestBody: `{"email": "updatedemail@g.c", "username": "updatedusername", "password": "updatedpassword", "admin": false}`,
			mockBehavior: func(r *mock_service.MockUser, id uuid.UUID, input models.UpdateProfileInput) {
				ctx := context.Background()
				r.EXPECT().UpdateProfile(ctx, id, input).Return(nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"message":"user updated successfully"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			mockUserService := mock_service.NewMockUser(c)
			test.mockBehavior(mockUserService, id, models.UpdateProfileInput{
				NewEmail:    "updatedemail@g.c",
				NewUsername: "updatedusername",
				NewPassword: "updatedpassword",
				Admin:       false,
			})

			handler := &Handler{userS: mockUserService}

			// Init Endpoint
			e := echo.New()
			req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/users/%s", id), strings.NewReader(test.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.SetParamNames("id")
			ctx.SetParamValues(id.String())

			// Perform Request
			err := handler.UpdateUser(ctx)

			// Assertions
			require.NoError(t, err)
			require.Equal(t, test.expectedStatusCode, rec.Code)
			require.Equal(t, strings.TrimSpace(test.expectedResponseBody), strings.TrimSpace(rec.Body.String()))
		})
	}
}

func TestHandler_DeleteUser(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockUser, id uuid.UUID)

	tests := []struct {
		name                 string
		requestID            string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "ok",
			mockBehavior: func(r *mock_service.MockUser, id uuid.UUID) {
				ctx := context.Background()
				r.EXPECT().DeleteProfile(ctx, id).Return(nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"message":"user deleted successfully"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			mockUserService := mock_service.NewMockUser(c)
			test.mockBehavior(mockUserService, id)

			handler := &Handler{userS: mockUserService}

			// Init Endpoint
			e := echo.New()
			req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/users/%s", id.String()), nil)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.SetParamNames("id")
			ctx.SetParamValues(id.String())

			// Perform Request
			err := handler.DeleteUser(ctx)

			// Assertions
			require.NoError(t, err)
			require.Equal(t, test.expectedStatusCode, rec.Code)
			require.Equal(t, strings.TrimSpace(test.expectedResponseBody), strings.TrimSpace(rec.Body.String()))
		})
	}
}
