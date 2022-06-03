package main

import (
	"github.com/SherzodAbdullajonov/restuarant-management/middleware"
	"github.com/SherzodAbdullajonov/restuarant-management/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	//port := "8000"
	router := gin.Default()
	router.Use(gin.Logger())
	routes.UseRoutes(router)
	router.Use(middleware.Authentication())

	routes.FoodRoutes(router)
	routes.MenuRoutes(router)
	routes.TableRoutes(router)
	routes.OrderRoutes(router)
	routes.OrderItemRoutes(router)
	routes.InvoiceRoutes(router)

	router.Run(":8000")
}
