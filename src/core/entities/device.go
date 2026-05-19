package entities

import "github.com/google/uuid"

type Device struct {
	ID                 int64     `json:"id" gorm:"primaryKey"`
	OrganizationID     uuid.UUID `json:"organization_id"`
	RampID             int64     `json:"ramp_id"`
	Name               string    `json:"name"`
	Description        *string   `json:"description,omitempty"`
	ImageURL           *string   `json:"image_url,omitempty"`
	QRCode             *string   `json:"qr_code,omitempty"`
	HowToUse           *string   `json:"how_to_use,omitempty"`
	Recommendations    *string   `json:"recommendations,omitempty"`
	Rationale          *string   `json:"rationale,omitempty"`
	ClassroomBenefit   *string   `json:"classroom_benefit,omitempty"`
	NeedsDescription   *string   `json:"needs_description,omitempty"`
	UsefulWhen         *string   `json:"useful_when,omitempty"`
	EvaluationCriteria *string   `json:"evaluation_criteria,omitempty"`
	Quantity           int       `json:"quantity"`
	SortOrder          int       `json:"sort_order"`
	Ramp               *Ramp     `json:"ramp,omitempty" gorm:"foreignKey:RampID"`
	TimeTrackedEntity
}
