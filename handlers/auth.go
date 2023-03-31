package handlers

import (
	"context"
	"crypto/sha256"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mamalmaleki/go_recipe_api/models"
	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

type AuthHandler struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewAuthHandler(ctx context.Context, collection *mongo.Collection) *AuthHandler {
	return &AuthHandler{collection, ctx}
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type JWTOutput struct {
	Token   string    `json:"token"`
	Expires time.Time `json:"expires"`
}

func (handler *AuthHandler) SignInHandler(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h := sha256.New()
	cur := handler.collection.FindOne(handler.ctx, bson.M{
		"username": user.Username,
		"password": string(h.Sum([]byte(user.Password))),
	})

	if cur.Err() != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	//expirationTime := time.Now().Add(10 * time.Minute)
	//claims := &Claims{
	//	Username: user.Username,
	//	RegisteredClaims: jwt.RegisteredClaims{
	//		ExpiresAt: jwt.NewNumericDate(expirationTime),
	//	},
	//}

	//token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//jwtSecret := os.Getenv("JWT_SECRET")
	//jwtSecret := "eUbP9shywUygMx7u"
	//tokenString, err := token.SignedString([]byte(jwtSecret))

	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//	return
	//}
	//
	//jwtOutput := JWTOutput{
	//	Token:   tokenString,
	//	Expires: expirationTime,
	//}
	//c.JSON(http.StatusOK, jwtOutput)

	sessionToken := xid.New().String()
	session := sessions.Default(c)
	session.Set("username", user.Username)
	session.Set("token", sessionToken)
	session.Save()
	c.JSON(http.StatusOK, gin.H{"message": "User signed in"})
}

func (handler *AuthHandler) RefreshHandler(c *gin.Context) {
	tokenValue := c.GetHeader("Authorization")
	claims := &Claims{}
	jwtSecret := "eUbP9shywUygMx7u"
	tkn, err := jwt.ParseWithClaims(tokenValue, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if tkn == nil || !tkn.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	if time.Unix(claims.ExpiresAt.Unix(), 0).Sub(time.Now()) > 30*time.Second {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token is not expired yet"})
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = jwt.NewNumericDate(expirationTime)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	jwtOutput := JWTOutput{
		Token:   tokenString,
		Expires: expirationTime,
	}

	c.JSON(http.StatusOK, jwtOutput)
}

func (handler *AuthHandler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//tokenValue := c.GetHeader("Authorization")
		//claims := &Claims{}
		//
		////jwtSecret := os.Getenv("JWT_SECRET")
		//jwtSecret := "eUbP9shywUygMx7u"
		//tkn, err := jwt.ParseWithClaims(tokenValue, claims, func(token *jwt.Token) (interface{}, error) {
		//	return []byte(jwtSecret), nil
		//})
		//if err != nil {
		//	c.AbortWithStatus(http.StatusUnauthorized)
		//}
		//if tkn == nil || !tkn.Valid {
		//	c.AbortWithStatus(http.StatusUnauthorized)
		//}
		//c.Next()

		session := sessions.Default(c)
		sessionToken := session.Get("token")
		if sessionToken == nil {
			c.JSON(http.StatusForbidden, gin.H{"message": "Not logged in"})
			c.Abort()
		}
		c.Next()
	}
}

func (handler *AuthHandler) SignOutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.JSON(http.StatusOK, gin.H{"message": "signed out..."})
}