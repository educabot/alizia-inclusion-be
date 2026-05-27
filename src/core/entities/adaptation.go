package entities

import "github.com/google/uuid"

type Adaptation struct {
	ID                  int64     `json:"id" gorm:"primaryKey"`
	OrganizationID      uuid.UUID `json:"organization_id"`
	StudentID           int64     `json:"student_id"`
	TeacherID           int64     `json:"teacher_id"`
	DeviceID            *int64    `json:"device_id,omitempty"`
	Title               string    `json:"title" gorm:"default:''"`
	Subject             string    `json:"subject"`
	ActivityDescription *string   `json:"activity_description,omitempty"`
	AdaptationStrategy  *string   `json:"adaptation_strategy,omitempty"`
	AdaptationType      string    `json:"adaptation_type" gorm:"default:''"`
	Outcome             *string   `json:"outcome,omitempty"`
	Notes               *string   `json:"notes,omitempty"`
	Status              string    `json:"status" gorm:"default:en_curso"`
	Student             *Student  `json:"student,omitempty" gorm:"foreignKey:StudentID"`
	Teacher             *User     `json:"teacher,omitempty" gorm:"foreignKey:TeacherID"`
	Device              *Device   `json:"device,omitempty" gorm:"foreignKey:DeviceID"`
	Devices             []Device  `json:"devices,omitempty" gorm:"many2many:adaptation_devices"`
	TimeTrackedEntity
}
