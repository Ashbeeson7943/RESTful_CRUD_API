package main

import (
	"flag"
	"fmt"
	"os"

	internalConfig "github.com/Ashbeeson7943/RESTful_CRUD_API/config"
	"github.com/Ashbeeson7943/RESTful_CRUD_API/data"
	database "github.com/Ashbeeson7943/RESTful_CRUD_API/db"
	album_routes "github.com/Ashbeeson7943/RESTful_CRUD_API/routes"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)


var collection *mongo.Collection

func main() {

	configPathPtr := flag.String("config", "./config.json", "file path to config file")

	flag.Parse()

	if *configPathPtr == ""{
		fmt.Println("no config path found")
		os.Exit(1)
	}

	config := internalConfig.LoadConfig(*configPathPtr)


	router := gin.Default()

	collection = database.LaunchDB(config)
	database.SeedDB(*collection, data.Albums)

	album_routes.Collection = collection

    router.GET(config.API_CONFIG.BASE_PATH, album_routes.GetAlbums)
	router.GET(fmt.Sprintf("%s/:id", config.API_CONFIG.BASE_PATH), album_routes.GetAlbumByID)
	router.POST(config.API_CONFIG.BASE_PATH, album_routes.PostAlbum)
	router.PUT(fmt.Sprintf("%s/:id", config.API_CONFIG.BASE_PATH), album_routes.UpdateAlbumByID)
	router.DELETE(fmt.Sprintf("%s/:id", config.API_CONFIG.BASE_PATH), album_routes.DeleteAlbumByID)

    router.Run(config.API_CONFIG.FullURL())

}
