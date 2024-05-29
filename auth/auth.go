package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Ashbeeson7943/RESTful_CRUD_API/data"
	database "github.com/Ashbeeson7943/RESTful_CRUD_API/db"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	FORBIDDEN string = "forbidden"
	ALLOWED   string = "allowed"
	NEW       string = "new"
)

type Access struct {
	TYPE string
}

var access Access

var DB_config database.DatabaseConfig

func FindAndValidateAPIKey(c *gin.Context) {
	var dbKey data.ApiKey
	requestAPIKey := c.Request.Header.Get("X-API-Key")
	if requestAPIKey != "" {
		filter := bson.D{{Key: "value", Value: requestAPIKey}}
		err := DB_config.KEYS.FindOne(context.TODO(), filter).Decode(&dbKey)
		if err != nil {
			fmt.Println("API Key not found in DB")
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Invalid API Key"})
			access = Access{
				TYPE: FORBIDDEN,
			}
			return
		}
		if !dbKey.VALID {
			fmt.Println("API Key not Valid")
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Invalid API Key"})
			access = Access{
				TYPE: FORBIDDEN,
			}
			return
		}
		access = Access{
			TYPE: ALLOWED,
		}
	} else {
		fmt.Println("no api key header present")
		access = Access{
			TYPE: NEW,
		}
	}
}

func AddAPIKey(c *gin.Context) {
	if !ValidateAccess(access, NEW) {
		fmt.Println("No new Key for you")
		return
	}
	username := c.Param("username")
	password := c.Param("password")

	var user database.User

	filter := bson.D{{Key: "username", Value: username}}
	err := DB_config.USERS.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		fmt.Println("User not found in DB")
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "user not found in db"})
		return
	}

	if password != user.PASSWORD {
		fmt.Println("password did not match")
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "incorrect password"})
		return
	}

	apiKey := &data.ApiKey{
		OWNER: user.USERNAME,
		VALUE: "key",
		VALID: true,
	}

	res, err := DB_config.KEYS.InsertOne(context.TODO(), apiKey)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
	c.IndentedJSON(http.StatusCreated, gin.H{"apiKey": apiKey.VALUE})
}

func ValidateAccess(access Access, supportedType string) bool {
	if access.TYPE != supportedType {
		return false
	} else {
		return true
	}
}

func GetAccess() *Access {
	return &access
}
