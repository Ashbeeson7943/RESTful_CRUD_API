package usage

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Ashbeeson7943/RESTful_CRUD_API/auth"
	"github.com/Ashbeeson7943/RESTful_CRUD_API/data"
	database "github.com/Ashbeeson7943/RESTful_CRUD_API/db"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

var DB_config database.DatabaseConfig

func LogKeyUsage(c *gin.Context) {
	apiKey := c.Request.Header.Get("X-API-Key")
	if apiKey != "" {
		filter := bson.D{{Key: "value", Value: apiKey}}
		var key auth.ApiKey
		err := DB_config.KEYS.FindOne(context.TODO(), filter).Decode(&key)
		if err != nil {
			fmt.Println("API Key not found in DB")
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "API Key not found"})
			return
		}
		if key.USES > 0 {
			update := bson.D{{Key: "$inc", Value: bson.D{{Key: "uses", Value: -1}}}}
			DB_config.KEYS.UpdateOne(context.TODO(), filter, update)
			stat := data.KeyStat{
				KEY_ID:   apiKey,
				DATETIME: string(time.Now().Format(time.DateTime)),
				METHOD:   c.Request.Method,
				PATH:     c.Request.URL.Path,
			}
			_, err := DB_config.USAGE.InsertOne(context.TODO(), stat)
			if err != nil {
				c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Key not found"})
			}
			auth.Expired = false
		} else {
			auth.Expired = true
			auth.SetAccess(auth.Access{
				TYPE: auth.FORBIDDEN,
			})
			c.IndentedJSON(http.StatusForbidden, gin.H{"message": "Key has no lives remaining"})
			update := bson.D{{Key: "$set", Value: bson.D{{Key: "valid", Value: false}}}}
			DB_config.KEYS.UpdateOne(context.TODO(), filter, update)
		}

	} else {
		s := data.KeyStat{
			ERR:      "UNAUTHORISED",
			DATETIME: string(time.Now().Format(time.DateTime)),
			METHOD:   c.Request.Method,
			PATH:     c.Request.URL.Path,
		}
		fmt.Println(s)
		_, err := DB_config.USAGE.InsertOne(context.TODO(), s)
		if err != nil {
			fmt.Println(err.Error())
		}
		auth.Expired = false
	}
}
