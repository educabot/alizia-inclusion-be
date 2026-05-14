package entities

import "github.com/google/uuid"

type Organization struct {
	ID   uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Name string    `json:"name"`
	TimeTrackedEntity
}
