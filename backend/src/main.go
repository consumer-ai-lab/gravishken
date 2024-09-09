package main

import (
	"os"
	config "server/config"
	route "server/src/routes"

	helmet "github.com/danielkov/gin-helmet"
	"github.com/joho/godotenv"

	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

// build time configuration. these get set using -ldflags in build script
var build_mode string

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file")
	}
	log.Println("I am updated code..............")
	log.Println("Loaded .env file")

	port, ok := os.LookupEnv("SERVER_PORT")
	if !ok {
		log.Fatalln("SERVER_PORT not set")
	}

	router := SetupRouter()
	log.Fatal(router.Run(":" + port))
}

func SetupRouter() *gin.Engine {
	db, err := config.Connection()

	if err != nil {
		log.Fatal("Error connecting to MongoDB: ", err)

		return nil
	}

	router := gin.Default()

	if build_mode == "DEV" {
		gin.SetMode(gin.DebugMode)
		// } else if build_mode == "test" {
		// 	gin.SetMode(gin.TestMode)
	} else if build_mode == "PROD" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		panic("invalid BUILD_MODE")
	}

	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"*"},
		AllowHeaders:  []string{"*"},
		AllowWildcard: true,
	}))
	router.Use(helmet.Default())
	router.Use(gzip.Gzip(gzip.BestCompression))

	route.InitAuthRoutes(db, router)
	// route.InitOtherRoutes(db, router)

	return router
}
