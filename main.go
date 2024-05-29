package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Ashbeeson7943/RESTful_CRUD_API/auth"
	internalConfig "github.com/Ashbeeson7943/RESTful_CRUD_API/config"
	"github.com/Ashbeeson7943/RESTful_CRUD_API/data"
	database "github.com/Ashbeeson7943/RESTful_CRUD_API/db"
	album_routes "github.com/Ashbeeson7943/RESTful_CRUD_API/routes"
	"github.com/gin-gonic/gin"
)

func main() {

	//Set flags for running API
	configPathPtr := flag.String("config", "./config.json", "file path to config file")
	flag.Parse()

	if *configPathPtr == "" {
		fmt.Println("no config path found")
		os.Exit(1)
	}

	//Load config file
	config := internalConfig.LoadConfig(*configPathPtr)

	//Create router
	router := gin.Default()

	//Launch and seed DB
	db_config := database.LaunchDB(config)
	database.SeedAlbumCollection(db_config, data.Albums)
	database.SeedUserCollection(db_config, database.Users)
	database.SeedKeyCollection(db_config, []data.ApiKey{
		{OWNER: "test_1", VALUE: "ABC", VALID: true},
	})

	//Give access to the collection for the route handlers
	album_routes.Collection = db_config.ALBUMS
	auth.DB_config = db_config

	//Set Middleware
	router.Use(func(ctx *gin.Context) {
		auth.FindAndValidateAPIKey(ctx)
		album_routes.Access = *auth.GetAccess()
	})

	//Set routes and handler
	router.GET(config.API_CONFIG.BASE_PATH, album_routes.GetAlbums)
	router.GET(fmt.Sprintf("%s/:id", config.API_CONFIG.BASE_PATH), album_routes.GetAlbumByID)
	router.GET(fmt.Sprintf("%s/:username/:password", "/key"), auth.AddAPIKey)
	router.POST(config.API_CONFIG.BASE_PATH, album_routes.PostAlbum)
	router.PUT(fmt.Sprintf("%s/:id", config.API_CONFIG.BASE_PATH), album_routes.UpdateAlbumByID)
	router.DELETE(fmt.Sprintf("%s/:id", config.API_CONFIG.BASE_PATH), album_routes.DeleteAlbumByID)

	//Start router
	router.Run(config.API_CONFIG.FullURL())

}
