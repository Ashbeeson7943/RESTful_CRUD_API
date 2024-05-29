package database

import (
	"context"
	"fmt"
	"log"

	internalConfig "github.com/Ashbeeson7943/RESTful_CRUD_API/config"
	"github.com/Ashbeeson7943/RESTful_CRUD_API/data"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DatabaseConfig struct {
	KEYS   *mongo.Collection
	USERS  *mongo.Collection
	ALBUMS *mongo.Collection
}

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

func SeedKeyCollection(db_config DatabaseConfig, seed []data.ApiKey) {
	//Seed the key collection
	db_config.KEYS.Drop(context.TODO())
	a := []interface{}{seed[0]}
	_, err := db_config.KEYS.InsertMany(context.TODO(), a)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("key info inserted")
}
