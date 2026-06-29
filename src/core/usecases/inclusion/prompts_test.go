package inclusion

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
)

func ptr(s string) *string { return &s }

func Test_stripInternalMarkers_RemovesTagsFromVisibleText(t *testing.T) {
	// Arrange: respuesta del modelo con marcadores internos intercalados en la prosa.
	in := "Para [STUDENT_ID:2] que le cuesta escribir, probá el cronómetro [DEVICE_ID:5]. Listo. [ADAPTATION_JSON:{\"title\":\"x\",\"type\":\"t\",\"strategy\":\"s\",\"device_ids\":[1],\"device_names\":[\"n\"]}]"

	// Act
	out := stripInternalMarkers(in)

	// Assert: ningún marcador interno queda visible y los espacios quedan prolijos.
	assert.NotContains(t, out, "[STUDENT_ID:")
	assert.NotContains(t, out, "[DEVICE_ID:")
	assert.NotContains(t, out, "[ADAPTATION_JSON:")
	assert.Equal(t, "Para que le cuesta escribir, probá el cronómetro. Listo.", out)
}

func Test_stripInternalMarkers_LeavesCleanTextUnchanged(t *testing.T) {
	in := "Probá anticipar la consigna en pasos cortos. ¿Querés que la guarde como recurso?"
	assert.Equal(t, in, stripInternalMarkers(in))
}

func Test_extractDeviceID_ExtractsValidID(t *testing.T) {
	input := "Use this [DEVICE_ID:42] for the student"

	result := extractDeviceID(input)

	require.NotNil(t, result)
	assert.Equal(t, int64(42), *result)
}

func Test_extractDeviceID_ReturnsNilForNoMatch(t *testing.T) {
	input := "no device here"

	result := extractDeviceID(input)

	assert.Nil(t, result)
}

func Test_extractDeviceID_ReturnsNilForInvalidNumber(t *testing.T) {
	input := "[DEVICE_ID:abc]"

	result := extractDeviceID(input)

	assert.Nil(t, result)
}

func Test_extractStudentID_ExtractsValidID(t *testing.T) {
	input := "Student [STUDENT_ID:7] needs help"

	result := extractStudentID(input)

	require.NotNil(t, result)
	assert.Equal(t, int64(7), *result)
}

func Test_extractStudentID_ReturnsNilForNoMatch(t *testing.T) {
	input := "no student"

	result := extractStudentID(input)

	assert.Nil(t, result)
}

func Test_extractAdaptationJSON_ExtractsValidJSON(t *testing.T) {
	input := `text [ADAPTATION_JSON:{"title":"Test","type":"actividad_adaptada","strategy":"do thing","device_ids":[1,2],"device_names":["A","B"]}] more text`

	result := extractAdaptationJSON(input)

	require.NotNil(t, result)
	assert.Equal(t, "Test", result.Title)
	assert.Equal(t, "actividad_adaptada", result.Type)
	assert.Len(t, result.DeviceIDs, 2)
}

func Test_extractAdaptationJSON_ReturnsNilForNoMatch(t *testing.T) {
	input := "no json here"

	result := extractAdaptationJSON(input)

	assert.Nil(t, result)
}

func Test_extractAdaptationJSON_ReturnsNilForMalformedJSON(t *testing.T) {
	input := "[ADAPTATION_JSON:{invalid json}]"

	result := extractAdaptationJSON(input)

	assert.Nil(t, result)
}

func Test_buildAssistSystemPrompt_ContainsQuestioningCriteria(t *testing.T) {
	// Arrange
	devices := []entities.Device{{ID: 1, Name: "Organizador visual"}}

	// Act: el bloque de preguntas se inyecta siempre (no depende de agentic).
	prompt := buildAssistSystemPrompt(devices, nil, false)

	// Assert: están los criterios definidos con pedagogía (Mercedes).
	assert.Contains(t, prompt, "CÓMO PREGUNTÁS")
	assert.Contains(t, prompt, "2-3 preguntas base en el MISMO turno")
	assert.Contains(t, prompt, "HASTA 4 opciones")
	assert.Contains(t, prompt, `"Otro"`, "siempre debe ofrecer la opción Otro de texto libre")
	assert.Contains(t, prompt, "Abierta")
	assert.Contains(t, prompt, "De opción única")
	assert.Contains(t, prompt, "De opción múltiple")
}

func Test_buildAssistSystemPrompt_ContainsFirstProposalAndWarmClose(t *testing.T) {
	prompt := buildAssistSystemPrompt(nil, nil, false)

	assert.Contains(t, prompt, "PROPONÉ, NO INTERROGUES")
	assert.Contains(t, prompt, "PRIMERA propuesta concreta")
	assert.Contains(t, prompt, "¿Continuamos?")
}

func Test_buildAssistSystemPrompt_AgenticDoesNotCiteSources(t *testing.T) {
	// Act: con agentic se inyecta FUNDAMENTOS (RAG).
	prompt := buildAssistSystemPrompt(nil, nil, true)

	// Assert: ya no instruye citar la fuente con el marker de contenido,
	// pero sí usar el RAG para repreguntar mejor.
	assert.NotContains(t, prompt, "[CONTENT_ID:")
	assert.Contains(t, prompt, "SIN citar la fuente")
	assert.Contains(t, prompt, "ANTES DE REPREGUNTAR")
}

func Test_aliziaPersona_ReinforcesNoDiagnosis(t *testing.T) {
	assert.Contains(t, aliziaPersona, "No diagnosticás ni insinuás un diagnóstico")
	assert.Contains(t, aliziaPersona, "No abrís con empatía en abstracto ni con soluciones genéricas")
}

func Test_buildRecommendSystemPrompt_ContainsDeviceInfo(t *testing.T) {
	devices := []entities.Device{
		{ID: 1, Name: "Timer Visual", NeedsDescription: ptr("Para alumnos con distraccion")},
	}

	prompt := buildRecommendSystemPrompt(devices)

	assert.True(t, strings.Contains(prompt, "Timer Visual"), "prompt should contain device name")
	assert.True(t, strings.Contains(prompt, "[ID:1]"), "prompt should contain device ID")
	assert.True(t, strings.Contains(prompt, "ADAPTATION_JSON"), "prompt should contain format instructions")
	assert.True(t, strings.Contains(prompt, "Para alumnos con distraccion"), "prompt should contain device needs description")
}
