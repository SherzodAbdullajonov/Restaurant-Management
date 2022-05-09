package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/SherzodAbdullajonov/restuarant-management/database"
	"github.com/SherzodAbdullajonov/restuarant-management/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var orderCollection *mongo.Collection = database.OpenCollection(database.Client, "order")

//var tableCollection *mongo.Collection = database.OpenCollection(database.Client, "table")

func GetOrders() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		result, err := orderCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing orders"})
			return
		}
		var allOrders []bson.M
		if err = result.All(c, &allOrders); err != nil {
			log.Fatal(err)
		}

		ctx.JSON(http.StatusOK, allOrders)

	}
}
func GetOrder() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		orderId := ctx.Param("order_id")
		var order models.Order

		err := orderCollection.FindOne(c, bson.M{"order_id": orderId}).Decode(&order)
		defer cancel()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching the orders"})
		}
		ctx.JSON(http.StatusOK, order)
	}
}
func CreateOrder() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var table models.Table
		var order models.Order
		err := ctx.BindJSON(&order)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		validationErr := validate.Struct(order)
		if validationErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		}
		if order.Table_id != nil {
			err := tableCollection.FindOne(c, bson.M{"table_id": order.Table_id}).Decode(&table)
			defer cancel()
			if err != nil {
				msg := fmt.Sprint("message: Table not found")
				ctx.JSON(http.StatusInternalServerError, msg)
				return
			}
		}
		order.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		order.ID = primitive.NewObjectID()
		order.Order_id = order.ID.Hex()

		result, err := orderCollection.InsertOne(c, order)
		if err != nil {
			msg := fmt.Sprint("message: error occured while inserting order")
			ctx.JSON(http.StatusInternalServerError, gin.H{"eror": msg})
			return
		}

		defer cancel()
		ctx.JSON(http.StatusOK, result)
	}
}
func UpdateOrder() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var table models.Table
		var order models.Order

		var updateObj primitive.D
		orderId := ctx.Param("order_id")
		err := ctx.BindJSON(&order)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if order.Table_id != nil {
			err := orderCollection.FindOne(c, bson.M{"table_id": order.Table_id}).Decode(&table)
			defer cancel()
			if err != nil {
				msg := fmt.Sprint("message: Menu was not found")
				ctx.JSON(http.StatusInternalServerError, msg)
				return
			}
			updateObj = append(updateObj, bson.E{"table_id", order.Table_id})
		}
		order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updated_at", order.Updated_at})

		upsert := true
		filter := bson.M{"order_id": orderId}
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}
		result, err := orderCollection.UpdateOne(
			c,
			filter,
			bson.D{
				{"$set", updateObj},
			},
			&opt,
		)
		if err != nil {
			msg := fmt.Sprint("orders update failed")
			ctx.JSON(http.StatusInternalServerError, msg)
			return
		}
		defer cancel()
		ctx.JSON(http.StatusOK, result)
	}

}
func OrderItemOrderCreator(order models.Order) string {
	var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	order.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.ID = primitive.NewObjectID()
	order.Order_id = order.ID.Hex()

	orderCollection.InsertOne(c, order)
	defer cancel()

	return order.Order_id
}
