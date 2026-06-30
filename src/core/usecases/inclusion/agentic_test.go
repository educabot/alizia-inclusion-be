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
	mockproviders "github.com/educabot/alizia-inclusion-be/src/core/providers/mocks"
)

func TestRunAgenticChat_ReturnsPlainChatResponseWhenNoToolsProvided(t *testing.T) {
	ctx := context.Background()
	orgID := uuid.New()
	ai := new(mockproviders.MockAIClient)
	ai.On("Chat", ctx, mock.AnythingOfType("[]providers.ChatMessage")).
		Return(&providers.ChatResponse{Content: "hola"}, nil)
	msgs := []providers.ChatMessage{{Role: "user", Content: "hola"}}

	resp, _, err := runAgenticChat(ctx, ai, msgs, nil, inclusionDispatcher{}, orgID, maxAgenticIterations, false)

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

	resp, _, err := runAgenticChat(ctx, ai, msgs, inclusionTools(), dispatcher, orgID, maxAgenticIterations, false)

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

	resp, _, err := runAgenticChat(ctx, ai, msgs, inclusionTools(), dispatcher, orgID, maxAgenticIterations, false)

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

	resp, _, err := runAgenticChat(ctx, ai, msgs, inclusionTools(), dispatcher, orgID, 2, false)

	require.NoError(t, err)
	assert.Equal(t, "respuesta forzada", resp.Content)
	ai.AssertExpectations(t)
}

// Red de seguridad RAG: si el modelo propone (bloque [STEPS]/[ADAPTATION_JSON]) sin haber
// llamado nunca a search_content_hibrido, el loop inyecta un mensaje que lo obliga a buscar
// y vuelve a consultar. requireSearchBeforeProposal=true (como en AssistClassroom).
func TestRunAgenticChat_ForcesSearchWhenProposalWithoutRAG(t *testing.T) {
	ctx := context.Background()
	orgID := uuid.New()
	var secondTurnMsgs []providers.ChatMessage

	ai := new(mockproviders.MockAIClient)
	// 1) El modelo cierra con una propuesta sin haber buscado.
	ai.On("ChatWithTools", ctx, mock.AnythingOfType("[]providers.ChatMessage"), mock.AnythingOfType("[]providers.ToolDefinition")).
		Return(&providers.ChatResponse{Content: "Acá van los pasos:\n[STEPS]\n1. Anticipá la consigna.\n[/STEPS]"}, nil).Once()
	// 2) Tras el empujón de la red de seguridad, vuelve a consultarse (capturamos los mensajes).
	ai.On("ChatWithTools", ctx, mock.AnythingOfType("[]providers.ChatMessage"), mock.AnythingOfType("[]providers.ToolDefinition")).
		Run(func(args mock.Arguments) {
			msgs, ok := args.Get(1).([]providers.ChatMessage)
			require.True(t, ok)
			secondTurnMsgs = msgs
		}).
		Return(&providers.ChatResponse{Content: "ok"}, nil).Once()

	msgs := []providers.ChatMessage{{Role: "user", Content: "se mueve mucho en clase"}}

	resp, _, err := runAgenticChat(ctx, ai, msgs, inclusionTools(), inclusionDispatcher{}, orgID, maxAgenticIterations, true)

	require.NoError(t, err)
	assert.Equal(t, "ok", resp.Content)
	var sawForce bool
	for _, m := range secondTurnMsgs {
		if m.Role == "user" && strings.Contains(m.Content, "search_content_hibrido") {
			sawForce = true
		}
	}
	assert.True(t, sawForce, "la red de seguridad debía inyectar un mensaje que fuerce la búsqueda")
	ai.AssertExpectations(t)
}

// Sin requireSearchBeforeProposal, una propuesta sin búsqueda se devuelve tal cual (no se
// fuerza nada): el mecanismo es opt-in por usecase.
func TestRunAgenticChat_DoesNotForceSearchWhenFlagOff(t *testing.T) {
	ctx := context.Background()
	orgID := uuid.New()

	ai := new(mockproviders.MockAIClient)
	ai.On("ChatWithTools", ctx, mock.AnythingOfType("[]providers.ChatMessage"), mock.AnythingOfType("[]providers.ToolDefinition")).
		Return(&providers.ChatResponse{Content: "[STEPS]\n1. x\n[/STEPS]"}, nil).Once()

	msgs := []providers.ChatMessage{{Role: "user", Content: "dame pasos"}}

	resp, _, err := runAgenticChat(ctx, ai, msgs, inclusionTools(), inclusionDispatcher{}, orgID, maxAgenticIterations, false)

	require.NoError(t, err)
	assert.Contains(t, resp.Content, "[STEPS]")
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

func TestInclusionDispatcher_CreateStudentCreatesInClassroomWithProfile(t *testing.T) {
	ctx := context.Background()
	orgID := uuid.New()
	students := new(mockproviders.MockStudentProvider)
	profiles := new(mockproviders.MockStudentProfileProvider)

	// El aula viene en los argumentos del modelo (classroom_id). Sin alumnos previos
	// en el aula → no es duplicado → da de alta.
	students.On("ListByClassroom", ctx, orgID, int64(9)).Return([]entities.Student{}, nil)
	students.On("Create", ctx, mock.AnythingOfType("*entities.Student")).
		Return(nil).
		Run(func(args mock.Arguments) {
			s, ok := args.Get(1).(*entities.Student)
			require.True(t, ok)
			s.ID = 42
		})
	profiles.On("Upsert", ctx, mock.AnythingOfType("*entities.StudentProfile")).Return(nil)

	d := inclusionDispatcher{students: students, profiles: profiles}

	result, err := d.Dispatch(ctx, orgID, providers.ToolCall{
		Name:      "create_student",
		Arguments: `{"name":"Lucas Pérez","classroom_id":9,"difficulties":["le cuesta sostener la atención"]}`,
	})

	require.NoError(t, err)
	var got entities.Student
	require.NoError(t, json.Unmarshal([]byte(result), &got))
	assert.Equal(t, int64(42), got.ID)
	assert.Equal(t, int64(9), got.ClassroomID)
	assert.Equal(t, "Lucas Pérez", got.Name)
	require.NotNil(t, got.Profile)
	assert.Equal(t, []string{"le cuesta sostener la atención"}, []string(got.Profile.Difficulties))
	students.AssertExpectations(t)
	profiles.AssertExpectations(t)
}

// Idempotencia: si ya existe un alumno con el mismo nombre (normalizado) en el aula,
// create_student devuelve el existente sin volver a crearlo.
func TestInclusionDispatcher_CreateStudentIsIdempotentWithinClassroom(t *testing.T) {
	ctx := context.Background()
	orgID := uuid.New()
	students := new(mockproviders.MockStudentProvider)

	students.On("ListByClassroom", ctx, orgID, int64(9)).
		Return([]entities.Student{{ID: 7, Name: "Lucas Pérez", ClassroomID: 9}}, nil)

	d := inclusionDispatcher{students: students}

	result, err := d.Dispatch(ctx, orgID, providers.ToolCall{
		Name:      "create_student",
		Arguments: `{"name":"  lucas perez  ","classroom_id":9}`,
	})

	require.NoError(t, err)
	var got entities.Student
	require.NoError(t, json.Unmarshal([]byte(result), &got))
	assert.Equal(t, int64(7), got.ID)
	students.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestInclusionDispatcher_CreateStudentRequiresClassroom(t *testing.T) {
	d := inclusionDispatcher{}

	_, err := d.Dispatch(context.Background(), uuid.New(), providers.ToolCall{
		Name:      "create_student",
		Arguments: `{"name":"Lucas"}`,
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "aula")
}

func TestInclusionDispatcher_CreateStudentRejectsEmptyName(t *testing.T) {
	students := new(mockproviders.MockStudentProvider)
	d := inclusionDispatcher{students: students}

	_, err := d.Dispatch(context.Background(), uuid.New(), providers.ToolCall{
		Name:      "create_student",
		Arguments: `{"name":"   ","classroom_id":1}`,
	})

	require.Error(t, err)
	students.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestInclusionDispatcher_FindStudentByNameMatchesNormalized(t *testing.T) {
	ctx := context.Background()
	orgID := uuid.New()
	students := new(mockproviders.MockStudentProvider)
	students.On("List", ctx, orgID).Return([]entities.Student{
		{ID: 1, Name: "Lucas Pérez", ClassroomID: 9},
		{ID: 2, Name: "Martina Gómez", ClassroomID: 9},
	}, nil)
	d := inclusionDispatcher{students: students}

	result, err := d.Dispatch(ctx, orgID, providers.ToolCall{
		Name:      "find_student_by_name",
		Arguments: `{"name":"lucas"}`,
	})

	require.NoError(t, err)
	var got struct {
		Students []struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		} `json:"students"`
	}
	require.NoError(t, json.Unmarshal([]byte(result), &got))
	require.Len(t, got.Students, 1)
	assert.Equal(t, int64(1), got.Students[0].ID)
}

func TestInclusionDispatcher_ListClassrooms(t *testing.T) {
	ctx := context.Background()
	orgID := uuid.New()
	grade := "3ro"
	classrooms := new(mockproviders.MockClassroomProvider)
	classrooms.On("List", ctx, orgID).Return([]entities.Classroom{
		{ID: 5, Name: "3ro A", Grade: &grade},
	}, nil)
	d := inclusionDispatcher{classrooms: classrooms}

	result, err := d.Dispatch(ctx, orgID, providers.ToolCall{Name: "list_classrooms", Arguments: `{}`})

	require.NoError(t, err)
	assert.Contains(t, result, "3ro A")
}

func TestInclusionDispatcher_CreateClassroomNormalizesGrade(t *testing.T) {
	ctx := context.Background()
	orgID := uuid.New()
	classrooms := new(mockproviders.MockClassroomProvider)
	// No existe aún → la crea.
	classrooms.On("List", ctx, orgID).Return([]entities.Classroom{}, nil)
	classrooms.On("Create", ctx, mock.AnythingOfType("*entities.Classroom")).
		Return(nil).
		Run(func(args mock.Arguments) {
			c, ok := args.Get(1).(*entities.Classroom)
			require.True(t, ok)
			c.ID = 11
			assert.Equal(t, "3ro B", c.Name)
			require.NotNil(t, c.Grade)
			assert.Equal(t, "3ro", *c.Grade)
			require.NotNil(t, c.Section)
			assert.Equal(t, "B", *c.Section)
		})
	d := inclusionDispatcher{classrooms: classrooms}

	result, err := d.Dispatch(ctx, orgID, providers.ToolCall{
		Name:      "create_classroom",
		Arguments: `{"grade":"tercero B"}`,
	})

	require.NoError(t, err)
	assert.Contains(t, result, "3ro B")
	classrooms.AssertExpectations(t)
}

func TestInclusionDispatcher_CreateClassroomIsIdempotent(t *testing.T) {
	ctx := context.Background()
	orgID := uuid.New()
	classrooms := new(mockproviders.MockClassroomProvider)
	classrooms.On("List", ctx, orgID).Return([]entities.Classroom{{ID: 3, Name: "3ro A"}}, nil)
	d := inclusionDispatcher{classrooms: classrooms}

	result, err := d.Dispatch(ctx, orgID, providers.ToolCall{
		Name:      "create_classroom",
		Arguments: `{"grade":"3ro A"}`,
	})

	require.NoError(t, err)
	assert.Contains(t, result, `"id":3`)
	classrooms.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestInclusionTools_ExposeCreateStudent(t *testing.T) {
	var found bool
	for _, tool := range inclusionTools() {
		if tool.Name != "create_student" {
			continue
		}
		found = true
		params, ok := tool.Parameters.(map[string]any)
		require.True(t, ok)
		req, _ := params["required"].([]string)
		assert.Contains(t, req, "name")
	}
	assert.True(t, found, "create_student debe estar expuesta como tool")
}

func TestInclusionDispatcher_RejectsUnknownTool(t *testing.T) {
	d := inclusionDispatcher{}

	_, err := d.Dispatch(context.Background(), uuid.New(), providers.ToolCall{Name: "delete_everything"})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown tool")
}

func TestInclusionDispatcher_RejectsMalformedArguments(t *testing.T) {
	d := inclusionDispatcher{}

	_, err := d.Dispatch(context.Background(), uuid.New(), providers.ToolCall{
		Name:      "get_student",
		Arguments: `{not json`,
	})

	require.Error(t, err)
}
