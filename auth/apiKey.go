package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	internalConfig "github.com/Ashbeeson7943/RESTful_CRUD_API/config"
	database "github.com/Ashbeeson7943/RESTful_CRUD_API/db"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var access Access

var DB_config database.DatabaseConfig

var Key_config internalConfig.KeyConfig

type ApiKey struct {
	ID    string
	OWNER string
	VALUE string
	VALID bool
}

func (k *ApiKey) IsValid() error {
	if !k.VALID {
		return errors.New("api key is no longer valid")
	} else {
		return nil
	}
}

// Checks the users key and validates it before assigning access "level"
// NOTE: Access Levels not implemented so just some basic checks are done
func FindAndValidateAPIKey(c *gin.Context) {
	var dbKey ApiKey
	requestAPIKey := c.Request.Header.Get("X-API-Key")
	if requestAPIKey != "" {
		filter := bson.D{{Key: "value", Value: requestAPIKey}}
		err := DB_config.KEYS.FindOne(context.TODO(), filter).Decode(&dbKey)
		if err != nil {
			fmt.Println("API Key not found in DB")
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "API Key not found"})
			access = Access{
				TYPE: FORBIDDEN,
			}
			return
		}
		if !dbKey.VALID {
			fmt.Println("API Key not Valid")
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Found API Key is invalid please Generate a new one"})
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

// Find Keys attached to User
func GetApiKey(c *gin.Context) {
	pass := (!ValidateAccess(access, NEW) || !ValidateAccess(access, ALLOWED))
	if !pass {
		fmt.Println("No Finding a key for you")
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

	var keys []*ApiKey
	filter = bson.D{{Key: "owner", Value: user.USERNAME}}
	findOptions := options.Find()

	cur, err := DB_config.KEYS.Find(context.TODO(), filter, findOptions)
	if err != nil {
		fmt.Println("API Key not found in DB")
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "API Key not found"})
		access = Access{
			TYPE: FORBIDDEN,
		}
		return
	}

	for cur.Next(context.TODO()) {
		var k ApiKey
		err := cur.Decode(&k)
		if err != nil {
			log.Fatal(err)
		}
		keys = append(keys, &k)
	}

	fmt.Println(keys)
	if len(keys) > 0 {
		c.IndentedJSON(http.StatusOK, gin.H{"ApiKeys": keys})
	} else {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "No keys assigned to this user"})
	}

}

// Generate new API Key for user
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

	apiKey := &ApiKey{
		ID:    uuid.NewString(),
		OWNER: user.USERNAME,
		VALUE: string(generateKey()),
		VALID: true,
	}

	res, err := DB_config.KEYS.InsertOne(context.TODO(), apiKey)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
	c.IndentedJSON(http.StatusCreated, gin.H{"apiKey": apiKey.VALUE})
}

// Invalidate a users key
func InvalidateKey(c *gin.Context) {
	if !ValidateAccess(access, ALLOWED) {
		fmt.Println("No invalidating Keys for you")
		return
	}

	requestAPIKey := c.Request.Header.Get("X-API-Key")
	filter := bson.D{{Key: "value", Value: requestAPIKey}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "valid", Value: false}}}}
	opts := options.Update().SetUpsert(true)
	_, err := DB_config.KEYS.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found", "error": err})
		return
	}
	c.IndentedJSON(http.StatusAccepted, gin.H{"message": "API Key invalidated"})
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

func generateKey() []byte {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	key := make([]byte, Key_config.KEY_LENGTH)
	for i := range key {
		if i%9 == 0 && i != 0 {
			key[i] = '-'
		} else {
			key[i] = charset[seededRand().Intn(len(charset))]
		}
	}
	return key
}

func seededRand() *rand.Rand {
	time.Sleep(50 * time.Millisecond)
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}
