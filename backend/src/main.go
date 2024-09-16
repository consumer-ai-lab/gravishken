package main

import (
	"io/fs"
	"net/http"
	"os"
	assets "server"
	config "server/config"
	route "server/src/routes"
	"strings"
	"time"

	types "common"
	"path/filepath"

	helmet "github.com/danielkov/gin-helmet"
	"github.com/joho/godotenv"

	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"io"
	"net/http/httptest"
	"strconv"
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

	if build_mode == "DEV" {
		root, ok := os.LookupEnv("PROJECT_ROOT")
		if !ok {
			panic("'PROJECT_ROOT' not set")
		}
		ts_dir := filepath.Join(root, "common", "ts")
		types.DumpTypes(ts_dir)
	}

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

	allowOrigins := getEnvOrDefault("CORS_ALLOW_ORIGINS", "https://solid-succotash-gwjp9pr7r59265g-3000.app.github.dev/")
	allowMethods := getEnvOrDefault("CORS_ALLOW_METHODS", "GET,POST,PUT,PATCH,DELETE,HEAD,OPTIONS")
	allowHeaders := getEnvOrDefault("CORS_ALLOW_HEADERS", "Origin,Content-Length,Content-Type,Authorization")
	allowCredentials := getEnvOrDefault("CORS_ALLOW_CREDENTIALS", "true") == "true"
	maxAge := 12 * 60 * 60 // 12 hours

	router.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Split(allowOrigins, ","),
		AllowMethods:     strings.Split(allowMethods, ","),
		AllowHeaders:     strings.Split(allowHeaders, ","),
		AllowCredentials: allowCredentials,
		MaxAge:           time.Duration(maxAge) * time.Second,
		AllowWildcard:    true,
		AllowWebSockets:  true,
		AllowFiles:       true,
	}))
	router.Use(helmet.Default())
	router.Use(gzip.Gzip(gzip.BestCompression))

	route.InitAuthRoutes(db, router)
	// route.InitOtherRoutes(db, router)

	AppRoutes(router)
	AdminUiRoutes(router)

	return router
}

func AdminUiRoutes(router *gin.Engine) {
	var contentReplacements = map[string]string{
		"%SERVER_URL%": os.Getenv("SERVER_URL"),
	}

	var httpFS http.FileSystem
	if build_mode == "PROD" {
		build, _ := fs.Sub(assets.Dist, "dist")
		httpFS = http.FS(build)
	} else if build_mode == "DEV" {
		httpFS = http.Dir("dist")
	} else {
		panic("invalid BUILD_MODE")
	}

	fileServer := http.FileServer(httpFS)

	modifiedFileServer := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Serve the file and capture its content
		recorder := httptest.NewRecorder()
		fileServer.ServeHTTP(recorder, r)

		content := recorder.Body.String()

		for oldString, newString := range contentReplacements {
			content = strings.ReplaceAll(content, oldString, newString)
		}

		for k, v := range recorder.Header() {
			w.Header()[k] = v
		}
		w.Header().Set("Content-Length", strconv.Itoa(len(content)))
		if build_mode == "DEV" {
			w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
		}
		w.WriteHeader(recorder.Code)
		io.Copy(w, strings.NewReader(content))
	})

	router.NoRoute(gin.WrapH(modifiedFileServer))
}

func getEnvOrDefault(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}
