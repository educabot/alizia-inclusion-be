package inclusion

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func int64Ptr(v int64) *int64 { return &v }

func TestSummarizeSources_MixedTraceFlagsEachSource(t *testing.T) {
	// Arrange: un turno que tocó valija, alumno (2 ids, uno repetido) y RAG.
	trace := []toolInvocation{
		{Name: "list_devices"},
		{Name: "get_student", StudentID: int64Ptr(7)},
		{Name: "get_past_adaptations", StudentID: int64Ptr(7)},
		{Name: "list_classroom_students", StudentID: int64Ptr(9)},
		{Name: "search_content_hibrido", Query: "TEA autorregulación", Hits: 3},
		{Name: "search_content", Query: "rampa silla", Hits: 2},
	}

	// Act
	s := summarizeSources(trace)

	// Assert
	assert.True(t, s.UsedValija)
	assert.True(t, s.UsedStudent)
	assert.True(t, s.UsedRAG)
	assert.Equal(t, []int64{7, 9}, s.StudentIDs)
	assert.Equal(t, 5, s.RAGHits)
	assert.Equal(t, []string{"TEA autorregulación", "rampa silla"}, s.RAGQueries)
	assert.Len(t, s.Tools, 6)
}

func TestSummarizeSources_EmptyTraceAllFalse(t *testing.T) {
	// Act
	s := summarizeSources(nil)

	// Assert
	assert.False(t, s.UsedValija)
	assert.False(t, s.UsedStudent)
	assert.False(t, s.UsedRAG)
	assert.Empty(t, s.StudentIDs)
	assert.Empty(t, s.RAGQueries)
	assert.Equal(t, 0, s.RAGHits)
	assert.Empty(t, s.Tools)
}

func TestContentRefsFromTrace_DedupesByIDPreservingOrder(t *testing.T) {
	// Arrange: el mismo recurso aparece en dos tools; debe quedar una sola vez.
	trace := []toolInvocation{
		{Name: "search_content_hibrido", ContentRefs: []ContentRef{
			{ID: 10, Title: "Guía TEA"},
			{ID: 20, Title: "Manual DUA"},
		}},
		{Name: "search_content", ContentRefs: []ContentRef{
			{ID: 20, Title: "Manual DUA"},
			{ID: 30, Title: "Rampa de aprendizaje"},
			{ID: 0, Title: "sin id, se ignora"},
		}},
	}

	// Act
	refs := contentRefsFromTrace(trace)

	// Assert
	assert.Equal(t, []ContentRef{
		{ID: 10, Title: "Guía TEA"},
		{ID: 20, Title: "Manual DUA"},
		{ID: 30, Title: "Rampa de aprendizaje"},
	}, refs)
}

func TestExtractContentRefs_HybridResultUsesResourceID(t *testing.T) {
	// Arrange
	result := `{"results":[{"resource_id":42,"title":"Guía X","score":0.9},{"resource_id":43,"title":"Guía Y"}]}`

	// Act
	refs := extractContentRefs("search_content_hibrido", result)

	// Assert
	assert.Equal(t, []ContentRef{{ID: 42, Title: "Guía X"}, {ID: 43, Title: "Guía Y"}}, refs)
}

func TestExtractContentRefs_NonRAGToolReturnsNil(t *testing.T) {
	assert.Nil(t, extractContentRefs("list_devices", `{"devices":[{"id":1}]}`))
}

func TestExtractToolQuery_PrefersSemanticQuestion(t *testing.T) {
	assert.Equal(t, "pregunta completa", extractToolQuery(`{"semantic_question":"pregunta completa","terms":["a"]}`))
	assert.Equal(t, "kw", extractToolQuery(`{"query":"kw"}`))
}
