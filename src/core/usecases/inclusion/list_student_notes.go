package inclusion

import (
	"context"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

type ListStudentNotesRequest struct {
	OrgID     uuid.UUID
	StudentID int64
	UserID    int64
}

func (r ListStudentNotesRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return errOrgIDRequired
	}
	if r.StudentID <= 0 {
		return errStudentIDRequired
	}
	if r.UserID <= 0 {
		return errUserIDRequired
	}
	return nil
}

type ListStudentNotes interface {
	Execute(ctx context.Context, req ListStudentNotesRequest) ([]entities.StudentNote, error)
}

type listStudentNotesImpl struct {
	notes providers.StudentNoteProvider
}

func NewListStudentNotes(notes providers.StudentNoteProvider) ListStudentNotes {
	return &listStudentNotesImpl{notes: notes}
}

func (uc *listStudentNotesImpl) Execute(ctx context.Context, req ListStudentNotesRequest) ([]entities.StudentNote, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return uc.notes.ListByStudent(ctx, req.OrgID, req.StudentID, req.UserID)
}
