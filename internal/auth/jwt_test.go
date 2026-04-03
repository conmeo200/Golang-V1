package auth

import (
	"os"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// Thiết lập môi trường cần thiết cho test
	os.Setenv("JWT_SECRET_KEY", "your-super-secret-key-for-testing")
	//secretKey = []byte(os.Getenv("JWT_SECRET_KEY")) // Gán lại secretKey vì nó được khởi tạo ở cấp package

	code := m.Run()
	os.Exit(code)
}

func TestGenerateTokens(t *testing.T) {
	userID := "test-user-123"

	accessToken, refreshToken, err := GenerateTokens(userID)

	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)
}

func TestValidateToken(t *testing.T) {
	userID := "test-user-123"

	accessToken, _, _ := GenerateTokens(userID)

	// Test hợp lệ
	token, err := ValidateToken(accessToken)
	assert.NoError(t, err)
	assert.NotNil(t, token)
	assert.True(t, token.Valid)

	// Kiểm tra claims
	claims, ok := token.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, userID, claims["user_id"])
	assert.Equal(t, "user", claims["role"])

	// Test token không hợp lệ
	invalidToken, err := ValidateToken("invalid.token.here")
	assert.Error(t, err)
	assert.NotNil(t, invalidToken)
	assert.False(t, invalidToken.Valid)
}
