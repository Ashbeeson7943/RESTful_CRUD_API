package auth

import (
	"context"
	"fmt"
	"log"

	database "github.com/Ashbeeson7943/RESTful_CRUD_API/db"
)

const (
	FORBIDDEN string = "forbidden"
	ALLOWED   string = "allowed"
	NEW       string = "new"
)

type Access struct {
	TYPE string
}

var Expired bool

func SeedKeyCollection(db_config database.DatabaseConfig, seed []ApiKey) {
	//Seed the key collection
	db_config.KEYS.Drop(context.TODO())
	a := []interface{}{seed[0]}
	_, err := db_config.KEYS.InsertMany(context.TODO(), a)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("key info inserted")
}

func GetAccess() *Access {
	return &access
}

func SetAccess(a Access) {
	access = a
}
