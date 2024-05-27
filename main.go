package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//region In memory data

type album struct {
	ID     string  `json:"id"`
	TITLE  string  `json:"title"`
	ARTIST string  `json:"artist"`
	PRICE  float32 `json:"price"`
}

var albums = []album{
	{ID: "1", TITLE: "Blue Train", ARTIST: "John Coltrane", PRICE: 56.99},
	{ID: "2", TITLE: "Jeru", ARTIST: "Gerry Mulligan", PRICE: 17.99},
	{ID: "3", TITLE: "Sarah Vaughan and Clifford Brown", ARTIST: "Sarah Vaughan", PRICE: 39.99},
}

//endregion


//region API func handlers

//Func to get all albums stored
func getAlbums(c *gin.Context) {
	//DB
	var results []*album

	findOptions := options.Find()
	
	cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {
		var a album
		err := cur.Decode(&a)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, &a)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.TODO())

	c.IndentedJSON(http.StatusOK, results)
}

// func to add an album
func postAlbum(c *gin.Context){
	var newAlbum album

	if err := c.BindJSON(&newAlbum); err != nil{
		return
	}

	//DB insert
	collection.InsertOne(context.TODO(), newAlbum)

	c.IndentedJSON(http.StatusCreated, newAlbum)
}

//func to get specific album by ID
func getAlbumByID(c *gin.Context){
	id := c.Param("id")
	filter := bson.D{{Key: "id", Value: id}}
	var res album

	err := collection.FindOne(context.TODO(), filter).Decode(&res)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
	} else {
		c.IndentedJSON(http.StatusOK, res)
	}
}

//func to update an album
func updateAlbumByID(c *gin.Context){
	id := c.Param("id")
	var newAlbum album
	if err := c.BindJSON(&newAlbum); err != nil{
		return
	}

	filter := bson.M{"id": id}
	update := bson.D{{ Key: "$set", Value: bson.D{{Key: "id", Value: newAlbum.ID},{ Key: "title", Value: newAlbum.TITLE}, {Key: "artist", Value: newAlbum.ARTIST}, {Key: "price", Value: newAlbum.PRICE}}}}
	opts := options.Update().SetUpsert(true)


	_, err := collection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found", "error" : err})
	} else {
		c.IndentedJSON(http.StatusAccepted, newAlbum)
	}
}

//func to delete album
func deleteAlbumByID(c *gin.Context){
	id := c.Param("id")
	filter := bson.M{"id": id}
	
	_, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	msg := fmt.Sprintf("document with ID:%v deleted", id)
	c.IndentedJSON(http.StatusAccepted, gin.H{"mesage": msg})
}

//endregion


//region DB stuff

func launchDB() {

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = c.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	collection = c.Database("tmp").Collection("albums")



}

func seedDB(){

	collection.Drop(context.TODO())

	a := []interface{}{albums[0], albums[1], albums[2]}
	r, err := collection.InsertMany(context.TODO(), a)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("album info inserted", r.InsertedIDs)
}

var collection *mongo.Collection

//endregion






func main() {
		
	router := gin.Default()

	launchDB()
	seedDB()



    router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbum)
	router.PUT("/albums/:id", updateAlbumByID)
	router.DELETE("/albums/:id", deleteAlbumByID)

    router.Run("localhost:8080")

}
