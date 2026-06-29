package inclusion

import (
	"context"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

const defaultNoteType = "seguimiento"

var validNoteTypes = map[string]struct{}{
	"diagnostico": {},
	"observable":  {},
	"seguimiento": {},
}

type CreateStudentNoteRequest struct {
	OrgID     uuid.UUID
	StudentID int64
	UserID    int64
	Content   string
	Type      string
	// Internal controla si la nota es visible en el front. Las notas internas
	// (diagnóstico/observables del alta) no se muestran. Default true.
	Internal *bool
}

func (r CreateStudentNoteRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return errOrgIDRequired
	}
	if r.StudentID <= 0 {
		return errStudentIDRequired
	}
	if r.UserID <= 0 {
		return errUserIDRequired
	}
	if r.Content == "" {
		return errContentRequired
	}
	if r.Type != "" {
		if _, ok := validNoteTypes[r.Type]; !ok {
			return errInvalidNoteType
		}
	}
	return nil
}

type CreateStudentNote interface {
	Execute(ctx context.Context, req CreateStudentNoteRequest) (*entities.StudentNote, error)
}

type createStudentNoteImpl struct {
	notes providers.StudentNoteProvider
}

func NewCreateStudentNote(notes providers.StudentNoteProvider) CreateStudentNote {
	return &createStudentNoteImpl{notes: notes}
}

func (uc *createStudentNoteImpl) Execute(ctx context.Context, req CreateStudentNoteRequest) (*entities.StudentNote, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	noteType := req.Type
	if noteType == "" {
		noteType = defaultNoteType
	}
	internal := true
	if req.Internal != nil {
		internal = *req.Internal
	}

	note := &entities.StudentNote{
		OrganizationID: req.OrgID,
		StudentID:      req.StudentID,
		UserID:         req.UserID,
		Content:        req.Content,
		Type:           noteType,
		Internal:       internal,
	}
	if err := uc.notes.Create(ctx, note); err != nil {
		return nil, err
	}
	return note, nil
}
