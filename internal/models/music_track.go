package models

import "gorm.io/gorm"

// Artist represents an artist in the music track.
type Artist struct {
	gorm.Model
	Id      string `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	URI     string `json:"uri"`
	TrackID uint
}

// MusicTrack represents a music track.
type MusicTrack struct {
	gorm.Model
	ISRC   string   `json:"isrc" gorm:"primaryKey"`
	Title  string   `json:"name"`
	URI    string   `json:"uri"`
	Artist []Artist `gorm:"foreignKey:TrackID"` // Foreign key
}
