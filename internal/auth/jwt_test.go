package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"os"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// 1. Kiểm tra xem đã có key chưa (từ env hoặc file)
	if !IsConfigured() {
		// 2. Nếu chưa có, tạo key tạm thời để chạy test
		reader := rand.Reader
		bitSize := 2048

		key, err := rsa.GenerateKey(reader, bitSize)
		if err == nil {
			privateKey = key
			publicKey = &key.PublicKey
		}
	}

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
