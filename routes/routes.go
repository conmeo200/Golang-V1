package routes


import(
	"github.com/conmeo200/Golang-V1/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/users/signup", controllers.SignUp())
	incomingRoutes.POST("/users/login", controllers.LogIn())
	incomingRoutes.POST("/admin/add-product", controllers.ProductViewerAdmin())
	incomingRoutes.POST("/users/product-view", controllers.SearchProduct())
	incomingRoutes.POST("/users/search", controllers.SearchProductByQuery())
}