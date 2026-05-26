package inclusion

import (
	"fmt"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

var (
	errOrgIDRequired        = fmt.Errorf("%w: organization_id is required", providers.ErrValidation)
	errStudentIDRequired    = fmt.Errorf("%w: student_id is required", providers.ErrValidation)
	errClassroomIDRequired  = fmt.Errorf("%w: classroom_id is required", providers.ErrValidation)
	errMessageRequired      = fmt.Errorf("%w: message is required", providers.ErrValidation)
	errSubjectRequired      = fmt.Errorf("%w: subject is required", providers.ErrValidation)
	errNameRequired         = fmt.Errorf("%w: name is required", providers.ErrValidation)
	errTeacherIDRequired    = fmt.Errorf("%w: teacher_id is required", providers.ErrValidation)
	errAdaptationIDRequired = fmt.Errorf("%w: adaptation_id is required", providers.ErrValidation)
	errModeRequired         = fmt.Errorf("%w: mode is required", providers.ErrValidation)
	errUserIDRequired       = fmt.Errorf("%w: user_id is required", providers.ErrValidation)
	errInvalidStatus        = fmt.Errorf("%w: status must be one of: en_curso, probado, funciono, para_ajustar", providers.ErrValidation)
	errInvalidExportFormat  = fmt.Errorf("%w: format must be one of: md, pdf", providers.ErrValidation)
)
