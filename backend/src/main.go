package main

import (
	helmet "github.com/danielkov/gin-helmet"
	config "server/config"
	route "server/src/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var build_mode string

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	router := SetupRouter()
	log.Fatal(router.Run(":" + os.Getenv("SERVER_PORT")))
}

func SetupRouter() *gin.Engine {
	db := config.Connection()

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
