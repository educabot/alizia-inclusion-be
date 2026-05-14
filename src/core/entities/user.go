package entities

import "github.com/google/uuid"

type User struct {
	ID             int64     `json:"id" gorm:"primaryKey"`
	OrganizationID uuid.UUID `json:"organization_id"`
	Email          string    `json:"email"`
	Name           string    `json:"name"`
	PasswordHash   string    `json:"-"`
	Role           string    `json:"role" gorm:"type:member_role"`
	TimeTrackedEntity
}
