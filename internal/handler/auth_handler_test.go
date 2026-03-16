package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/conmeo200/Golang-V1/internal/dto"
	"github.com/conmeo200/Golang-V1/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuthService is a mock type for the AuthService
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(user dto.RegisterRequest) (*dto.UserResponse, error) {
	args := m.Called(user)
	return args.Get(0).(*dto.UserResponse), args.Error(1)
}

func (m *MockAuthService) Login(credentials dto.LoginRequest) (string, string, error) {
	args := m.Called(credentials)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *MockAuthService) ForgotPassword(email string) error {
	args := m.Called(email)
	return args.Error(0)
}

func (m *MockAuthService) RefreshToken(refreshToken string) (string, error) {
	args := m.Called(refreshToken)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) Logout(token string) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockAuthService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	args := m.Called(userID, oldPassword, newPassword)
	return args.Error(0)
}

func (m *MockAuthService) RevokeToken(token string) error {
	args := m.Called(token)
	return args.Error(0)
}

func TestRegisterHandler(t *testing.T) {
	mockAuthService := new(MockAuthService)
	authHandler := NewAuthHandler(mockAuthService)

	// Setup the expected call and return value
	registerReq := dto.RegisterRequest{Username: "testuser", Password: "password", Email: "test@example.com"}
	expectedUser := &dto.UserResponse{ID: 1, Username: "testuser", Email: "test@example.com"}
	mockAuthService.On("Register", registerReq).Return(expectedUser, nil)

	// Create a request body
	body, _ := json.Marshal(registerReq)
	req, err := http.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(body))
	assert.NoError(t, err)

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	authHandler.RegisterHandler(rr, req)

	// Assert the status code and response body
	assert.Equal(t, http.StatusCreated, rr.Code)

	var response dto.APIResponse
	json.Unmarshal(rr.Body.Bytes(), &response)

	// You might need to adjust the assertions based on the actual structure of your APIResponse
	assert.True(t, response.Success)
	assert.Equal(t, "User registered successfully", response.Message)

	// Assert that the mock was called as expected
	mockAuthService.AssertExpectations(t)
}
