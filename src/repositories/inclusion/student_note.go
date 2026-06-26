package inclusion

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type studentNoteRepo struct {
	db *gorm.DB
}

func NewStudentNoteRepo(db *gorm.DB) providers.StudentNoteProvider {
	return &studentNoteRepo{db: db}
}

func (r *studentNoteRepo) ListByStudent(ctx context.Context, orgID uuid.UUID, studentID int64) ([]entities.StudentNote, error) {
	var notes []entities.StudentNote
	err := r.db.WithContext(ctx).
		Where("organization_id = ? AND student_id = ?", orgID, studentID).
		Order("created_at DESC").
		Find(&notes).Error
	if err != nil {
		return nil, err
	}
	return notes, nil
}

func (r *studentNoteRepo) Create(ctx context.Context, note *entities.StudentNote) error {
	return r.db.WithContext(ctx).Create(note).Error
}
