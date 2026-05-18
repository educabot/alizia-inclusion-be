package providers

import bcerrors "github.com/educabot/team-ai-toolkit/errors"

var (
	ErrNotFound     = bcerrors.ErrNotFound
	ErrValidation   = bcerrors.ErrValidation
	ErrUnauthorized = bcerrors.ErrUnauthorized
	ErrForbidden    = bcerrors.ErrForbidden
	ErrDuplicate    = bcerrors.ErrDuplicate
	ErrConflict     = bcerrors.ErrConflict
)

var (
	ErrProfileNotFound    = bcerrors.New("student inclusion profile not found")
	ErrServiceUnavailable = bcerrors.New("AI service unavailable")
	ErrInvalidCredentials = bcerrors.New("invalid email or password")
	ErrClassroomNotFound  = bcerrors.New("classroom not found")
	ErrAdaptationNotFound = bcerrors.New("adaptation not found")
)
