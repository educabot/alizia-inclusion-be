package entities

import "github.com/google/uuid"

// Situation is an observable classroom behaviour (the ~15 MVP entries: "does not
// start the task", "gets distracted constantly", etc.). It is the primary
// pedagogical input — grounded in observation, not diagnosis.
// organization_id NULL means global (Educabot catalogue); non-NULL means the
// organisation has its own definition.
type Situation struct {
	ID             int64      `json:"id" gorm:"primaryKey"`
	OrganizationID *uuid.UUID `json:"organization_id,omitempty"`
	Code           string     `json:"code"`
	Name           string     `json:"name"`
	Description    *string    `json:"description,omitempty"`
	Phase          *string    `json:"phase,omitempty"`
	SortOrder      int        `json:"sort_order"`
}

func (Situation) TableName() string {
	return "situations_catalog"
}
