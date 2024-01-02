package handlers

import (
	"music_api/internal/logger"
	"music_api/internal/models"
	"music_api/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// TrackHandler represents the handler with access to the database
type TrackHandler struct {
	DB *gorm.DB
}

// GetSuccessMessage returns a success message
func (h *TrackHandler) GetSuccessMessage(c *gin.Context) {
	logger.Info("GetSuccessMessage")
	c.JSON(http.StatusOK, gin.H{"message": "login success"})
}

// CreateTrack creates a new track with metadata from Spotify API
func (h *TrackHandler) CreateTrack(c *gin.Context) {
	// Extract ISRC from request body or query parameter
	isrc := c.Param("isrc")
	// Fetch metadata using Spotify API based on the ISRC
	spotifyMusic, err := utils.GetSpotifyTrackDetails(isrc, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Store track metadata in the database using GORM
	track := models.MusicTrack{
		ISRC:   spotifyMusic.External_ids["isrc"],
		Title:  spotifyMusic.Title,
		Artist: spotifyMusic.Artists,
		URI:    spotifyMusic.URI,
	}

	if err := h.DB.Create(&track).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// For demonstration, return a success message
	c.JSON(http.StatusCreated, gin.H{"message": "Track created successfully"})
}

// GetTrackByISRC retrieves a track by its ISRC
func (h *TrackHandler) GetTrackByISRC(c *gin.Context) {
	isrc := c.Param("isrc") // Assuming isrc is part of the route path
	// Retrieve track from the database using GORM
	var track models.MusicTrack
	if err := h.DB.Where("isrc = ?", isrc).First(&track).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Track not found"})
		return
	}
	var artists []models.Artist
	if err := h.DB.Where("track_id = ?", track.ID).Find(&artists).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Artist not found"})
		return
	}
	track.Artist = artists
	c.JSON(http.StatusOK, gin.H{"message": track})
}

// GetTracksByArtist retrieves tracks by artist name
func (h *TrackHandler) GetTracksByArtist(c *gin.Context) {
	artistName := c.Query("artist") // Assuming artist name is a query parameter
	// Retrieve tracks by artist name from the database using GORM
	// var tracks map[uint]models.MusicTrack
	tracks := make(map[uint]models.MusicTrack)

	var artists []models.Artist
	if err := h.DB.Where("name LIKE ?", "%"+artistName+"%").Find(&artists).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Artist not found"})
		return
	}
	for _, artist := range artists {
		var track models.MusicTrack
		_, ok := tracks[artist.TrackID]
		if ok {
			track = tracks[artist.TrackID]
			track.Artist = append(tracks[artist.TrackID].Artist, artist)
			tracks[track.ID] = track
			continue
		}
		if err := h.DB.Where("id = ?", artist.TrackID).First(&track).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Track not found"})
			return
		}
		track.Artist = []models.Artist{artist}
		tracks[track.ID] = track
	}

	var tracksData []models.MusicTrack

	for _, track := range tracks {
		tracksData = append(tracksData, track)
	}

	// For demonstration, return the retrieved tracks data
	// c.JSON(http.StatusOK, tracks)
	c.JSON(http.StatusOK, gin.H{"message": tracksData})
}
