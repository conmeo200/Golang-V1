package auth

import (
	"crypto/rsa"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var privateKey *rsa.PrivateKey
var publicKey *rsa.PublicKey

func init() {
	// 1. Try to load from Environment Variables (useful for CI/CD)
	privPEM := os.Getenv("JWT_PRIVATE_KEY")
	pubPEM := os.Getenv("JWT_PUBLIC_KEY")

	if privPEM != "" && pubPEM != "" {
		var err error
		privateKey, err = jwt.ParseRSAPrivateKeyFromPEM([]byte(privPEM))
		if err != nil {
			fmt.Printf("Error parsing private key from env: %v\n", err)
		}
		publicKey, err = jwt.ParseRSAPublicKeyFromPEM([]byte(pubPEM))
		if err != nil {
			fmt.Printf("Error parsing public key from env: %v\n", err)
		}
		if privateKey != nil && publicKey != nil {
			return
		}
	}

	// 2. Fallback to Files
	privData, err := os.ReadFile("certs/private.pem")
	if err == nil {
		privateKey, _ = jwt.ParseRSAPrivateKeyFromPEM(privData)
	}

	pubData, err := os.ReadFile("certs/public.pem")
	if err == nil {
		publicKey, _ = jwt.ParseRSAPublicKeyFromPEM(pubData)
	}
}

func IsConfigured() bool {
	return privateKey != nil && publicKey != nil
}


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
			Issuer:    "app-issuer",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenExpire)),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, accessClaims)
	accessString, err := accessToken.SignedString(privateKey)
	if err != nil {
		return "", "", err
	}

	refreshClaims := RefreshClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "app-issuer",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(RefreshTokenExpire)),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshClaims)
	refreshString, err := refreshToken.SignedString(privateKey)
	if err != nil {
		return "", "", err
	}

	return accessString, refreshString, nil
}

func ValidateToken(tokenString string) (*jwt.Token, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	return token, err
}