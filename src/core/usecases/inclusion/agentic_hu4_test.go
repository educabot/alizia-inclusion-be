package inclusion

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	mockproviders "github.com/educabot/alizia-inclusion-be/src/mocks/providers"
)

func TestInclusionDispatcher_CreateRecursoWithoutConfirmedDoesNotPersist(t *testing.T) {
	// Arrange
	adaptations := new(mockproviders.MockAdaptationProvider)
	d := inclusionDispatcher{adaptations: adaptations, userID: 1, conversationID: 50}

	// Act
	result, err := d.Dispatch(context.Background(), uuid.New(), providers.ToolCall{
		Name:      "create_recurso",
		Arguments: `{"student_id": 9001, "subject": "Lengua", "title": "Texto narrativo"}`,
	})

	// Assert
	require.NoError(t, err)
	var got struct {
		Pending bool `json:"pending_confirmation"`
	}
	require.NoError(t, json.Unmarshal([]byte(result), &got))
	assert.True(t, got.Pending)
	adaptations.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestInclusionDispatcher_CreateRecursoConfirmedPersistsWithOrigin(t *testing.T) {
	// Arrange
	ctx := context.Background()
	orgID := uuid.New()
	adaptations := new(mockproviders.MockAdaptationProvider)
	adaptations.On("Create", ctx, mock.MatchedBy(func(a *entities.Adaptation) bool {
		return a.StudentID == 9001 && a.TeacherID == 7 &&
			a.SourceConversationID != nil && *a.SourceConversationID == 50 && a.WasEdited
	})).Return(nil)
	adaptations.On("SetDevices", ctx, mock.Anything, []int64{3}).Return(nil)
	d := inclusionDispatcher{adaptations: adaptations, userID: 7, conversationID: 50}

	// Act
	result, err := d.Dispatch(ctx, orgID, providers.ToolCall{
		Name:      "create_recurso",
		Arguments: `{"student_id": 9001, "subject": "Lengua", "device_ids": [3], "was_edited": true, "confirmed": true}`,
	})

	// Assert
	require.NoError(t, err)
	var got struct {
		Saved bool `json:"saved"`
	}
	require.NoError(t, json.Unmarshal([]byte(result), &got))
	assert.True(t, got.Saved)
	adaptations.AssertExpectations(t)
}

func TestInclusionDispatcher_CreateStudentConfirmedPersists(t *testing.T) {
	// Arrange
	ctx := context.Background()
	orgID := uuid.New()
	students := new(mockproviders.MockStudentProvider)
	students.On("Create", ctx, mock.MatchedBy(func(s *entities.Student) bool {
		return s.Name == "Nuevo" && s.ClassroomID == 9001
	})).Return(nil)
	d := inclusionDispatcher{students: students}

	// Act
	result, err := d.Dispatch(ctx, orgID, providers.ToolCall{
		Name:      "create_student",
		Arguments: `{"name": "Nuevo", "classroom_id": 9001, "confirmed": true}`,
	})

	// Assert
	require.NoError(t, err)
	assert.Contains(t, result, `"saved":true`)
	students.AssertExpectations(t)
}

func TestInclusionDispatcher_CreateStudentWithoutConfirmedDoesNotPersist(t *testing.T) {
	// Arrange
	students := new(mockproviders.MockStudentProvider)
	d := inclusionDispatcher{students: students}

	// Act
	result, err := d.Dispatch(context.Background(), uuid.New(), providers.ToolCall{
		Name:      "create_student",
		Arguments: `{"name": "Nuevo", "classroom_id": 9001}`,
	})

	// Assert
	require.NoError(t, err)
	assert.Contains(t, result, `"pending_confirmation":true`)
	students.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestInclusionDispatcher_RelateStudentRecursoConfirmedUpdates(t *testing.T) {
	// Arrange
	ctx := context.Background()
	orgID := uuid.New()
	adaptations := new(mockproviders.MockAdaptationProvider)
	adaptations.On("Get", ctx, orgID, int64(5)).Return(&entities.Adaptation{ID: 5, StudentID: 1}, nil)
	adaptations.On("Update", ctx, mock.MatchedBy(func(a *entities.Adaptation) bool {
		return a.ID == 5 && a.StudentID == 9002
	})).Return(nil)
	d := inclusionDispatcher{adaptations: adaptations}

	// Act
	result, err := d.Dispatch(ctx, orgID, providers.ToolCall{
		Name:      "relate_student_recurso",
		Arguments: `{"student_id": 9002, "recurso_id": 5, "confirmed": true}`,
	})

	// Assert
	require.NoError(t, err)
	assert.Contains(t, result, `"saved":true`)
	adaptations.AssertExpectations(t)
}

func TestInclusionTools_ExposeActionTools(t *testing.T) {
	// Arrange / Act
	names := make(map[string]bool)
	for _, tool := range inclusionTools() {
		names[tool.Name] = true
	}

	// Assert
	assert.True(t, names["create_student"])
	assert.True(t, names["create_recurso"])
	assert.True(t, names["relate_student_recurso"])
}
