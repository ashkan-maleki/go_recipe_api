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
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"net/http"
	"strings"
	"time"
)

var recipes []Recipe

var ctx context.Context
var err error
var client *mongo.Client
var collection *mongo.Collection

// mongodb://admin:password@127.0.0.1:27017/
const mongoUri = "mongodb://admin:password@127.0.0.1:27017/"
const mongoDatabase = "demo"

// MONGO_URI="mongodb://admin:password@127.0.0.1:27017/" MONGO_DATABASE=demo go run main.go

func init() {
	//recipes = make([]Recipe, 0)
	//file, _ := os.ReadFile("recipe.json")
	//_ = json.Unmarshal(file, &recipes)

	ctx = context.Background()
	//client, err = mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoUri))

	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDb")
	//collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")
	collection = client.Database(mongoDatabase).Collection("recipes")
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

func findRecipeIndexByID(id string) int {
	index := -1
	for i := 0; i < len(recipes); i++ {
		//if recipes[i].ID == id {
		if recipes[i].ID.String() == id {
			index = i
		}
	}
	return index
}

func findRecipeByID(id string) (*Recipe, error) {
	index := findRecipeIndexByID(id)
	if index == -1 {
		return nil, errors.New("recipe not found")
	}
	return &(recipes[index]), nil
}

func main() {
	router := gin.Default()
	router.GET("/", IndexHandler)
	router.POST("/recipes", NewRecipeHandler)
	router.GET("/recipes", ListRecipeHandler)
	router.PUT("/recipes/:id", UpdateRecipeHandler)
	router.DELETE("/recipes/:id", DeleteRecipeHandler)
	router.GET("/recipes/:id", GetRecipeHandler)
	router.GET("/recipes/search", SearchRecipeHandler)
	router.Run()
	fmt.Println("serving on http://localhost:8080/")
}

type Recipe struct {
	//swagger:ignore
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Name         string             `json:"name" bson:"name"`
	Tags         []string           `json:"tags" bson:"tags"`
	Ingredients  []string           `json:"ingredients" bson:"ingredients"`
	Instructions []string           `json:"instructions" bson:"instructions"`
	//swagger:ignore
	PublishedAt time.Time `json:"publishedAt"`
}

// swagger:parameters newRecipe
type CreateRecipe struct {
	// in: body
	Body *Recipe
}

//swagger:parameters updateRecipe
type UpdateRecipe struct {
	// in: body
	Body *Recipe
	// in: path
	Id int `json:"id"`
}

// swagger:operation GET / index indexPage
// Index
// ---
// produces:
// - application/json
// responses:
//
//	'200':
//	    description: Successful operation
func IndexHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "hello world",
	})
}

// swagger:operation GET /recipes recipes listRecipes
// Returns list of recipes
// ---
// produces:
// - application/json
// responses:
//
//	'200':
//	    description: Successful operation
func ListRecipeHandler(c *gin.Context) {
	cur, err := collection.Find(ctx, bson.M{})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer cur.Close(ctx)

	recipes := make([]Recipe, 0)
	for cur.Next(ctx) {
		var recipe Recipe
		err := cur.Decode(&recipe)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		recipes = append(recipes, recipe)
	}

	c.JSON(http.StatusOK, recipes)
}

// swagger:operation POST /recipes recipes newRecipe
// Create a new recipe
// ---
// produces:
// - application/json
//
// responses:
//
//	'200':
//	    description: Successful operation
//	'400':
//	    description: Invalid input
func NewRecipeHandler(c *gin.Context) {
	var recipe Recipe

	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	recipe.ID = primitive.NewObjectID()
	recipe.PublishedAt = time.Now()

	_, err = collection.InsertOne(ctx, recipe)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting a new recipe"})
		return
	}

	c.JSON(http.StatusOK, recipe)
}

// swagger:operation GET /recipes/{id} recipes oneRecipe
// Get one recipe
// ---
// produces:
// - application/json
// parameters:
//   - name: id
//     in: path
//     description: ID of the recipe
//     required: true
//     type: string
//
// responses:
//
//	'200':
//	    description: Successful operation
//	'404':
//	    description: Invalid recipe ID
func GetRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	recipe, err := findRecipeByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, recipe)
}

// swagger:operation PUT /recipes/{id} recipes updateRecipe
// Update an existing recipe
// ---
//
// produces:
// - application/json
// responses:
//
//	'200':
//	    description: Successful operation
//	'400':
//	    description: Invalid input
//	'404':
//	    description: Invalid recipe ID
func UpdateRecipeHandler(c *gin.Context) {
	var recipe Recipe

	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	id := c.Param("id")
	index := -1
	for i := 0; i < len(recipes); i++ {
		//if recipes[i].ID == id {
		if recipes[i].ID.String() == id {
			index = i
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
		return
	}
	//recipe.ID = id
	recipe.ID = primitive.NewObjectID()
	recipes[index] = recipe
	c.JSON(http.StatusOK, recipe)
}

// swagger:operation DELETE /recipes/{id} recipes deleteRecipe
// Delete an existing recipe
// ---
// produces:
// - application/json
// parameters:
//   - name: id
//     in: path
//     description: ID of the recipe
//     required: true
//     type: string
//
// responses:
//
//	'200':
//	    description: Successful operation
//	'404':
//	    description: Invalid recipe ID
func DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	index := -1

	for i := 0; i < len(recipes); i++ {
		//if recipes[i].ID == id {
		if recipes[i].ID.String() == id {
			index = i
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
		return
	}

	recipes = append(recipes[:index], recipes[index+1:]...)
	c.JSON(http.StatusOK, gin.H{"message": "Recipe has been deleted"})
}

// swagger:operation GET /recipes/search recipes findRecipe
// Search recipes based on tags
// ---
// produces:
// - application/json
// parameters:
//   - name: tag
//     in: query
//     description: recipe tag
//     required: true
//     type: string
//
// responses:
//
//	'200':
//	    description: Successful operation
func SearchRecipeHandler(c *gin.Context) {
	tag := c.Query("tag")
	listOfRecipes := make([]Recipe, 0)

	for i := 0; i < len(recipes); i++ {
		found := false
		for _, t := range recipes[i].Tags {
			if strings.EqualFold(t, tag) {
				found = true
			}
		}
		if found {
			listOfRecipes = append(listOfRecipes, recipes[i])
		}
	}

	c.JSON(http.StatusOK, listOfRecipes)
}
