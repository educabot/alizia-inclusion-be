package inclusion

import "encoding/json"

// Buckets de tools por fuente, para la trazabilidad de chat.sources_used.
var (
	studentTools = map[string]bool{
		"list_classroom_students": true,
		"get_student":             true,
		"get_student_history":     true,
		"get_past_adaptations":    true,
	}
	valijaTools = map[string]bool{
		"list_devices": true,
	}
	ragTools = map[string]bool{
		"search_content":         true,
		"search_content_hibrido": true,
		"get_content":            true,
	}
)

// sourcesSummary resume, para un turno del chat, de qué fuentes se sacó info
// activamente (vía tools): valija, datos de alumno y/o el corpus RAG.
type sourcesSummary struct {
	UsedValija  bool
	UsedStudent bool
	UsedRAG     bool
	StudentIDs  []int64
	RAGQueries  []string
	RAGHits     int
	Tools       []string
}

// summarizeSources recorre el trace de tools del turno y arma el resumen de fuentes.
func summarizeSources(trace []toolInvocation) sourcesSummary {
	var s sourcesSummary
	seenStudent := map[int64]bool{}
	for i := range trace {
		inv := &trace[i]
		s.Tools = append(s.Tools, inv.Name)
		switch {
		case valijaTools[inv.Name]:
			s.UsedValija = true
		case studentTools[inv.Name]:
			s.UsedStudent = true
			if inv.StudentID != nil && !seenStudent[*inv.StudentID] {
				seenStudent[*inv.StudentID] = true
				s.StudentIDs = append(s.StudentIDs, *inv.StudentID)
			}
		case ragTools[inv.Name]:
			s.UsedRAG = true
			s.RAGHits += inv.Hits
			if inv.Query != "" {
				s.RAGQueries = append(s.RAGQueries, inv.Query)
			}
		}
	}
	return s
}

// contentRefsFromTrace junta las referencias de contenido pedagógico que las tools
// RAG trajeron en el turno, deduplicadas por id y en orden de aparición. Alimenta
// referenced_content para que el FE resuelva el título de los chips [CONTENT_ID:X].
func contentRefsFromTrace(trace []toolInvocation) []ContentRef {
	seen := map[int64]bool{}
	var out []ContentRef
	for i := range trace {
		for _, ref := range trace[i].ContentRefs {
			if ref.ID == 0 || seen[ref.ID] {
				continue
			}
			seen[ref.ID] = true
			out = append(out, ref)
		}
	}
	return out
}

// extractToolStudentID lee el student_id de los argumentos de una tool (si lo trae).
func extractToolStudentID(args string) *int64 {
	var a struct {
		StudentID *int64 `json:"student_id"`
	}
	if json.Unmarshal([]byte(args), &a) != nil {
		return nil
	}
	return a.StudentID
}

// extractToolQuery devuelve la consulta de una tool RAG: la pregunta semántica
// (search_content_hibrido) o la query de palabras clave (search_content).
func extractToolQuery(args string) string {
	var a struct {
		Query            string `json:"query"`
		SemanticQuestion string `json:"semantic_question"`
	}
	_ = json.Unmarshal([]byte(args), &a)
	if a.SemanticQuestion != "" {
		return a.SemanticQuestion
	}
	return a.Query
}

// countResults cuenta los elementos del array "results" de un resultado de tool.
func countResults(result string) int {
	var r struct {
		Results []json.RawMessage `json:"results"`
	}
	if json.Unmarshal([]byte(result), &r) != nil {
		return 0
	}
	return len(r.Results)
}

// extractContentRefs saca las referencias (id + título) del resultado de una tool RAG.
// Para search_content_hibrido el id es resource_id (rag_resources); para search_content
// best-effort sobre id/resource_id. Otras tools no aportan refs.
func extractContentRefs(toolName, result string) []ContentRef {
	if !ragTools[toolName] {
		return nil
	}
	var r struct {
		Results []struct {
			ResourceID int64  `json:"resource_id"`
			ID         int64  `json:"id"`
			Title      string `json:"title"`
		} `json:"results"`
	}
	if json.Unmarshal([]byte(result), &r) != nil {
		return nil
	}
	out := make([]ContentRef, 0, len(r.Results))
	for _, h := range r.Results {
		id := h.ResourceID
		if id == 0 {
			id = h.ID
		}
		if id == 0 {
			continue
		}
		out = append(out, ContentRef{ID: id, Title: h.Title})
	}
	return out
}
