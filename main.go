// Recipes API
//
// This is a sample recipes API. You can find out more about the API at https://github.com/PacktPublishing/Building-Distributed-Applications-in-Gin.
//
//				Schemes: http
//		 Host: localhost:8080
//				BasePath: /
//				Version: 1.0.0
//				Contact: Mohamed Labouardy <mohamed@labouardy.com> https://labouardy.com
//	 SecurityDefinitions:
//	 api_key:
//	   type: apiKey
//	   name: Authorization
//	   in: header
//
//				Consumes:
//				- application/json
//
//				Produces:
//				- application/json
//
// swagger:meta
package main

import (
	"context"
	"fmt"
	"github.com/gin-contrib/sessions"
	redisStore "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/mamalmaleki/go_recipe_api/handlers"
	"github.com/mamalmaleki/go_recipe_api/users"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
)

var authHandler *handlers.AuthHandler
var recipeHandler *handlers.RecipesHandler

// mongodb://admin:password@127.0.0.1:27017/
const mongoUri = "mongodb://admin:password@127.0.0.1:27017/"
const mongoDatabase = "demo"

// MONGO_URI="mongodb://admin:password@127.0.0.1:27017/" MONGO_DATABASE=demo go run main.go

func init() {
	ctx := context.Background()
	//client, err = mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUri))

	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDb")
	//collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")
	collection := client.Database(mongoDatabase).Collection("recipes")

	// begin

	//recipes := make([]models.Recipe, 0)
	//file, _ := os.ReadFile("recipe.json")
	//err = json.Unmarshal(file, &recipes)
	//
	//var listOfRecipes []interface{}
	//
	//for _, recipe := range recipes {
	//	listOfRecipes = append(listOfRecipes, recipe)
	//}
	//
	//log.Println(fmt.Printf("number of recipes: %v", err))
	//
	//insertManyResult, err := collection.InsertMany(ctx, listOfRecipes)
	//
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//log.Println(fmt.Printf("Inserted recipes: %v", len(insertManyResult.InsertedIDs)))

	// end

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	status := redisClient.Ping(ctx)
	log.Println(status)
	recipeHandler = handlers.NewRecipesHandler(ctx, collection, redisClient)

	collectionUsers := client.Database(mongoDatabase).Collection("users")
	authHandler = handlers.NewAuthHandler(ctx, collectionUsers)

	userContext := users.NewUserContext(ctx, collectionUsers)
	userContext.SetUp()
}

func main() {
	router := gin.Default()

	store, _ := redisStore.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))

	router.Use(sessions.Sessions("recipes_api", store))

	router.POST("/sign-in", authHandler.SignInHandler)
	router.POST("/sign-out", authHandler.SignOutHandler)
	router.POST("/refresh", authHandler.RefreshHandler)
	router.GET("/", recipeHandler.IndexHandler)

	authorized := router.Group("/")
	authorized.Use(authHandler.AuthMiddleware())
	{
		authorized.POST("/recipes", recipeHandler.NewRecipeHandler)
		authorized.GET("/recipes", recipeHandler.ListRecipeHandler)
		authorized.PUT("/recipes/:id", recipeHandler.UpdateRecipeHandler)
		authorized.DELETE("/recipes/:id", recipeHandler.DeleteRecipeHandler)
		authorized.GET("/recipes/:id", recipeHandler.GetRecipeHandler)
		authorized.GET("/recipes/search", recipeHandler.SearchRecipeHandler)
	}
	fmt.Println("serving on http://localhost:8080/")
	fmt.Println("serving on https://localhost:443/")
	router.Run()
	//router.RunTLS(":443", "certs/localhost.crt",
	//	"certs/localhost.key")

}
