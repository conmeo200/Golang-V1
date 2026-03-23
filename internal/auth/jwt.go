package auth

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
	"os"
)

var secretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

type AccessClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

var AccessTokenExpire  = 60 * time.Minute
var RefreshTokenExpire = 7 * 24 * time.Hour

func GenerateTokens(userID string) (string, string, error) {

	accessClaims := AccessClaims{
		UserID: userID,
		Role:   "user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenExpire)),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessString, err := accessToken.SignedString(secretKey)
	if err != nil {
		return "", "", err
	}

	refreshClaims := RefreshClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(RefreshTokenExpire)),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshString, err := refreshToken.SignedString(secretKey)
	if err != nil {
		return "", "", err
	}

	return accessString, refreshString, nil
}

func ValidateToken(tokenString string) (*jwt.Token, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	return token, err
}