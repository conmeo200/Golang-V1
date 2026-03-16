package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/conmeo200/Golang-V1/internal/dto"
	"github.com/conmeo200/Golang-V1/internal/model"
	"github.com/conmeo200/Golang-V1/internal/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserService is a mock type for the UserService
type MockUserService struct {
	mock.Mock
}

// Mock the UserService methods that are used in the handler
func (m *MockUserService) CreateUser(req dto.CreateUserRequest) (*model.User, error) {
	args := m.Called(req)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) GetUserByID(id string) (*model.User, error) {
	args := m.Called(id)
	return args.Get(0).(*model.User), args.Error(1)
}

// Add other mocked methods as you expand tests

func TestCreateUser(t *testing.T) {
	// 1. Setup
	mockUserService := new(MockUserService)
	userHandler := NewUserHandler(mockUserService)

	// 2. Define Input and Expected Output
	createReq := dto.CreateUserRequest{
		Email:    "test@example.com",
		Password: "password123",
		Balance:  100.0,
	}

	expectedUser := &model.User{
		ID:       uuid.New(),
		Email:    "test@example.com",
		Balance:  100.0,
		Role:     "user",
		Status:   "active",
	}

	// 3. Setup Mock
	// We need to adjust the argument to match what the handler will pass to the service
	// For now, let's assume the handler converts its internal request struct to dto.CreateUserRequest
	mockUserService.On("CreateUser", mock.AnythingOfType("dto.CreateUserRequest")).Return(expectedUser, nil)

	// 4. Create HTTP Request
	// The handler's internal struct expects `password`, not `PasswordHash`
	payload := map[string]interface{}{
        "email":    createReq.Email,
        "password": createReq.Password, // This matches the `json:"password"` tag in handler
        "balance":  createReq.Balance,
    }
	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	// 5. Call Handler
	// Note: The actual endpoint path is defined in the router, not here.
	// We are unit-testing the handler function directly.
	userHandler.CreateUser(rr, req)

	// 6. Assert Results
	assert.Equal(t, http.StatusCreated, rr.Code)

	var returnedUser model.User
    err = json.Unmarshal(rr.Body.Bytes(), &returnedUser)
    assert.NoError(t, err)

	assert.Equal(t, expectedUser.ID, returnedUser.ID)
    assert.Equal(t, expectedUser.Email, returnedUser.Email)
    assert.Equal(t, expectedUser.Balance, returnedUser.Balance)

	// 7. Verify that the mock was called
	mockUserService.AssertExpectations(t)
}
