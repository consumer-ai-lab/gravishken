package main

import (
	config "server/config"
	route "server/src/routes"

	helmet "github.com/danielkov/gin-helmet"
	"github.com/joho/godotenv"

	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	router := SetupRouter()
	log.Fatal(router.Run(":" + os.Getenv("GO_PORT")))
	
}

func SetupRouter() *gin.Engine {
	db, err := config.Connection()

	if err != nil {
		log.Fatal("Error connecting to MongoDB: ", err)

		return nil
	}

	router := gin.Default()

	if os.Getenv("GO_ENV") == "dev" {
		gin.SetMode(gin.DebugMode)
	} else if os.Getenv("GO_ENV") == "test" {
		gin.SetMode(gin.TestMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
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
