package routes

import (
	"github.com/SherzodAbdullajonov/restuarant-management/controller"
	"github.com/gin-gonic/gin"
)

func FoodRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/foods", controller.GetFood())
	incomingRoutes.GET("/foods/:food_id", controller.GetFood())
	incomingRoutes.POST("foods", controller.CreateFood())
	incomingRoutes.PATCH("/foods/:food_id", controller.UpdateFood())
}
