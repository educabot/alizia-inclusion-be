package management_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/management"
)

var (
	errDB                = errors.New("db error")
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
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if err == nil {
				t.Fatal("expected validation error, got nil")
			}
			if !errors.Is(err, providers.ErrValidation) {
				t.Errorf("expected ErrValidation, got: %v", err)
			}
		})
	}
}
