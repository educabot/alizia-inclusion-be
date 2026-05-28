package management_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/management"
)

var (
	errDB                = fmt.Errorf("db error")
	errClassroomNotFound = fmt.Errorf("%w: classroom 99", providers.ErrNotFound)
)

func TestValidationErrors(t *testing.T) {
	tests := []struct {
		name string
		req  interface{ Validate() error }
	}{
		{"ListClassrooms_empty", management.ListClassroomsRequest{}},
		{"GetClassroom_empty", management.GetClassroomRequest{}},
		{"CreateClassroom_empty", management.CreateClassroomRequest{}},
		{"UpdateClassroom_empty", management.UpdateClassroomRequest{}},
		{"DeleteClassroom_empty", management.DeleteClassroomRequest{}},
		{"ListTeachers_empty", management.ListTeachersRequest{}},
	}
	for _, tt := range tests {
		err := tt.req.Validate()
		assert.ErrorIs(t, err, providers.ErrValidation, tt.name)
	}
}
