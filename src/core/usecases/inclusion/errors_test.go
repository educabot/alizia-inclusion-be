package inclusion_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/inclusion"
)

var (
	errDB                 = errors.New("db error")
	errStudentNotFound    = fmt.Errorf("%w: student 99", providers.ErrNotFound)
	errAdaptationNotFound = fmt.Errorf("%w: adaptation 99", providers.ErrNotFound)
)

func TestValidationErrors(t *testing.T) {
	tests := []struct {
		name string
		req  interface{ Validate() error }
	}{
		{"ListStudents_empty", inclusion.ListStudentsRequest{}},
		{"CreateStudent_empty", inclusion.CreateStudentRequest{}},
		{"UpdateStudent_empty", inclusion.UpdateStudentRequest{}},
		{"DeleteStudent_empty", inclusion.DeleteStudentRequest{}},
		{"GetStudentProfile_empty", inclusion.GetStudentProfileRequest{}},
		{"UpsertStudentProfile_empty", inclusion.UpsertStudentProfileRequest{}},
		{"ListClassroomStudents_empty", inclusion.ListClassroomStudentsRequest{}},
		{"ListAdaptations_empty", inclusion.ListAdaptationsRequest{}},
		{"GetAdaptation_empty", inclusion.GetAdaptationRequest{}},
		{"CreateAdaptation_empty", inclusion.CreateAdaptationRequest{}},
		{"UpdateAdaptation_empty", inclusion.UpdateAdaptationRequest{}},
		{"DeleteAdaptation_empty", inclusion.DeleteAdaptationRequest{}},
		{"ListAdaptationResources_empty", inclusion.ListAdaptationResourcesRequest{}},
		{"GetChatHistory_empty", inclusion.GetChatHistoryRequest{}},
		{"RecommendDevice_empty", inclusion.RecommendDeviceRequest{}},
		{"AssistClassroom_empty", inclusion.AssistClassroomRequest{}},
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
