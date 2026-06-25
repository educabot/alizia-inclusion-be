package inclusion

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/inclusion/prompts"
)

type AssistClassroomRequest struct {
	OrgID          uuid.UUID
	UserID         int64
	ConversationID int64
	ClassroomID    int64
	StudentID      *int64
	Message        string
	Mode           string
	History        []providers.ChatMessage
}

func (r AssistClassroomRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return errOrgIDRequired
	}
	if r.Message == "" {
		return errMessageRequired
	}
	return nil
}

type AssistClassroomResponse struct {
	Response          string               `json:"response"`
	ConversationID    int64                `json:"conversation_id"`
	IdentifiedStudent *int64               `json:"identified_student,omitempty"`
	RecommendedDevice *int64               `json:"recommended_device,omitempty"`
	Adaptation        *GeneratedAdaptation `json:"adaptation,omitempty"`
	// ReferencedContent lists the pedagogical materials the assistant cited this turn
	// (via [CONTENT_ID:X]), so the frontend can deep-link a chip to the exact document
	// instead of the materials list. Omitted when nothing was cited.
	ReferencedContent []ContentRef `json:"referenced_content,omitempty"`
}

// ContentRef is a lightweight reference to a pedagogical content document the
// assistant cited, carrying just what the frontend needs to render and link a chip.
type ContentRef struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

type AssistClassroom interface {
	Execute(ctx context.Context, req AssistClassroomRequest) (*AssistClassroomResponse, error)
}

type assistClassroomImpl struct {
	ai            providers.AIClient
	students      providers.StudentProvider
	devices       providers.DeviceProvider
	conversations providers.ConversationProvider
	summaries     providers.ConversationSummaryProvider
	adaptations   providers.AdaptationProvider
	content       providers.PedagogicalContentProvider
	usage         providers.AIUsageProvider
	agentic       bool
}

func NewAssistClassroom(ai providers.AIClient, students providers.StudentProvider, devices providers.DeviceProvider, conversations providers.ConversationProvider, summaries providers.ConversationSummaryProvider, adaptations providers.AdaptationProvider, content providers.PedagogicalContentProvider, usage providers.AIUsageProvider, agentic bool) AssistClassroom {
	return &assistClassroomImpl{ai: ai, students: students, devices: devices, conversations: conversations, summaries: summaries, adaptations: adaptations, content: content, usage: usage, agentic: agentic}
}

func (uc *assistClassroomImpl) Execute(ctx context.Context, req AssistClassroomRequest) (*AssistClassroomResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	devices, err := uc.devices.ListDevices(ctx, req.OrgID, nil)
	if err != nil {
		return nil, err
	}

	allStudents, _ := uc.students.ListByClassroom(ctx, req.OrgID, req.ClassroomID)

	// Pedagogical framework is owned by the prompts package.
	systemPrompt := prompts.AssistSystem(devices, allStudents)
	messages := buildChatMessages(systemPrompt, req.History, req.Message)

	var tools []providers.ToolDefinition
	if uc.agentic {
		tools = inclusionTools()
	}
	dispatcher := inclusionDispatcher{students: uc.students, devices: uc.devices, summaries: uc.summaries, adaptations: uc.adaptations, content: uc.content, userID: req.UserID, conversationID: req.ConversationID}

	start := time.Now()
	resp, toolCalls, err := runAgenticChat(ctx, uc.ai, messages, tools, dispatcher, req.OrgID, maxAgenticIterations)
	if err != nil {
		return nil, wrapServiceUnavailable(err)
	}
	latencyMs := int(time.Since(start).Milliseconds())

	// Hard guardrail: if the response references a non-existent DEVICE_ID
	// or crosses another hard limit, suppress it and fall back to the off-ramp.
	guardAnswer(ctx, resp, devices, "usecase", "assist_classroom", "user_id", req.UserID, "mode", req.Mode)

	studentID := extractStudentID(resp.Content)
	deviceID := extractDeviceID(resp.Content)
	adaptation := extractAdaptationJSON(resp.Content)
	referencedContent := uc.resolveReferencedContent(ctx, req.OrgID, extractContentIDs(resp.Content))

	// Never surface a "save resource" card on the opening turn: a pedagogical resource
	// is only offered after the teacher has discussed it and confirmed across turns.
	// This is a hard backend guarantee independent of the prompt's confirmation protocol.
	if len(req.History) == 0 {
		adaptation = nil
	}

	// Carry the identified student into the adaptation so the frontend can build the
	// save request from a single object (it previously had to read the sibling field).
	if adaptation != nil && adaptation.StudentID == nil {
		adaptation.StudentID = studentID
	}

	convID, persistErr := uc.persistTurn(ctx, req, resp.Content, studentID, deviceID, adaptation)
	if persistErr != nil {
		slog.WarnContext(ctx, "assist_classroom: persist turn failed", "error", persistErr, "user_id", req.UserID, "mode", req.Mode)
		convID = req.ConversationID
	}

	// Per-turn trace: IDs only, no PII. Best-effort.
	recordAIUsage(ctx, uc.usage, aiTrace{
		orgID: req.OrgID, userID: req.UserID, mode: modeAssist,
		model: uc.ai.Model(), latencyMs: latencyMs, toolCalls: toolCalls,
		conversationID: convID, usage: resp.Usage,
		context: assistContextSnapshot(req, studentID, deviceID),
	})

	return &AssistClassroomResponse{
		Response:          resp.Content,
		ConversationID:    convID,
		IdentifiedStudent: studentID,
		RecommendedDevice: deviceID,
		Adaptation:        adaptation,
		ReferencedContent: referencedContent,
	}, nil
}

// resolveReferencedContent turns cited content ids into {id, title} refs, scoped to
// the org. Missing or cross-org ids are silently skipped (best-effort): a broken chip
// reference must never fail the chat turn.
func (uc *assistClassroomImpl) resolveReferencedContent(ctx context.Context, orgID uuid.UUID, ids []int64) []ContentRef {
	if len(ids) == 0 || uc.content == nil {
		return nil
	}
	refs := make([]ContentRef, 0, len(ids))
	for _, id := range ids {
		doc, err := uc.content.GetContent(ctx, orgID, id)
		if err != nil || doc == nil {
			continue
		}
		title := ""
		if doc.Title != nil {
			title = *doc.Title
		}
		refs = append(refs, ContentRef{ID: id, Title: title})
	}
	if len(refs) == 0 {
		return nil
	}
	return refs
}

// assistContextSnapshot builds the per-turn context snapshot using IDs only (no PII):
// classroom, student, and recommended device references.
func assistContextSnapshot(req AssistClassroomRequest, studentID, deviceID *int64) map[string]any {
	snap := map[string]any{}
	if req.ClassroomID > 0 {
		snap["classroom_id"] = req.ClassroomID
	}
	if req.StudentID != nil {
		snap["student_id"] = *req.StudentID
	}
	if studentID != nil {
		snap["identified_student_id"] = *studentID
	}
	if deviceID != nil {
		snap["recommended_device_id"] = *deviceID
	}
	return snap
}

func (uc *assistClassroomImpl) persistTurn(ctx context.Context, req AssistClassroomRequest, assistantContent string, studentID, deviceID *int64, adaptation *GeneratedAdaptation) (int64, error) {
	if uc.conversations == nil || req.UserID == 0 {
		return req.ConversationID, nil
	}
	mode := req.Mode
	if mode == "" {
		mode = modeAssist
	}
	metadata := map[string]any{}
	if studentID != nil {
		metadata[metaKeyIdentifiedStudent] = *studentID
	}
	if deviceID != nil {
		metadata[metaKeyRecommendedDevice] = *deviceID
	}
	if adaptation != nil {
		metadata[metaKeyAdaptation] = adaptation
	}
	return uc.conversations.AppendTurn(ctx, providers.AppendTurnParams{
		ConversationID:   req.ConversationID,
		OrgID:            req.OrgID,
		UserID:           req.UserID,
		Mode:             mode,
		StudentID:        req.StudentID,
		UserContent:      req.Message,
		AssistantContent: assistantContent,
		Metadata:         metadata,
	})
}
