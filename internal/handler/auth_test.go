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
	"net/http"
	"net/http/httptest"
	"testing"
)

var id = uuid.New()

func TestHandler_SignUp(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockAuthorization, user *models.SignUpInput)

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
			mockBehavior: func(r *mock_service.MockAuthorization, user *models.SignUpInput) {
				ctx := context.Background()
				r.EXPECT().SignUp(ctx, user).Return(id, nil)
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

			repo := mock_service.NewMockAuthorization(c)
			test.mockBehavior(repo, &test.inputUser)

			authService := service.NewAuthService(repo)
			handler := Handler{authS: authService}

			// Init Endpoint
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/sign_up", bytes.NewBufferString(test.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			// Perform Request
			err := handler.SignUp(ctx)

			// Assertions
			assert.NoError(t, err)
			assert.JSONEq(t, test.expectedResponseBody, rec.Body.String())
		})
	}
}

func TestHandler_SignIn(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockAuthorization, user *models.SignInInput)

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            models.SignInInput
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"username": "username", "password": "qwerty"}`,
			inputUser: models.SignInInput{
				Username: "username",
				Password: "qwerty",
			},
			mockBehavior: func(r *mock_service.MockAuthorization, user *models.SignInInput) {
				ctx := context.Background()
				r.EXPECT().SignIn(ctx, user).Return(nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"message":"you successfully signed in"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockAuthorization(c)
			test.mockBehavior(repo, &test.inputUser)

			authService := service.NewAuthService(repo)
			handler := Handler{authS: authService}

			// Init Endpoint
			// Init Endpoint
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/sign_up", bytes.NewBufferString(test.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			// Perform Request
			err := handler.SignIn(ctx)

			// Assertions
			assert.NoError(t, err)
			assert.Equal(t, test.expectedStatusCode, rec.Code)
			assert.JSONEq(t, test.expectedResponseBody, rec.Body.String())
		})
	}
}
