package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

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
