package database

import (
	"context"
	"fmt"
	"log"
	"time"

	internalConfig "github.com/Ashbeeson7943/RESTful_CRUD_API/config"
	"github.com/Ashbeeson7943/RESTful_CRUD_API/data"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DatabaseConfig struct {
	KEYS   *mongo.Collection
	USERS  *mongo.Collection
	ALBUMS *mongo.Collection
	USAGE  *mongo.Collection
}

var DB_config DatabaseConfig

func LaunchDB(conf internalConfig.Config) DatabaseConfig {

	clientOptions := options.Client().ApplyURI(conf.DB_URI)

	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = c.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	var db_config DatabaseConfig

	db_config.ALBUMS = c.Database("tmp").Collection("albums")
	db_config.KEYS = c.Database("tmp").Collection("keys")
	db_config.USERS = c.Database("tmp").Collection("users")
	db_config.USAGE = c.Database("tmp").Collection("usage")
	db_config.USAGE.Drop(context.TODO())

	DB_config = db_config
	return db_config

}

func SeedAlbumCollection(db_config DatabaseConfig, seed []data.Album) {
	//Seed the album collection
	db_config.ALBUMS.Drop(context.TODO())
	a := []interface{}{seed[0], seed[1], seed[2]}
	_, err := db_config.ALBUMS.InsertMany(context.TODO(), a)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("album info inserted")
}

func SeedUserCollection(db_config DatabaseConfig, seed []User) {
	//Seed the user collection
	db_config.USERS.Drop(context.TODO())
	a := []interface{}{seed[0], seed[1], seed[2]}
	_, err := db_config.USERS.InsertMany(context.TODO(), a)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("user info inserted")
}

func LogKeyUsage(c *gin.Context) {
	apiKey := c.Request.Header.Get("X-API-Key")
	if apiKey != "" {
		stat := data.KeyStat{
			KEY_ID:   apiKey,
			DATETIME: string(time.Now().Format(time.DateTime)),
			METHOD:   c.Request.Method,
			PATH:     c.Request.URL.Path,
		}
		_, err := DB_config.USAGE.InsertOne(context.TODO(), stat)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		stat := data.KeyStat{
			ERR:      "UNAUTHORISED",
			DATETIME: string(time.Now().Format(time.DateTime)),
			METHOD:   c.Request.Method,
			PATH:     c.Request.URL.Path,
		}
		_, err := DB_config.USAGE.InsertOne(context.TODO(), stat)
		if err != nil {
			fmt.Println(err)
		}
	}
}
