package main

import (
	"music_api/config"
	"music_api/internal/database"
	"music_api/internal/handlers"
	"music_api/internal/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ConfigFile = "app_config.json"
)

func main() {
	// Create a log file

	logger.Init()

	// Connect to the database
	_, err := LoadApplicationConfig(ConfigFile)
	if err != nil {
		panic("Failed to load configs")
	}
	handlers.Init()

	db, err := database.ConnectDB()
	if err != nil {
		logger.Error("Failed to connect to the database:", err)
	}
	defer database.CloseDB()

	// Initialize Gin router
	r := gin.Default()

	// Initialize handlers with access to the database
	trackHandler := handlers.TrackHandler{DB: db}

	// Group API routes that require authentication
	api := r.Group("/api")
	api.Use(handlers.SpotifyAuthMiddleware()) // Add OAuth2 authentication middleware

	{
		api.POST("/tracks/:isrc", trackHandler.CreateTrack)
		api.GET("/tracks/:isrc", trackHandler.GetTrackByISRC)
		api.GET("/tracks", trackHandler.GetTracksByArtist)
		api.GET("/", trackHandler.GetSuccessMessage)
		// Add other routes as needed
	}

	auth_api := r.Group("/auth")
	auth_api.GET("/spotify/callback", handlers.CallbackHandler())
	// Start the server
	logger.Info("Starting server on port 8081...")
	if err := http.ListenAndServe(":8081", r); err != nil {
		logger.Error("Failed to start the server: ", err)
	}
}

func LoadApplicationConfig(configFilePath string) (*config.Configuration, error) {
	configs := &config.Configuration{}

	err := config.Load(configFilePath)
	if err != nil {
		return nil, err
	}
	configs = config.Config
	return configs, nil
}
