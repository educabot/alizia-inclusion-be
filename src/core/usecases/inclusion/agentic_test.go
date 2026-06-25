package inclusion

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	mockproviders "github.com/educabot/alizia-inclusion-be/src/mocks/providers"
)

func TestRunAgenticChat_ReturnsPlainChatResponseWhenNoToolsProvided(t *testing.T) {
	ctx := context.Background()
	orgID := uuid.New()
	ai := new(mockproviders.MockAIClient)
	ai.On("Chat", ctx, mock.AnythingOfType("[]providers.ChatMessage")).
		Return(&providers.ChatResponse{Content: "hola"}, nil)
	msgs := []providers.ChatMessage{{Role: "user", Content: "hola"}}

	resp, _, err := runAgenticChat(ctx, ai, msgs, nil, inclusionDispatcher{}, orgID, maxAgenticIterations)

	require.NoError(t, err)
	assert.Equal(t, "hola", resp.Content)
	ai.AssertExpectations(t)
	ai.AssertNotCalled(t, "ChatWithTools", mock.Anything, mock.Anything, mock.Anything)
}

func TestRunAgenticChat_ExecutesToolThenReturnsTheFinalAnswer(t *testing.T) {
	ctx := context.Background()
	orgID := uuid.New()
	var secondTurnMsgs []providers.ChatMessage

	ai := new(mockproviders.MockAIClient)
	ai.On("ChatWithTools", ctx, mock.AnythingOfType("[]providers.ChatMessage"), mock.AnythingOfType("[]providers.ToolDefinition")).
		Return(&providers.ChatResponse{
			ToolCalls: []providers.ToolCall{{ID: "call_1", Name: "list_devices", Arguments: "{}"}},
			Usage:     &providers.TokenUsage{TotalTokens: 10},
		}, nil).Once()
	ai.On("ChatWithTools", ctx, mock.AnythingOfType("[]providers.ChatMessage"), mock.AnythingOfType("[]providers.ToolDefinition")).
		Run(func(args mock.Arguments) {
			msgs, ok := args.Get(1).([]providers.ChatMessage)
			require.True(t, ok)
			secondTurnMsgs = msgs
		}).
		Return(&providers.ChatResponse{
			Content: "Te recomiendo el Timer Visual",
			Usage:   &providers.TokenUsage{TotalTokens: 5},
		}, nil).Once()

	devices := new(mockproviders.MockDeviceProvider)
	devices.On("ListDevices", ctx, orgID, (*int64)(nil)).
		Return([]entities.Device{{ID: 1, Name: "Timer Visual"}}, nil)
	dispatcher := inclusionDispatcher{devices: devices}
	msgs := []providers.ChatMessage{{Role: "user", Content: "que dispositivo uso?"}}

	resp, _, err := runAgenticChat(ctx, ai, msgs, inclusionTools(), dispatcher, orgID, maxAgenticIterations)

	require.NoError(t, err)
	assert.Equal(t, "Te recomiendo el Timer Visual", resp.Content)
	require.NotNil(t, resp.Usage)
	assert.Equal(t, 15, resp.Usage.TotalTokens)
	var foundToolMsg bool
	for _, m := range secondTurnMsgs {
		if m.Role == "tool" && m.ToolCallID == "call_1" && strings.Contains(m.Content, "Timer Visual") {
			foundToolMsg = true
		}
	}
	assert.True(t, foundToolMsg, "expected a tool result message wired back into the conversation")
	ai.AssertExpectations(t)
	devices.AssertExpectations(t)
}

func TestRunAgenticChat_FeedsErrorResultBackWhenToolFails(t *testing.T) {
	ctx := context.Background()
	orgID := uuid.New()
	var secondTurnMsgs []providers.ChatMessage

	ai := new(mockproviders.MockAIClient)
	ai.On("ChatWithTools", ctx, mock.AnythingOfType("[]providers.ChatMessage"), mock.AnythingOfType("[]providers.ToolDefinition")).
		Return(&providers.ChatResponse{
			ToolCalls: []providers.ToolCall{{ID: "call_x", Name: "list_devices", Arguments: "{}"}},
		}, nil).Once()
	ai.On("ChatWithTools", ctx, mock.AnythingOfType("[]providers.ChatMessage"), mock.AnythingOfType("[]providers.ToolDefinition")).
		Run(func(args mock.Arguments) {
			msgs, ok := args.Get(1).([]providers.ChatMessage)
			require.True(t, ok)
			secondTurnMsgs = msgs
		}).
		Return(&providers.ChatResponse{Content: "lo siento"}, nil).Once()

	devices := new(mockproviders.MockDeviceProvider)
	devices.On("ListDevices", ctx, orgID, (*int64)(nil)).Return(nil, errors.New("db down"))
	dispatcher := inclusionDispatcher{devices: devices}
	msgs := []providers.ChatMessage{{Role: "user", Content: "dispositivos?"}}

	resp, _, err := runAgenticChat(ctx, ai, msgs, inclusionTools(), dispatcher, orgID, maxAgenticIterations)

	require.NoError(t, err, "loop must not abort on tool error")
	assert.Equal(t, "lo siento", resp.Content)
	var sawErrResult bool
	for _, m := range secondTurnMsgs {
		if m.Role == "tool" && strings.Contains(m.Content, "db down") {
			sawErrResult = true
		}
	}
	assert.True(t, sawErrResult, "expected the tool error to be fed back to the model")
}

func TestRunAgenticChat_ForcesFinalAnswerWhenIterationBudgetExhausted(t *testing.T) {
	ctx := context.Background()
	orgID := uuid.New()

	ai := new(mockproviders.MockAIClient)
	ai.On("ChatWithTools", ctx, mock.AnythingOfType("[]providers.ChatMessage"), mock.AnythingOfType("[]providers.ToolDefinition")).
		Return(&providers.ChatResponse{
			ToolCalls: []providers.ToolCall{{ID: "loop", Name: "list_devices", Arguments: "{}"}},
		}, nil).Times(2)
	ai.On("Chat", ctx, mock.AnythingOfType("[]providers.ChatMessage")).
		Return(&providers.ChatResponse{Content: "respuesta forzada"}, nil).Once()

	devices := new(mockproviders.MockDeviceProvider)
	devices.On("ListDevices", ctx, orgID, (*int64)(nil)).Return(nil, nil)
	dispatcher := inclusionDispatcher{devices: devices}
	msgs := []providers.ChatMessage{{Role: "user", Content: "loop"}}

	resp, _, err := runAgenticChat(ctx, ai, msgs, inclusionTools(), dispatcher, orgID, 2)

	require.NoError(t, err)
	assert.Equal(t, "respuesta forzada", resp.Content)
	ai.AssertExpectations(t)
}

func TestInclusionDispatcher_GetStudentReturnsStudentPayload(t *testing.T) {
	ctx := context.Background()
	orgID := uuid.New()
	students := new(mockproviders.MockStudentProvider)
	students.On("GetStudent", ctx, orgID, int64(7)).
		Return(&entities.Student{ID: 7, Name: "Lucas"}, nil)
	d := inclusionDispatcher{students: students}

	result, err := d.Dispatch(ctx, orgID, providers.ToolCall{
		Name:      "get_student",
		Arguments: `{"student_id": 7}`,
	})

	require.NoError(t, err)
	var got entities.Student
	require.NoError(t, json.Unmarshal([]byte(result), &got))
	assert.Equal(t, int64(7), got.ID)
	assert.Equal(t, "Lucas", got.Name)
	students.AssertExpectations(t)
}

func TestInclusionDispatcher_RejectsUnknownTool(t *testing.T) {
	d := inclusionDispatcher{}

	_, err := d.Dispatch(context.Background(), uuid.New(), providers.ToolCall{Name: "delete_everything"})

	require.ErrorIs(t, err, errUnknownTool)
}

func TestInclusionDispatcher_RejectsMalformedArguments(t *testing.T) {
	d := inclusionDispatcher{}

	_, err := d.Dispatch(context.Background(), uuid.New(), providers.ToolCall{
		Name:      "get_student",
		Arguments: `{not json`,
	})

	require.Error(t, err)
}
