package middlewares

import (
	"log"
	"net/http"
	"zate/constants"
	"zate/database"
	"zate/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var authDumpCollection *mongo.Collection = database.OpenCollection(database.Client, constants.AUTH_DUMP_COLLECTION)

func VerifyEmailAuthentication() gin.HandlerFunc {

	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No authorization header"})
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(token)

		if err != "" {
			defer deleteDumpUser(c)

			c.JSON(http.StatusUnauthorized, gin.H{"error": err})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("uid", claims.Uid)

		c.Next()
	}

}

func deleteDumpUser(ctx *gin.Context) {
	var email string = ctx.GetString("email")
	_, err := authDumpCollection.DeleteOne(ctx.Request.Context(), bson.M{"email": email})
	if err != nil {
		log.Println("Unable to delete dump user with email" + email)
	}
}
