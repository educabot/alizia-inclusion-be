package inclusion

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	mockproviders "github.com/educabot/alizia-inclusion-be/src/core/providers/mocks"
)

func TestInclusionDispatcher_GetStudentHistoryReturnsSummaries(t *testing.T) {
	// Arrange
	ctx := context.Background()
	orgID := uuid.New()
	summaries := new(mockproviders.MockConversationSummaryProvider)
	summaries.On("RecentByStudent", ctx, orgID, int64(7), maxPriorSummaries).
		Return([]entities.ConversationSummary{
			{ConversationID: 99, Summary: "Veníamos trabajando lectura", TopicKeywords: []string{"lectura"}},
		}, nil)
	d := inclusionDispatcher{summaries: summaries}

	// Act
	result, err := d.Dispatch(ctx, orgID, providers.ToolCall{
		Name:      "get_student_history",
		Arguments: `{"student_id": 7}`,
	})

	// Assert
	require.NoError(t, err)
	var got struct {
		History []struct {
			ConversationID int64    `json:"conversation_id"`
			Summary        string   `json:"summary"`
			TopicKeywords  []string `json:"topic_keywords"`
		} `json:"history"`
	}
	require.NoError(t, json.Unmarshal([]byte(result), &got))
	require.Len(t, got.History, 1)
	assert.Equal(t, int64(99), got.History[0].ConversationID)
	assert.Equal(t, "Veníamos trabajando lectura", got.History[0].Summary)
	summaries.AssertExpectations(t)
}

func TestInclusionDispatcher_GetStudentHistoryUnavailableWhenProviderNil(t *testing.T) {
	// Arrange
	d := inclusionDispatcher{}

	// Act
	_, err := d.Dispatch(context.Background(), uuid.New(), providers.ToolCall{
		Name:      "get_student_history",
		Arguments: `{"student_id": 7}`,
	})

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no disponible")
}

func TestInclusionDispatcher_GetPastAdaptationsReturnsStatusAndOutcome(t *testing.T) {
	// Arrange
	ctx := context.Background()
	orgID := uuid.New()
	outcome := "funcionó muy bien con time timer"
	adaptations := new(mockproviders.MockAdaptationProvider)
	studentID := int64(7)
	adaptations.On("List", ctx, providers.AdaptationFilter{OrgID: orgID, StudentID: &studentID}).
		Return([]entities.Adaptation{
			{ID: 1, Subject: "Lengua", Status: "funcionó", Outcome: &outcome},
		}, nil)
	d := inclusionDispatcher{adaptations: adaptations}

	// Act
	result, err := d.Dispatch(ctx, orgID, providers.ToolCall{
		Name:      "get_past_adaptations",
		Arguments: `{"student_id": 7}`,
	})

	// Assert
	require.NoError(t, err)
	var got struct {
		Adaptations []struct {
			ID      int64  `json:"id"`
			Subject string `json:"subject"`
			Status  string `json:"status"`
			Outcome string `json:"outcome"`
		} `json:"adaptations"`
	}
	require.NoError(t, json.Unmarshal([]byte(result), &got))
	require.Len(t, got.Adaptations, 1)
	assert.Equal(t, "funcionó", got.Adaptations[0].Status)
	assert.Equal(t, outcome, got.Adaptations[0].Outcome)
	adaptations.AssertExpectations(t)
}

func TestInclusionTools_ExposeHistoryAndPastAdaptations(t *testing.T) {
	// Arrange / Act
	names := make(map[string]bool)
	for _, tool := range inclusionTools() {
		names[tool.Name] = true
	}

	// Assert
	assert.True(t, names["get_student_history"])
	assert.True(t, names["get_past_adaptations"])
}
