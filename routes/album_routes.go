package album_routes

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Ashbeeson7943/RESTful_CRUD_API/data"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


var Collection *mongo.Collection

//Func to get all albums stored
func GetAlbums(c *gin.Context) {
	//DB
	var results []*data.Album

	findOptions := options.Find()
	
	cur, err := Collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {
		var a data.Album
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
func PostAlbum(c *gin.Context){
	var newAlbum data.Album

	if err := c.BindJSON(&newAlbum); err != nil{
		return
	}

	//DB insert
	Collection.InsertOne(context.TODO(), newAlbum)

	c.IndentedJSON(http.StatusCreated, newAlbum)
}

//func to get specific album by ID
func GetAlbumByID(c *gin.Context){
	id := c.Param("id")
	filter := bson.D{{Key: "id", Value: id}}
	var res data.Album

	err := Collection.FindOne(context.TODO(), filter).Decode(&res)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
	} else {
		c.IndentedJSON(http.StatusOK, res)
	}
}

//func to update an album
func UpdateAlbumByID(c *gin.Context){
	id := c.Param("id")
	var newAlbum data.Album
	if err := c.BindJSON(&newAlbum); err != nil{
		return
	}

	filter := bson.M{"id": id}
	update := bson.D{{ Key: "$set", Value: bson.D{{Key: "id", Value: newAlbum.ID},{ Key: "title", Value: newAlbum.TITLE}, {Key: "artist", Value: newAlbum.ARTIST}, {Key: "price", Value: newAlbum.PRICE}}}}
	opts := options.Update().SetUpsert(true)


	_, err := Collection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found", "error" : err})
	} else {
		c.IndentedJSON(http.StatusAccepted, newAlbum)
	}
}

//func to delete album
func DeleteAlbumByID(c *gin.Context){
	id := c.Param("id")
	filter := bson.M{"id": id}
	
	_, err := Collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	msg := fmt.Sprintf("document with ID:%v deleted", id)
	c.IndentedJSON(http.StatusAccepted, gin.H{"mesage": msg})
}

