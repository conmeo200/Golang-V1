package controllers


import (

)

func HashPassword (password string) string {}

func VerifyPassword (userPassword string, givenPassword string) (bool, string) {}

func Signup () gin.HandlerFunc {

	return func(c *gin.Context) {
		
	}
}

func Login () gin.HandlerFunc {}

func ProductViewerAdmin () gin.HandlerFunc {}

func SearchProduct () gin.HandlerFunc {}

func SearchProductByQuery () gin.HandlerFunc {}

