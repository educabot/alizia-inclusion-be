package inclusion

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

// maxSummaryInputTokens acota cuánto historial mandamos al modelo para compactar.
// Si la conversación lo excede, preservamos los turnos más recientes (lo que el
// docente necesita para retomar el hilo al volver).
const maxSummaryInputTokens = 6000

// maxTopicKeywords acota las keywords que guardamos como tags de tema (HU-5).
const maxTopicKeywords = 8

// CloseSessionRequest dispara la compactación de una conversación al cerrarla.
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

// CloseSessionResponse devuelve el resumen compactado y sus tags a 3 dimensiones
// (alumno / tema / valija) para que el front pueda mostrarlo.
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

// Execute compacta una conversación al cerrarla (HU-5, §6.4): trae sus mensajes,
// genera un resumen comprimido vía LLM con tags a 3 dimensiones (alumno / tema /
// valija) y lo persiste de forma idempotente por conversation_id.
func (uc *closeSessionImpl) Execute(ctx context.Context, req CloseSessionRequest) (*CloseSessionResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	conv, err := uc.conversations.GetWithMessages(ctx, req.OrgID, req.ConversationID)
	if err != nil {
		return nil, err
	}

	// Sin mensajes no hay nada que compactar: evitamos un llamado al modelo y un
	// resumen vacío en DB.
	if len(conv.Messages) == 0 {
		return &CloseSessionResponse{ConversationID: conv.ID}, nil
	}

	transcript := buildTranscript(conv.Messages, maxSummaryInputTokens)

	resp, err := uc.ai.Chat(ctx, []providers.ChatMessage{
		{Role: "system", Content: summarizerSystemPrompt},
		{Role: "user", Content: transcript},
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", providers.ErrServiceUnavailable, err)
	}
	recordAIUsage(ctx, uc.usage, req.OrgID, req.UserID, "close", resp.Usage)

	summaryText, keywords := parseSummary(resp.Content)
	studentIDs, deviceIDs := collectEntities(conv)

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

const summarizerSystemPrompt = "Sos Alizia, asistente de inclusión. Resumí la siguiente conversación entre un docente y vos " +
	"para poder retomar el hilo más adelante sin recontextualizar todo. Devolvé EXCLUSIVAMENTE un JSON con esta forma:\n" +
	"{\"summary\": \"un par de párrafos en español con lo trabajado, decisiones y próximos pasos\", " +
	"\"topic_keywords\": [\"palabras clave en minúscula: temas y discapacidades tratadas\"]}\n" +
	"No incluyas nombres propios de alumnos ni diagnósticos en las keywords; usá temas (ej. \"autismo\", \"lectura\", \"autorregulación\"). " +
	"No agregues texto fuera del JSON."

// buildTranscript arma el texto de la conversación para resumir, preservando los
// turnos más recientes hasta agotar maxTokens (la recencia es lo que importa para
// retomar el hilo). Recorre de atrás hacia adelante y arma el texto en orden.
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

// parseSummary extrae el resumen y las keywords del JSON que devuelve el modelo.
// Es tolerante: si el modelo envuelve el JSON en texto o falla el parseo, cae al
// contenido crudo como resumen sin keywords (nunca rompe el cierre de sesión).
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

// extractJSONObject devuelve el primer objeto JSON balanceado dentro de s (el modelo
// a veces lo envuelve en fences o prosa). Vacío si no encuentra uno.
func extractJSONObject(s string) string {
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")
	if start == -1 || end == -1 || end < start {
		return ""
	}
	return s[start : end+1]
}

// normalizeKeywords limpia, deduplica y acota las keywords a minúscula.
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

// collectEntities reúne los ids de alumnos y devices vinculados a la conversación,
// para taggear el resumen por dimensión (alumno / valija). El alumno de cabecera de
// la conversación más los detectados turno a turno en la metadata.
func collectEntities(conv *entities.Conversation) (studentIDs, deviceIDs []int64) {
	students := newIDSet()
	devices := newIDSet()

	if conv.StudentID != nil {
		students.add(*conv.StudentID)
	}
	for i := range conv.Messages {
		meta := decodeMessageMeta(conv.Messages[i].Metadata)
		if id, ok := metaInt64(meta, "identified_student"); ok {
			students.add(id)
		}
		if id, ok := metaInt64(meta, "recommended_device"); ok {
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

// metaInt64 lee un id numérico de la metadata. El JSON lo decodifica como float64.
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

// idSet preserva orden de inserción y deduplica ids positivos.
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
