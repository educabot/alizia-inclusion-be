package inclusion

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/lib/pq"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

// SummarizeConversations resume en lote las conversaciones "cerradas" (inactivas)
// que no tienen resumen actualizado, para que el chat tenga memoria entre clases.
// Lo dispara el cron (cmd/summarizer); no hay endpoint HTTP.
type SummarizeConversations interface {
	Execute(ctx context.Context) (SummarizeResult, error)
}

type SummarizeResult struct {
	Processed int
	Failed    int
}

type summarizeConversationsImpl struct {
	conversations providers.ConversationProvider
	summaries     providers.ConversationSummaryProvider
	ai            providers.AIClient
	usage         providers.AIUsageProvider
	idleMinutes   int
	batchLimit    int
}

func NewSummarizeConversations(
	conversations providers.ConversationProvider,
	summaries providers.ConversationSummaryProvider,
	ai providers.AIClient,
	usage providers.AIUsageProvider,
	idleMinutes, batchLimit int,
) SummarizeConversations {
	return &summarizeConversationsImpl{
		conversations: conversations,
		summaries:     summaries,
		ai:            ai,
		usage:         usage,
		idleMinutes:   idleMinutes,
		batchLimit:    batchLimit,
	}
}

func (uc *summarizeConversationsImpl) Execute(ctx context.Context) (SummarizeResult, error) {
	idleBefore := time.Now().Add(-time.Duration(uc.idleMinutes) * time.Minute)
	pending, err := uc.conversations.ListPendingSummary(ctx, idleBefore, uc.batchLimit)
	if err != nil {
		return SummarizeResult{}, fmt.Errorf("list pending: %w", err)
	}

	var res SummarizeResult
	for i := range pending {
		conv := &pending[i]
		if err := uc.summarizeOne(ctx, conv); err != nil {
			slog.WarnContext(ctx, "summarize conversation failed", "error", err, "conversation_id", conv.ID)
			res.Failed++
			continue
		}
		res.Processed++
	}
	return res, nil
}

func (uc *summarizeConversationsImpl) summarizeOne(ctx context.Context, conv *entities.Conversation) error {
	if len(conv.Messages) == 0 {
		return nil
	}

	messages := []providers.ChatMessage{
		{Role: "system", Content: buildSummaryPrompt()},
		{Role: "user", Content: renderConversationForSummary(conv.Messages)},
	}
	resp, err := uc.ai.Chat(ctx, messages)
	if err != nil {
		return fmt.Errorf("ai chat: %w", err)
	}

	parsed, err := parseSummaryJSON(resp.Content)
	if err != nil {
		return fmt.Errorf("parse summary: %w", err)
	}
	if strings.TrimSpace(parsed.Summary) == "" {
		return fmt.Errorf("empty summary")
	}

	recordAIUsage(ctx, uc.usage, conv.OrganizationID, conv.UserID, "summary", resp.Usage)

	studentIDs, deviceIDs := collectEntityIDs(conv)
	tokenCount := 0
	if resp.Usage != nil {
		tokenCount = resp.Usage.TotalTokens
	}

	return uc.summaries.Upsert(ctx, entities.ConversationSummary{
		ConversationID: conv.ID,
		Summary:        strings.TrimSpace(parsed.Summary),
		TopicKeywords:  pq.StringArray(parsed.TopicKeywords),
		TokenCount:     tokenCount,
		UpdatedAt:      time.Now(),
	}, studentIDs, deviceIDs)
}

// renderConversationForSummary serializa el diálogo en texto plano "rol: contenido".
func renderConversationForSummary(msgs []entities.ConversationMessage) string {
	var b strings.Builder
	for i := range msgs {
		m := &msgs[i]
		fmt.Fprintf(&b, "%s: %s\n", m.Role, m.Content)
	}
	return b.String()
}

type summaryPayload struct {
	Summary       string   `json:"summary"`
	TopicKeywords []string `json:"topic_keywords"`
}

// parseSummaryJSON tolera que el modelo envuelva el JSON en fences o prosa:
// extrae del primer '{' al último '}' y deserializa.
func parseSummaryJSON(content string) (summaryPayload, error) {
	var p summaryPayload
	start := strings.Index(content, "{")
	end := strings.LastIndex(content, "}")
	if start == -1 || end == -1 || end < start {
		return p, fmt.Errorf("no json object in response")
	}
	if err := json.Unmarshal([]byte(content[start:end+1]), &p); err != nil {
		return p, fmt.Errorf("unmarshal: %w", err)
	}
	return p, nil
}

// collectEntityIDs reúne los alumnos/dispositivos involucrados en la conversación:
// el alumno de apertura (conversations.student_id) + los identified_student /
// recommended_device guardados en el metadata de cada mensaje. Dedup en memoria.
func collectEntityIDs(conv *entities.Conversation) (studentIDs, deviceIDs []int64) {
	studentSet := map[int64]struct{}{}
	deviceSet := map[int64]struct{}{}

	if conv.StudentID != nil && *conv.StudentID > 0 {
		studentSet[*conv.StudentID] = struct{}{}
	}

	for i := range conv.Messages {
		meta := map[string]any{}
		raw := conv.Messages[i].Metadata
		if len(raw) == 0 {
			continue
		}
		if err := json.Unmarshal(raw, &meta); err != nil {
			continue
		}
		if id := metaInt64(meta, "identified_student"); id > 0 {
			studentSet[id] = struct{}{}
		}
		if id := metaInt64(meta, "recommended_device"); id > 0 {
			deviceSet[id] = struct{}{}
		}
	}

	for id := range studentSet {
		studentIDs = append(studentIDs, id)
	}
	for id := range deviceSet {
		deviceIDs = append(deviceIDs, id)
	}
	return studentIDs, deviceIDs
}

// metaInt64 lee una clave numérica del metadata jsonb (los números llegan como float64).
func metaInt64(meta map[string]any, key string) int64 {
	v, ok := meta[key]
	if !ok {
		return 0
	}
	switch n := v.(type) {
	case float64:
		return int64(n)
	case int64:
		return n
	default:
		return 0
	}
}
