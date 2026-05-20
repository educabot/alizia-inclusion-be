package entities

import "time"

type AdaptationResource struct {
	ID           int64     `json:"id" gorm:"primaryKey"`
	AdaptationID int64     `json:"adaptation_id"`
	Title        string    `json:"title"`
	FileURL      string    `json:"file_url"`
	FileType     string    `json:"file_type" gorm:"default:pdf"`
	CreatedAt    time.Time `json:"created_at"`
}
