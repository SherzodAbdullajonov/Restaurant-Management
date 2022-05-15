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

var tableCollection *mongo.Collection = database.OpenCollection(database.Client, "table")

func GetTables() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		result, err := tableCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing tables"})
			return
		}
		var allTables []bson.M
		if err = result.All(c, &allTables); err != nil {
			log.Fatal(err)
		}
		defer cancel()
		ctx.JSON(http.StatusOK, allTables)
	}
}
func GetTable() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		tableId := ctx.Param("table_id")
		var table models.Table

		err := tableCollection.FindOne(c, bson.M{"table_id": tableId}).Decode(&table)
		defer cancel()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching the tables"})
		}
		ctx.JSON(http.StatusOK, table)
	}
}
func CreateTable() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var table models.Table

		err := ctx.BindJSON(&table)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			defer cancel()
			return
		}
		validationErr := validate.Struct(table)
		if validationErr != nil {
			defer cancel()
			ctx.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})

			return
		}
		table.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		table.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		table.ID = primitive.NewObjectID()
		table.Table_id = table.ID.Hex()

		result, insertErr := tableCollection.InsertOne(c, table)
		if insertErr != nil {
			msg := fmt.Sprintln("message: error occured while inserting order")
			defer cancel()
			ctx.JSON(http.StatusInternalServerError, gin.H{"eror": msg})

			return
		}

		defer cancel()
		ctx.JSON(http.StatusOK, result)
	}
}
func UpdateTable() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var table models.Table
		var updateObj primitive.D
		tableId := ctx.Param("table_id")
		err := ctx.BindJSON(&table)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			defer cancel()
			return
		}
		if table.Number_of_guests != nil {
			updateObj = append(updateObj, bson.E{Key: "number_of_guests", Value: table.Number_of_guests})
		}
		if table.Table_number != nil {
			updateObj = append(updateObj, bson.E{Key: "table_number", Value: table.Table_number})

		}
		table.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: table.Updated_at})
		upsert := true
		filter := bson.M{"table_id": tableId}
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}
		result, err := tableCollection.UpdateOne(
			c,
			filter,
			bson.D{
				{Key: "$set", Value: updateObj},
			},
			&opt,
		)
		if err != nil {
			msg := fmt.Sprintln("table update failed")
			ctx.JSON(http.StatusInternalServerError, msg)
			defer cancel()
			return
		}
		defer cancel()
		ctx.JSON(http.StatusOK, result)
	}
}
