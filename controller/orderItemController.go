package controller

import (
	"context"
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

type OrderItemPack struct {
	Table_id    *string
	Order_items []models.OrderItem
}

var orderItemCollection *mongo.Collection = database.OpenCollection(database.Client, "order_item")

func GetOrderItems() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		result, err := orderItemCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing order items"})
			return
		}
		var allOrderItems []bson.M
		if err = result.All(c, &allOrderItems); err != nil {
			log.Fatal(err)
		}

		ctx.JSON(http.StatusOK, allOrderItems)
	}
}
func GetOrderItemByOrder() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		orderId := ctx.Param("order_id")
		allOrderItems, err := ItemByOrder(orderId)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing all order items by order ID"})
			return
		}
		ctx.JSON(http.StatusOK, allOrderItems)
	}
}
func ItemByOrder(id string) (OrderItems []primitive.M, err error) {
	var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	matchStage := bson.D{{"$mathc", bson.D{{"order_id", id}}}}
	lookupStage := bson.D{{"$lookup", bson.D{{"from", "food"}, {"localFiled", "food_id"}, {"foreignField", "food_id"}, {"as", "food"}}}}
	unwindStage := bson.D{{"$unwind", bson.D{{"path", "$food"}, {"preserveNullAndEmptyArrays", true}}}}

	lookupOrderStage := bson.D{{"$lookup", bson.D{{"from", "order"}, {"localFiled", "order_id"}, {"foreignField", "order_id"}, {"as", "order"}}}}
	unwindOrderStage := bson.D{{"$unwind", bson.D{{"path", "$order"}, {"preserveNullAndEmptyArrays", true}}}}

	lookupTableStage := bson.D{{"$lookup", bson.D{{"from", "table"}, {"localField", "order.table_id"}, {"foreignField", "table_id"}, {"as", "table"}}}}
	unwindTableStage := bson.D{{"unwind", bson.D{{"path", "$table"}, {"preserveNullAndEmptyArrays", true}}}}

	projectStage := bson.D{
		{"$project".bson.D{
			{"id", 0},
			{"amount", "$food.price"},
			{"total_count", 1},
			{"food_name", "$food.name"},
			{"food_image", "$food.food_image"},
			{"table_number", "$table.table_number"},
			{"table_id", "$table.table_id"},
			{"order_id", "$order.order_id"},
			{"price", "$food.price"},
			{"quantity,", 1},
		}}}
	groupStage := bson.D{{"$group", bson.D{{"_id", bson.D{{"order_id", "order_id"}, {"table_id", "$table_id"}, {"table_number", "$table_number"}}}, {"payment_due", bson.D{{"$sum", "$amount"}}}, {"total_count", bson.D{{"$sum", 1}}}, {"order_items", bson.D{{"$sum", "$order_items"}}}}}}

	projectStage2 := bson.D{
		{"$project", bson.D{
			{"id", 0},
			{"payment_due", 1},
			{"total_count", 1},
			{"total_number", "$_id.table_number"},
			{"order_items", 1},
		}}}
	result, err := orderItemCollection.Aggregate(c, mongo.Pipeline{
		matchStage,
		lookupStage,
		unwindStage,
		lookupOrderStage,
		unwindOrderStage,
		lookupTableStage,
		unwindTableStage,
		projectStage,
		groupStage,
		projectStage2,
	})
	if err != nil {
		panic(err)
	}
	err = result.All(c, &orderItems)
	if err != nil {
		panic(err)
	}
	defer cancel()
	return OrderItems, err

}
func GetOrderItem() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		orderItemId := ctx.Param("order_item_id")
		var orderItem models.OrderItem

		err := orderItemCollection.FindOne(c, bson.M{"orderItem_id": orderItemId})
		defer cancel()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing order item"})
			return
		}
		ctx.JSON(http.StatusOK, orderItem)
	}
}
func CreateOrderItem() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var orderItemPack OrderItemPack
		var order models.Order

		err := ctx.BindJSON(&orderItemPack)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}
		order.Order_Data, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC1123))

		orderItemsBeInserted := []interface{}{}
		order.Table_id = orderItemPack.Table_id
		order_id := OrderItemOrderCreator(order)

		for _, orderItem := range orderItemPack.Order_items {
			orderItem.Order_id = order_id
			validationErr := validate.Struct(orderItem)
			if validationErr != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
				return
			}
			orderItem.ID = primitive.NewObjectID()
			orderItem.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC1123))
			orderItem.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC1123))

			orderItem.Order_item_id = orderItem.ID.Hex()
			var num = toFixed(*orderItem.Unit_price, 1)
			orderItem.Unit_price = &num
			orderItemsBeInserted = append(orderItemsBeInserted, orderItem)
		}
		insertedOrderItems, err := orderItemCollection.InsertMany(c, orderItemsBeInserted)
		if err != nil {
			log.Fatal(err)
		}
		defer cancel()

		ctx.JSON(http.StatusOK, insertedOrderItems)
	}
}
func UpdateOrderItem() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var orderItem models.OrderItem

		orderItemId := ctx.Param("order_item_id")

		filter := bson.M{"order_item_id": orderItemId}

		var updateObj primitive.D
		if orderItem.Unit_price != nil {
			updateObj = append(updateObj, bson.E{"unit_price", orderItem.Unit_price})
		}
		if orderItem.Quantity != nil {
			updateObj = append(updateObj, bson.E{"quantity", orderItem.Quantity})
		}
		if orderItem.Food_id != nil {
			updateObj = append(updateObj, bson.E{"food_id", orderItem.Food_id})
		}

		orderItem.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC1123))
		updateObj = append(updateObj, bson.E{"updated_at", orderItem.Updated_at})

		upsert := true

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}
		result, err := orderItemCollection.UpdateOne(
			c,
			filter,
			bson.D{
				{"$set", updateObj},
			},
			&opt,
		)
		if err != nil {
			msg := "Order item update failed"
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		ctx.JSON(http.StatusOK, result)
	}
}
