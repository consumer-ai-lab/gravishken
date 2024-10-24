package main

import (
	// "common"
	// "context"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	assets "server"
	"strings"
	"time"

	types "common"
	"path/filepath"

	helmet "github.com/danielkov/gin-helmet"
	"github.com/joho/godotenv"

	// "go.mongodb.org/mongo-driver/bson"

	"log"

	"io"
	"net/http/httptest"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

// build time configuration. these get set using -ldflags in build script
var build_mode string

func main() {
	// NOTE: we need a .env with SERVER_URL for the admin panel to work correctly
	_ = godotenv.Overload()

	if build_mode == "DEV" {
		root, ok := os.LookupEnv("PROJECT_ROOT")
		if !ok {
			panic("'PROJECT_ROOT' not set")
		}
		ts_dir := filepath.Join(root, "common", "ts")
		types.DumpTypes(ts_dir)
	}

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
	db, err := connectDatabase()
	// db.UserCollection.InsertOne(context.TODO(), common.User{
	// 	Username:  "test4",
	// 	Password:  "test",
	// 	BatchName: "testbatch4",
	// })
	// db.TestCollection.InsertOne(context.TODO(), common.Test{
	// 	Id:         "typing test id4",
	// 	Type:       common.TypingTest,
	// 	Duration:   500,
	// 	TypingText: "some text to type4",
	// })
	// db.TestCollection.InsertOne(context.TODO(), common.Test{
	// 	Id:         "typing test id41",
	// 	Type:       common.TypingTest,
	// 	Duration:   500,
	// 	TypingText: "some text to type41",
	// })
	// db.BatchCollection.InsertOne(context.TODO(), common.Batch{
	// 	Id:    "testbatchid4",
	// 	Name:  "testbatch4",
	// 	Tests: []string{"typing test id4", "typing test id41"},
	// })

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

	allowOrigins := getEnvOrDefault("CORS_ALLOW_ORIGINS", "http://localhost:6200")
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

	// Add a middleware to log request details
	router.Use(func(c *gin.Context) {
		log.Printf("Received request: %s %s", c.Request.Method, c.Request.URL.Path)
		log.Printf("Request headers: %v", c.Request.Header)
		c.Next()
		log.Printf("Response status: %d", c.Writer.Status())
		log.Printf("Response headers: %v", c.Writer.Header())
	})

	router.Use(helmet.Default())
	router.Use(gzip.Gzip(gzip.BestCompression))

	InitAuthRoutes(db, router)
	// route.InitOtherRoutes(db, router)

	AppRoutes(router)
	WebsiteRoutes(router)

	return router
}

func DownloadRoutes(route *gin.Engine) {
	releaseRoute := route.Group("/release")

	releaseRoute.GET("/latest/:os", func(ctx *gin.Context) {
		owner := "consumer-ai-lab"
		repo := "gravishken"
		targetOS := strings.ToLower(ctx.Param("os"))
		apiUrl := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)

		resp, err := http.Get(apiUrl)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch latest release"})
			return
		}
		defer resp.Body.Close()

		log.Println(resp.Status)

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
			return
		}
		log.Println(string(body))

		var release struct {
			Assets []struct {
				Name               string `json:"name"`
				BrowserDownloadURL string `json:"browser_download_url"`
			} `json:"assets"`
		}

		if err := json.Unmarshal(body, &release); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse release data"})
			return
		}
		log.Println(release)

		var targetAsset struct {
			Name string
			URL  string
		}

		filename := ""
		if targetOS == "windows" {
			filename = "GravishkenSetup.exe"
		}

		log.Println(filename)

		for _, asset := range release.Assets {
			if asset.Name == filename {
				targetAsset.Name = asset.Name
				targetAsset.URL = asset.BrowserDownloadURL
				break
			}
		}

		log.Printf("redirecting to %s: %s\n", targetAsset.Name, targetAsset.URL)

		if targetAsset.URL == "" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("No release found for OS: %s", targetOS)})
			return
		}

		ctx.Redirect(http.StatusFound, targetAsset.URL)
	})
}

func WebsiteRoutes(router *gin.Engine) {
	DownloadRoutes(router)

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
