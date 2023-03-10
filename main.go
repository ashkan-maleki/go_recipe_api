// Recipes API
//
// This is a sample recipes API. You can find out more about the API at https://github.com/PacktPublishing/Building-Distributed-Applications-in-Gin.
//
//			Schemes: http
//	 Host: localhost:8080
//			BasePath: /
//			Version: 1.0.0
//			Contact: Mohamed Labouardy <mohamed@labouardy.com> https://labouardy.com
//
//			Consumes:
//			- application/json
//
//			Produces:
//			- application/json
//
// swagger:meta
package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mamalmaleki/go_recipe_api/handlers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
)

var recipeHandler *handlers.RecipesHandler

// mongodb://admin:password@127.0.0.1:27017/
const mongoUri = "mongodb://admin:password@127.0.0.1:27017/"
const mongoDatabase = "demo"

// MONGO_URI="mongodb://admin:password@127.0.0.1:27017/" MONGO_DATABASE=demo go run main.go

func init() {
	//recipes = make([]Recipe, 0)
	//file, _ := os.ReadFile("recipe.json")
	//_ = json.Unmarshal(file, &recipes)

	ctx := context.Background()
	//client, err = mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUri))

	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDb")
	//collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")
	collection := client.Database(mongoDatabase).Collection("recipes")
	recipeHandler = handlers.NewRecipesHandler(ctx, collection)
	//var listOfRecipes []interface{}
	//
	//for _, recipe := range recipes {
	//	listOfRecipes = append(listOfRecipes, recipe)
	//}
	//

	//
	//insertManyResult, err := collection.InsertMany(ctx, listOfRecipes)
	//
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//log.Println(fmt.Printf("Inserted recipes: %v", len(insertManyResult.InsertedIDs)))
}

func main() {
	router := gin.Default()
	router.GET("/", recipeHandler.IndexHandler)
	router.POST("/recipes", recipeHandler.NewRecipeHandler)
	router.GET("/recipes", recipeHandler.ListRecipeHandler)
	router.PUT("/recipes/:id", recipeHandler.UpdateRecipeHandler)
	router.DELETE("/recipes/:id", recipeHandler.DeleteRecipeHandler)
	router.GET("/recipes/:id", recipeHandler.GetRecipeHandler)
	router.GET("/recipes/search", recipeHandler.SearchRecipeHandler)
	router.Run()
	fmt.Println("serving on http://localhost:8080/")
}
