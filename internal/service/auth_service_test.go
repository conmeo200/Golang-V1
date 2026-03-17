package service

import (
	"context"
	"testing"

	"github.com/conmeo200/Golang-V1/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// MockUserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetUser(ctx context.Context, id string) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) ListUser(ctx context.Context) ([]model.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.User), args.Error(1)
}

func (m *MockUserRepository) UpdateBalance(ctx context.Context, id uint, newBalance float64) error {
	args := m.Called(ctx, id, newBalance)
	return args.Error(0)
}

func (m *MockUserRepository) UpdatePassword(ctx context.Context, id string, newHash string) error {
	args := m.Called(ctx, id, newHash)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockTokenRepository
type MockTokenRepository struct {
	mock.Mock
}

func (m *MockTokenRepository) BlacklistToken(ctx context.Context, token *model.TokenBlacklist) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockTokenRepository) IsBlacklisted(ctx context.Context, tokenString string) bool {
	args := m.Called(ctx, tokenString)
	return args.Bool(0)
}

func TestRegisterUser(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	service := NewAuthService(nil, mockUserRepo, nil)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		email := "new@example.com"
		password := "password123"
		
		// Giả lập người dùng chưa tồn tại
		mockUserRepo.On("FindByEmail", ctx, email).Return(nil, nil)
		mockUserRepo.On("CreateUser", ctx, mock.AnythingOfType("*model.User")).Return(&model.User{Email: email}, nil)
		
		user, err := service.RegisterUser(ctx, email, password)
		
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, email, user.Email)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Missing Email", func(t *testing.T) {
		user, err := service.RegisterUser(ctx, "", "password123")
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "email is required", err.Error())
	})
}

func TestLoginUser(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	service := NewAuthService(nil, mockUserRepo, nil)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		email := "test@example.com"
		password := "password123"
		hashedBytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		hashedPassword := string(hashedBytes)
		
		dbUser := &model.User{
			ID:           uuid.New(),
			Email:        email,
			PasswordHash: hashedPassword,
		}

		mockUserRepo.On("FindByEmail", ctx, email).Return(dbUser, nil)
		
		user, err := service.LoginUser(ctx, email, password)
		
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, email, user.Email)
	})

	t.Run("Invalid Password", func(t *testing.T) {
		email := "test@example.com"
		mockUserRepo.On("FindByEmail", ctx, email).Return(&model.User{Email: email, PasswordHash: "wrong-hash"}, nil)
		
		user, err := service.LoginUser(ctx, email, "wrong-password")
		
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "invalid email or password", err.Error())
	})
}
