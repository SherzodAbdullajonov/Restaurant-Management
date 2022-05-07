package routes

import (
	"github.com/SherzodAbdullajonov/restuarant-management/controller"
	"github.com/gin-gonic/gin"
)

func UseRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/users", controller.GetUsers())
	incomingRoutes.GET("/users/:user_id", controller.GetUser())
	incomingRoutes.POST("/users/signup", controller.SignUp())
	incomingRoutes.POST("/users/login", controller.Login())
}
