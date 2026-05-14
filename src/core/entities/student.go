package entities

import "github.com/google/uuid"

type Student struct {
	ID             int64           `json:"id" gorm:"primaryKey"`
	OrganizationID uuid.UUID       `json:"organization_id"`
	ClassroomID    int64           `json:"classroom_id"`
	Name           string          `json:"name"`
	Profile        *StudentProfile `json:"profile,omitempty" gorm:"foreignKey:StudentID"`
	TimeTrackedEntity
}
