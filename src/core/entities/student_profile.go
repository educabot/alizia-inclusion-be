package entities

import "github.com/lib/pq"

type StudentProfile struct {
	ID              int64          `json:"id" gorm:"primaryKey"`
	StudentID       int64          `json:"student_id"`
	IsTransitory    bool           `json:"is_transitory"`
	Difficulties    pq.StringArray `json:"difficulties" gorm:"type:text[]"`
	FreeDescription *string        `json:"free_description,omitempty"`
	// Rich needs layer (HU-2, all optional). situation_codes uses the controlled
	// vocabulary from situations_catalog; remaining fields enrich student context
	// (strengths, interests, triggers, strategies, environment) with no required values.
	SupportLevel            *string        `json:"support_level,omitempty"`
	Strengths               pq.StringArray `json:"strengths,omitempty" gorm:"type:text[]"`
	Interests               pq.StringArray `json:"interests,omitempty" gorm:"type:text[]"`
	Triggers                pq.StringArray `json:"triggers,omitempty" gorm:"type:text[]"`
	EffectiveStrategies     pq.StringArray `json:"effective_strategies,omitempty" gorm:"type:text[]"`
	IneffectiveStrategies   pq.StringArray `json:"ineffective_strategies,omitempty" gorm:"type:text[]"`
	SituationCodes          pq.StringArray `json:"situation_codes,omitempty" gorm:"type:text[]"`
	HasTherapeuticCompanion *bool          `json:"has_therapeutic_companion,omitempty"`
	EnvironmentNotes        *string        `json:"environment_notes,omitempty"`
	TimeTrackedEntity
}

func (StudentProfile) TableName() string {
	return "student_profiles"
}
