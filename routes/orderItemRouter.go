package routes

import (
	controller "github.com/SherzodAbdullajonov/restuarant-management/controller"
	"github.com/gin-gonic/gin"
)

func OrderItemRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/orderItems", controller.GetOrderItems())
	incomingRoutes.GET("/orderItems/:orderItems_id", controller.GetOrderItem())
	incomingRoutes.GET("/orderItem-order/:order_id", controller.GetOrderItemByOrder())
	incomingRoutes.POST("/orderItmes", controller.CreateOrderItem())
	incomingRoutes.PATCH("/orderItems/:orderItems_id", controller.UpdateItem())
}
