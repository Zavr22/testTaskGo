package handler

import (
	"bytes"
	"encoding/json"
	"github.com/Zavr22/testTaskGo/cmd/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestSignUpHandler(t *testing.T) {
	mockUser := new(MockUser)
	mockAuth := new(MockAuthorization)
	handler := NewHandler(mockUser, mockAuth)
	requestBody := models.SignInInput{
		Username: "testuser",
		Password: "testpassword",
	}
	jsonBytes, err := json.Marshal(requestBody)
	reader := bytes.NewReader(jsonBytes)
	assert.NoError(t, err)
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/signup", reader)
	recorder := httptest.NewRecorder()
	c := e.NewContext(req, recorder)
	expectedUserID := uuid.New()
	mockUser.On("SignUp", mock.Anything, "testuser", "testpassword").Return(expectedUserID, nil)
	err := handler.SignUp(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, recorder.Code)
	mockUser.AssertCalled(t, "SignUp", mock.Anything, "testuser", "testpassword")
}

func TestSignInHandler(t *testing.T) {
	mockUser := new(MockUser)
	mockAuth := new(MockAuthorization)
	handler := NewHandler(mockUser, mockAuth)
	requestBody := models.SignInInput{
		Username: "testuser",
		Password: "testpassword",
	}
	jsonBytes, err := json.Marshal(requestBody)
	reader := bytes.NewReader(jsonBytes)
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/signin", reader)
	recorder := httptest.NewRecorder()
	c := e.NewContext(req, recorder)
	expectedToken := "testtoken"
	mockAuth.On("SignIn", mock.Anything, "testuser", "testpassword").Return(expectedToken, nil)
	err := handler.SignIn(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, recorder.Code)
	mockAuth.AssertCalled(t, "SignIn", mock.Anything, "testuser", "testpassword")
}
