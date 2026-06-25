package inclusion

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/inclusion/prompts"
)

// maxSummaryInputTokens caps how much history is sent to the model for compaction.
// When the conversation exceeds this limit, the most recent turns are preserved
// so the teacher can resume the thread without re-establishing full context.
const maxSummaryInputTokens = 6000

// maxTopicKeywords caps the number of topic-tag keywords stored per summary.
const maxTopicKeywords = 8

// CloseSessionRequest triggers compaction of a conversation when it is closed.
type CloseSessionRequest struct {
	OrgID          uuid.UUID
	UserID         int64
	ConversationID int64
}

func (r CloseSessionRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return errOrgIDRequired
	}
	if r.ConversationID <= 0 {
		return errConversationIDRequired
	}
	return nil
}

// CloseSessionResponse returns the compacted summary and its tags across three
// dimensions (student / topic / device) for display on the frontend.
type CloseSessionResponse struct {
	ConversationID int64    `json:"conversation_id"`
	Summary        string   `json:"summary"`
	TopicKeywords  []string `json:"topic_keywords"`
	StudentIDs     []int64  `json:"student_ids,omitempty"`
	DeviceIDs      []int64  `json:"device_ids,omitempty"`
}

type CloseSession interface {
	Execute(ctx context.Context, req CloseSessionRequest) (*CloseSessionResponse, error)
}

type closeSessionImpl struct {
	ai            providers.AIClient
	conversations providers.ConversationProvider
	summaries     providers.ConversationSummaryProvider
	usage         providers.AIUsageProvider
}

func NewCloseSession(ai providers.AIClient, conversations providers.ConversationProvider, summaries providers.ConversationSummaryProvider, usage providers.AIUsageProvider) CloseSession {
	return &closeSessionImpl{ai: ai, conversations: conversations, summaries: summaries, usage: usage}
}

// Execute compacts a conversation on close: fetches its messages,
// generates a compressed summary via LLM with three-dimension tags (student / topic /
// device), and persists it idempotently by conversation_id.
func (uc *closeSessionImpl) Execute(ctx context.Context, req CloseSessionRequest) (*CloseSessionResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	conv, err := uc.conversations.GetWithMessages(ctx, req.OrgID, req.ConversationID)
	if err != nil {
		return nil, err
	}

	// No messages means nothing to compact: avoids a model call and an empty summary in DB.
	if len(conv.Messages) == 0 {
		return &CloseSessionResponse{ConversationID: conv.ID}, nil
	}

	transcript := buildTranscript(conv.Messages, maxSummaryInputTokens)

	start := time.Now()
	resp, err := uc.ai.Chat(ctx, []providers.ChatMessage{
		{Role: "system", Content: summarizerSystemPrompt},
		{Role: "user", Content: transcript},
	})
	if err != nil {
		return nil, wrapServiceUnavailable(err)
	}

	summaryText, keywords := parseSummary(resp.Content)
	studentIDs, deviceIDs := collectEntities(conv)

	// Per-turn trace: IDs only, no PII. Best-effort.
	recordAIUsage(ctx, uc.usage, aiTrace{
		orgID: req.OrgID, userID: req.UserID, mode: modeClose,
		model: uc.ai.Model(), latencyMs: int(time.Since(start).Milliseconds()),
		conversationID: conv.ID, usage: resp.Usage,
		context: map[string]any{"student_ids": studentIDs, "device_ids": deviceIDs},
	})

	summary := &entities.ConversationSummary{
		ConversationID: conv.ID,
		Summary:        summaryText,
		TopicKeywords:  pq.StringArray(keywords),
		TokenCount:     estimateTokens(summaryText),
	}
	if err := uc.summaries.Upsert(ctx, summary, studentIDs, deviceIDs); err != nil {
		return nil, fmt.Errorf("persist conversation summary: %w", err)
	}

	return &CloseSessionResponse{
		ConversationID: conv.ID,
		Summary:        summaryText,
		TopicKeywords:  keywords,
		StudentIDs:     studentIDs,
		DeviceIDs:      deviceIDs,
	}, nil
}

// summarizerSystemPrompt reuses the shared identity (prompts.RolAlizia) so the role is
// declared in one place, then layers its own JSON-only contract — the summarizer is an
// internal prompt, so it deliberately omits the conversational voice/format rules.
const summarizerSystemPrompt = prompts.RolAlizia + "\n\n" +
	"Resumí la siguiente conversación entre un docente y vos " +
	"para poder retomar el hilo más adelante sin recontextualizar todo. Devolvé EXCLUSIVAMENTE un JSON con esta forma:\n" +
	"{\"summary\": \"un par de párrafos en español con lo trabajado, decisiones y próximos pasos\", " +
	"\"topic_keywords\": [\"palabras clave en minúscula: temas y discapacidades tratadas\"]}\n" +
	"No incluyas nombres propios de alumnos ni diagnósticos en las keywords; usá temas (ej. \"autismo\", \"lectura\", \"autorregulación\"). " +
	"No agregues texto fuera del JSON."

// buildTranscript assembles the conversation text for summarization, preserving
// the most recent turns up to maxTokens (recency matters most for resuming context).
// Iterates backwards then reconstructs in forward order.
func buildTranscript(messages []entities.ConversationMessage, maxTokens int) string {
	kept := make([]string, 0, len(messages))
	used := 0
	for i := len(messages) - 1; i >= 0; i-- {
		line := messages[i].Role + ": " + messages[i].Content
		t := estimateTokens(line)
		if used+t > maxTokens && len(kept) > 0 {
			break
		}
		used += t
		kept = append([]string{line}, kept...)
	}
	return strings.Join(kept, "\n")
}

// parseSummary extracts the summary text and keywords from the JSON returned by the model.
// Fault-tolerant: if the model wraps the JSON in prose or parsing fails, it falls back
// to the raw content as the summary with no keywords, so session close never breaks.
func parseSummary(content string) (summary string, keywords []string) {
	raw := extractJSONObject(content)
	if raw == "" {
		return strings.TrimSpace(content), nil
	}
	var parsed struct {
		Summary       string   `json:"summary"`
		TopicKeywords []string `json:"topic_keywords"`
	}
	if err := json.Unmarshal([]byte(raw), &parsed); err != nil || strings.TrimSpace(parsed.Summary) == "" {
		return strings.TrimSpace(content), nil
	}
	return strings.TrimSpace(parsed.Summary), normalizeKeywords(parsed.TopicKeywords)
}

// extractJSONObject returns the first balanced JSON object found in s (the model
// sometimes wraps it in fences or prose). Returns empty string if none found.
func extractJSONObject(s string) string {
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")
	if start == -1 || end == -1 || end < start {
		return ""
	}
	return s[start : end+1]
}

// normalizeKeywords lowercases, deduplicates, and caps the keyword list.
func normalizeKeywords(in []string) []string {
	seen := make(map[string]struct{}, len(in))
	out := make([]string, 0, len(in))
	for _, k := range in {
		k = strings.ToLower(strings.TrimSpace(k))
		if k == "" {
			continue
		}
		if _, dup := seen[k]; dup {
			continue
		}
		seen[k] = struct{}{}
		out = append(out, k)
		if len(out) >= maxTopicKeywords {
			break
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// collectEntities gathers student and device IDs linked to the conversation for
// tagging the summary by dimension (student / device). Includes the conversation's
// primary student plus any detected per-turn in message metadata.
func collectEntities(conv *entities.Conversation) (studentIDs, deviceIDs []int64) {
	students := newIDSet()
	devices := newIDSet()

	if conv.StudentID != nil {
		students.add(*conv.StudentID)
	}
	for i := range conv.Messages {
		meta := decodeMessageMeta(conv.Messages[i].Metadata)
		if id, ok := metaInt64(meta, metaKeyIdentifiedStudent); ok {
			students.add(id)
		}
		if id, ok := metaInt64(meta, metaKeyRecommendedDevice); ok {
			devices.add(id)
		}
	}
	return students.slice(), devices.slice()
}

func decodeMessageMeta(raw []byte) map[string]any {
	if len(raw) == 0 {
		return nil
	}
	var meta map[string]any
	if err := json.Unmarshal(raw, &meta); err != nil {
		return nil
	}
	return meta
}

// metaInt64 reads a numeric ID from message metadata. JSON decodes numbers as float64.
func metaInt64(meta map[string]any, key string) (int64, bool) {
	v, ok := meta[key]
	if !ok {
		return 0, false
	}
	switch n := v.(type) {
	case float64:
		return int64(n), true
	case int64:
		return n, true
	default:
		return 0, false
	}
}

// idSet preserves insertion order and deduplicates positive IDs.
type idSet struct {
	seen  map[int64]struct{}
	order []int64
}

func newIDSet() *idSet {
	return &idSet{seen: make(map[int64]struct{})}
}

func (s *idSet) add(id int64) {
	if id <= 0 {
		return
	}
	if _, ok := s.seen[id]; ok {
		return
	}
	s.seen[id] = struct{}{}
	s.order = append(s.order, id)
}

func (s *idSet) slice() []int64 {
	if len(s.order) == 0 {
		return nil
	}
	return s.order
}
