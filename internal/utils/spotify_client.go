package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"music_api/config"
	"music_api/internal/logger"
	"music_api/internal/models"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
)

type ByPopularity []models.SpotifyMusicTrack

// Simulated function to fetch track details from Spotify based on ISRC
func GetSpotifyTrackDetails(isrc string, c *gin.Context) (*models.SpotifyMusicTrack, error) {
	// Simulated logic to fetch track details from Spotify API
	// Construct query parameters with the ISRC
	queryParams := url.Values{}
	queryParams.Set("q", fmt.Sprintf("isrc:%s", isrc))
	queryParams.Set("type", "track")
	fullURL := fmt.Sprintf("%s?%s", config.Config.SpotifyAPIURL+"/v1/search", queryParams.Encode())

	// Create a new GET requests
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		logger.Error("Error creating request:", err)
		return nil, err
	}
	// allowing to call spotify api with token from cookie or header
	authToken, err := c.Request.Cookie("access_token")
	token := ""
	if err != nil {
		token = strings.Split(c.Request.Header.Get("Authorization"), " ")[1]
	} else {
		token = authToken.Value
	}

	// Set access token in headers
	if token == "" {
		logger.Error("Error getting access token:", err)
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	// Create an HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("Error sending request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Error reading response body:", err)
		return nil, err
	}

	var track models.SpotifyTrackResponse

	// Unmarshal JSON response into the custom struct
	err = json.Unmarshal(body, &track)
	if err != nil {
		logger.Error("Error decoding JSON:", err)
		return nil, err
	}
	music_track := track.Tracks.Items
	sort.Sort(ByPopularity(music_track))

	if resp.StatusCode == http.StatusOK {
		return &track.Tracks.Items[0], nil
	}
	return nil, errors.New("Track details not found")
}

func (a ByPopularity) Len() int {
	return len(a)
}
func (a ByPopularity) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a ByPopularity) Less(i, j int) bool {
	return a[i].Popularity > a[j].Popularity
}
