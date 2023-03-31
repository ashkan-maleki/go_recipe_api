package users

import (
	"context"
	"crypto/sha256"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type UserContext struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewUserContext(ctx context.Context,
	collection *mongo.Collection) *UserContext {
	return &UserContext{collection: collection, ctx: ctx}
}

func (userContext *UserContext) SetUp() {
	if userContext.usersExist() == true {
		log.Println("Users exist.")
		return
	}
	userContext.createUsers()
}

func (userContext *UserContext) usersExist() bool {
	return true
}

func (userContext *UserContext) createUsers() {
	users := map[string]string{
		"admin":      "fCRmh4Q2J7Rseqkz",
		"packt":      "RE4zfHB35VPtTkbT",
		"mlabouardy": "L3nSFRcZzNQ67bcc",
	}

	h := sha256.New()

	for username, password := range users {
		userContext.collection.InsertOne(userContext.ctx, bson.M{
			"username": username,
			"password": string(h.Sum([]byte(password))),
		})
	}
}
