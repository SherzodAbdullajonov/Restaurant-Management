package controller

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/SherzodAbdullajonov/restuarant-management/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

func GetUsers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		recordPerPage, err := strconv.Atoi(ctx.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		page, err1 := strconv.Atoi(ctx.Query("page"))
		if err1 != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage
		startIndex, err = strconv.Atoi(ctx.Query("startIndex"))
		// if err != nil {

		// }

		matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}
		projectStage := bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "total_count", Value: 1},
				{Key: "user_items", Value: bson.D{{Key: "$slice", Value: []interface{}{"$data", startIndex}}}},
			}}}
		result, err := userCollection.Aggregate(c, mongo.Pipeline{
			matchStage, projectStage})
		defer cancel()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing users"})
			return
		}
		var allUsers []bson.M
		if err = result.All(c, &allUsers); err != nil {
			log.Fatal(err)
		}
	}
}
func GetUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}
func SignUp() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//convert the JSON data coming from postman to something that golang understands

		//validate the data based on user struct

		//you'll check if the email has already been used by antoher user

		// hash password

		// you'll also check if the phone no. has already been used by another user

		//create some extra details for the user object - created_at, updated_at, ID

		//generate token and refresh token (generate all tokens funciton from )

		//if all ok, then you insert this new user into the user collection

		//return status OK and send the result back
	}
}
func Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//convert the login data from postman which is in JSON to golang readable format

		//find a user with that email and see if that user even exists

		//if all goes well, then you'll generate tokens

		//upadate tokens - token and refresh tokens

	}
}
func HashPassword(password string) string {
	return password
}
func VerifyPassword(userPassword string, providePassword string) (bool, string) {
	return true, providePassword
}
