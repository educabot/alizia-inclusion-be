package inclusion

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	"github.com/educabot/alizia-inclusion-be/src/core/providers/mocks"
)

func TestRunAgenticChat_ReturnsPlainChatResponseWhenNoToolsProvided(t *testing.T) {
	ctx := context.Background()
	orgID := uuid.New()

	chatCalled := false
	ai := &mocks.MockAIClient{
		ChatFn: func(_ context.Context, _ []providers.ChatMessage) (*providers.ChatResponse, error) {
			chatCalled = true
			return &providers.ChatResponse{Content: "hola"}, nil
		},
		ChatWithToolsFn: func(_ context.Context, _ []providers.ChatMessage, _ []providers.ToolDefinition) (*providers.ChatResponse, error) {
			t.Fatal("ChatWithTools must not be called when there are no tools")
			return nil, nil
		},
	}
	msgs := []providers.ChatMessage{{Role: "user", Content: "hola"}}

	resp, err := runAgenticChat(ctx, ai, msgs, nil, inclusionDispatcher{}, orgID, maxAgenticIterations)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !chatCalled {
		t.Error("expected Chat to be called")
	}
	if resp.Content != "hola" {
		t.Errorf("expected content %q, got %q", "hola", resp.Content)
	}
}

func TestRunAgenticChat_ExecutesToolThenReturnsTheFinalAnswer(t *testing.T) {
	ctx := context.Background()
	orgID := uuid.New()

	var turns int
	var toolMessages []providers.ChatMessage
	ai := &mocks.MockAIClient{
		ChatWithToolsFn: func(_ context.Context, msgs []providers.ChatMessage, _ []providers.ToolDefinition) (*providers.ChatResponse, error) {
			turns++
			if turns == 1 {
				return &providers.ChatResponse{
					ToolCalls: []providers.ToolCall{
						{ID: "call_1", Name: "list_devices", Arguments: "{}"},
					},
					Usage: &providers.TokenUsage{TotalTokens: 10},
				}, nil
			}
			toolMessages = msgs
			return &providers.ChatResponse{
				Content: "Te recomiendo el Timer Visual",
				Usage:   &providers.TokenUsage{TotalTokens: 5},
			}, nil
		},
	}
	devices := &mocks.MockDeviceProvider{
		ListDevicesFn: func(_ context.Context, _ uuid.UUID, _ *int64) ([]entities.Device, error) {
			return []entities.Device{{ID: 1, Name: "Timer Visual"}}, nil
		},
	}
	dispatcher := inclusionDispatcher{devices: devices}
	msgs := []providers.ChatMessage{{Role: "user", Content: "que dispositivo uso?"}}

	resp, err := runAgenticChat(ctx, ai, msgs, inclusionTools(), dispatcher, orgID, maxAgenticIterations)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if turns != 2 {
		t.Errorf("expected 2 model turns, got %d", turns)
	}
	if resp.Content != "Te recomiendo el Timer Visual" {
		t.Errorf("unexpected final content: %q", resp.Content)
	}
	if resp.Usage == nil || resp.Usage.TotalTokens != 15 {
		t.Errorf("expected accumulated usage 15, got %+v", resp.Usage)
	}
	var foundToolMsg bool
	for _, m := range toolMessages {
		if m.Role == "tool" && m.ToolCallID == "call_1" && strings.Contains(m.Content, "Timer Visual") {
			foundToolMsg = true
		}
	}
	if !foundToolMsg {
		t.Error("expected a tool result message wired back into the conversation")
	}
}

func TestRunAgenticChat_FeedsErrorResultBackWhenToolFails(t *testing.T) {
	ctx := context.Background()
	orgID := uuid.New()

	var secondTurnMsgs []providers.ChatMessage
	var turns int
	ai := &mocks.MockAIClient{
		ChatWithToolsFn: func(_ context.Context, msgs []providers.ChatMessage, _ []providers.ToolDefinition) (*providers.ChatResponse, error) {
			turns++
			if turns == 1 {
				return &providers.ChatResponse{
					ToolCalls: []providers.ToolCall{
						{ID: "call_x", Name: "list_devices", Arguments: "{}"},
					},
				}, nil
			}
			secondTurnMsgs = msgs
			return &providers.ChatResponse{Content: "lo siento"}, nil
		},
	}
	devices := &mocks.MockDeviceProvider{
		ListDevicesFn: func(_ context.Context, _ uuid.UUID, _ *int64) ([]entities.Device, error) {
			return nil, errors.New("db down")
		},
	}
	dispatcher := inclusionDispatcher{devices: devices}
	msgs := []providers.ChatMessage{{Role: "user", Content: "dispositivos?"}}

	resp, err := runAgenticChat(ctx, ai, msgs, inclusionTools(), dispatcher, orgID, maxAgenticIterations)

	if err != nil {
		t.Fatalf("loop must not abort on tool error: %v", err)
	}
	if resp.Content != "lo siento" {
		t.Errorf("unexpected content: %q", resp.Content)
	}
	var sawErrResult bool
	for _, m := range secondTurnMsgs {
		if m.Role == "tool" && strings.Contains(m.Content, "db down") {
			sawErrResult = true
		}
	}
	if !sawErrResult {
		t.Error("expected the tool error to be fed back to the model")
	}
}

func TestRunAgenticChat_ForcesFinalAnswerWhenIterationBudgetExhausted(t *testing.T) {
	ctx := context.Background()
	orgID := uuid.New()

	toolCalls := 0
	chatFallbackCalled := false
	ai := &mocks.MockAIClient{
		ChatWithToolsFn: func(_ context.Context, _ []providers.ChatMessage, _ []providers.ToolDefinition) (*providers.ChatResponse, error) {
			toolCalls++
			return &providers.ChatResponse{
				ToolCalls: []providers.ToolCall{{ID: "loop", Name: "list_devices", Arguments: "{}"}},
			}, nil
		},
		ChatFn: func(_ context.Context, _ []providers.ChatMessage) (*providers.ChatResponse, error) {
			chatFallbackCalled = true
			return &providers.ChatResponse{Content: "respuesta forzada"}, nil
		},
	}
	devices := &mocks.MockDeviceProvider{
		ListDevicesFn: func(_ context.Context, _ uuid.UUID, _ *int64) ([]entities.Device, error) {
			return nil, nil
		},
	}
	dispatcher := inclusionDispatcher{devices: devices}
	msgs := []providers.ChatMessage{{Role: "user", Content: "loop"}}

	resp, err := runAgenticChat(ctx, ai, msgs, inclusionTools(), dispatcher, orgID, 2)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if toolCalls != 2 {
		t.Errorf("expected 2 tool rounds (the cap), got %d", toolCalls)
	}
	if !chatFallbackCalled {
		t.Error("expected a final plain Chat fallback after the cap")
	}
	if resp.Content != "respuesta forzada" {
		t.Errorf("unexpected content: %q", resp.Content)
	}
}

func TestInclusionDispatcher_GetStudentReturnsStudentPayload(t *testing.T) {
	ctx := context.Background()
	orgID := uuid.New()

	students := &mocks.MockStudentProvider{
		GetStudentFn: func(_ context.Context, _ uuid.UUID, id int64) (*entities.Student, error) {
			return &entities.Student{ID: id, Name: "Lucas"}, nil
		},
	}
	d := inclusionDispatcher{students: students}

	result, err := d.Dispatch(ctx, orgID, providers.ToolCall{
		Name:      "get_student",
		Arguments: `{"student_id": 7}`,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var got entities.Student
	if jerr := json.Unmarshal([]byte(result), &got); jerr != nil {
		t.Fatalf("result is not valid JSON: %v", jerr)
	}
	if got.ID != 7 || got.Name != "Lucas" {
		t.Errorf("unexpected student payload: %+v", got)
	}
}

func TestInclusionDispatcher_RejectsUnknownTool(t *testing.T) {
	ctx := context.Background()
	orgID := uuid.New()

	d := inclusionDispatcher{}

	_, err := d.Dispatch(ctx, orgID, providers.ToolCall{Name: "delete_everything"})

	if err == nil {
		t.Fatal("expected error for unknown tool")
	}
	if !strings.Contains(err.Error(), "unknown tool") {
		t.Errorf("expected unknown-tool error, got: %v", err)
	}
}

func TestInclusionDispatcher_RejectsMalformedArguments(t *testing.T) {
	ctx := context.Background()
	orgID := uuid.New()

	d := inclusionDispatcher{}

	_, err := d.Dispatch(ctx, orgID, providers.ToolCall{
		Name:      "get_student",
		Arguments: `{not json`,
	})

	if err == nil {
		t.Fatal("expected error for malformed arguments")
	}
}
