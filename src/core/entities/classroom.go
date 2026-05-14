package entities

import "github.com/google/uuid"

type Classroom struct {
	ID             int64     `json:"id" gorm:"primaryKey"`
	OrganizationID uuid.UUID `json:"organization_id"`
	Name           string    `json:"name"`
	Grade          *string   `json:"grade,omitempty"`
	Section        *string   `json:"section,omitempty"`
	Students       []Student `json:"students,omitempty" gorm:"foreignKey:ClassroomID"`
	TimeTrackedEntity
}
