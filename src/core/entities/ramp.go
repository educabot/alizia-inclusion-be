package entities

import "github.com/google/uuid"

type Ramp struct {
	ID               int64     `json:"id" gorm:"primaryKey"`
	OrganizationID   uuid.UUID `json:"organization_id"`
	Name             string    `json:"name"`
	Description      *string   `json:"description,omitempty"`
	ShortDescription *string   `json:"short_description,omitempty"`
	SortOrder        int       `json:"sort_order"`
	Devices          []Device  `json:"devices,omitempty" gorm:"foreignKey:RampID"`
	TimeTrackedEntity
}
