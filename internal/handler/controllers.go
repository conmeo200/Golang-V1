package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/conmeo200/Golang-V1/models"
)

func HashPassword (password string) string {}

func VerifyPassword (userPassword string, givenPassword string) (bool, string) {}

func Signup () gin.HandlerFunc {

	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validatorErr := Validate.Struct(user)
		if validatorErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validatorErr.Error()})
			return
		}
	}
}

func Login () gin.HandlerFunc {}

func ProductViewerAdmin () gin.HandlerFunc {}

func SearchProduct () gin.HandlerFunc {}

func SearchProductByQuery () gin.HandlerFunc {}

