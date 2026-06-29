package entities

import "time"

type DeviceVideo struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	DeviceID  int64     `json:"device_id"`
	Title     *string   `json:"title,omitempty"`
	URL       string    `json:"url"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
}
