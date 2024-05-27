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

func LaunchDB(conf internalConfig.Config) *mongo.Collection{

	clientOptions := options.Client().ApplyURI(conf.DB_URI)

	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = c.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	return c.Database("tmp").Collection("albums")

}

func SeedDB(col mongo.Collection, seed []data.Album) {

	col.Drop(context.TODO())

	a := []interface{}{seed[0], seed[1], seed[2]}
	r, err := col.InsertMany(context.TODO(), a)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("album info inserted", r.InsertedIDs)
}