package inclusion

import (
	"fmt"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

var (
	errOrgIDRequired       = fmt.Errorf("%w: organization_id is required", providers.ErrValidation)
	errStudentIDRequired   = fmt.Errorf("%w: student_id is required", providers.ErrValidation)
	errClassroomIDRequired = fmt.Errorf("%w: classroom_id is required", providers.ErrValidation)
	errMessageRequired     = fmt.Errorf("%w: message is required", providers.ErrValidation)
	errSubjectRequired     = fmt.Errorf("%w: subject is required", providers.ErrValidation)
	errObjectiveRequired   = fmt.Errorf("%w: objective is required", providers.ErrValidation)
)
