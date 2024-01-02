package models

// SpotifyTrackResponse represents the structure of a simulated response from the Spotify API.
type SpotifyTrackResponse struct {
	// Tracks represents a collection of Spotify music tracks.
	Tracks struct {
		Items    []SpotifyMusicTrack `json:"items"`
		Href     string              `json:"href"`
		Limit    int                 `json:"limit"`
		Next     string              `json:"next"`
		Offset   int                 `json:"offset"`
		Previous string              `json:"previous"`
		Total    int                 `json:"total"`
	} `json:"tracks"`
}

// External_ids represents a map of external identifiers for a Spotify music track.
type External_ids map[string]string

// SpotifyMusicTrack represents a Spotify music track.
type SpotifyMusicTrack struct {
	ID           string       `json:"id" gorm:"primaryKey"` // Primary key
	Title        string       `json:"name"`
	ImageURI     string       `json:"imageUri"`
	Popularity   int          `json:"popularity"`
	Type         string       `json:"type"`
	URI          string       `json:"uri"`
	External_ids External_ids `json:"external_ids"`
	Artists      []Artist     `json:"artists"`
}
