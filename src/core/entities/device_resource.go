package entities

import "time"

type DeviceResource struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	DeviceID  int64     `json:"device_id"`
	Title     string    `json:"title"`
	FileURL   string    `json:"file_url"`
	FileType  string    `json:"file_type" gorm:"default:pdf"`
	CreatedAt time.Time `json:"created_at"`
}
