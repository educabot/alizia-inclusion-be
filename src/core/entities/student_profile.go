package entities

import "github.com/lib/pq"

type StudentProfile struct {
	ID              int64          `json:"id" gorm:"primaryKey"`
	StudentID       int64          `json:"student_id"`
	IsTransitory    bool           `json:"is_transitory"`
	Difficulties    pq.StringArray `json:"difficulties" gorm:"type:text[]"`
	FreeDescription *string        `json:"free_description,omitempty"`
	TimeTrackedEntity
}

func (StudentProfile) TableName() string {
	return "student_profiles"
}
