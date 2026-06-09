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
	// Origen IA (Capa C, HU-4): liga la sugerencia de Alizia con su resultado real.
	// was_edited = ¿el docente la modificó antes de guardar? (aceptación implícita).
	SourceConversationID *int64   `json:"source_conversation_id,omitempty"`
	SourceMessageID      *int64   `json:"source_message_id,omitempty"`
	WasEdited            bool     `json:"was_edited"`
	Student              *Student `json:"student,omitempty" gorm:"foreignKey:StudentID"`
	Teacher              *User    `json:"teacher,omitempty" gorm:"foreignKey:TeacherID"`
	Device               *Device  `json:"device,omitempty" gorm:"foreignKey:DeviceID"`
	Devices              []Device `json:"devices,omitempty" gorm:"many2many:adaptation_devices"`
	TimeTrackedEntity
}
