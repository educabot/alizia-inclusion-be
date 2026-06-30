package inclusion

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
)

func ptr(s string) *string { return &s }

func Test_writeStudentContext_NilEscribeNada(t *testing.T) {
	var b strings.Builder
	writeStudentContext(&b, nil)
	assert.Empty(t, b.String())
}

func Test_writeStudentContext_RenderizaContextoRico(t *testing.T) {
	pc := &PromptContext{
		Dimension: DimensionStudent,
		TargetStudent: &entities.Student{
			ID:   7,
			Name: "Lucas",
			Profile: &entities.StudentProfile{
				IsTransitory:        false,
				Difficulties:        []string{"sostener la atención"},
				Strengths:           []string{"memoria visual"},
				EffectiveStrategies: []string{"consignas cortas"},
			},
		},
		Diagnoses: []entities.StudentDiagnosis{
			{Diagnosis: &entities.Diagnosis{Name: "TDAH"}, Severity: ptr("leve")},
		},
		PPI: &entities.PPI{Objectives: []string{"mejorar autorregulación"}},
		PastAdaptations: []entities.Adaptation{
			{Title: "Tablero de anticipación", Subject: "Lengua"},
		},
		PriorSummaries: []entities.ConversationSummary{
			{Summary: "Trabajamos rutinas visuales la clase pasada."},
		},
		MissingData: []string{missingPPI},
	}

	var b strings.Builder
	writeStudentContext(&b, pc)
	out := b.String()

	assert.Contains(t, out, "ALUMNO FOCO: Lucas [ID:7]")
	assert.Contains(t, out, "memoria visual")
	assert.Contains(t, out, "consignas cortas")
	// El diagnóstico viaja como contexto, enmarcado para que NO se afirme como propio.
	assert.Contains(t, out, "TDAH (leve)")
	assert.Contains(t, out, "no los afirmes como propios")
	assert.Contains(t, out, "Tablero de anticipación (Lengua)")
	assert.Contains(t, out, "Trabajamos rutinas visuales")
}

func Test_writeStudentContext_SinAlumnoSugiereDatosFaltantes(t *testing.T) {
	pc := &PromptContext{MissingData: []string{missingStudentProfile, missingDiagnoses}}

	var b strings.Builder
	writeStudentContext(&b, pc)
	out := b.String()

	assert.NotContains(t, out, "ALUMNO FOCO")
	assert.Contains(t, out, "DATOS QUE FALTAN")
	assert.Contains(t, out, "el perfil del alumno")
}

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

func Test_extractQuestions_ExtractsAllThreeTypes(t *testing.T) {
	input := `Para ayudarte con María, necesito entender un poco más: [QUESTIONS_JSON:{"questions":[` +
		`{"id":"edad","text":"¿Qué edad tiene?","type":"open"},` +
		`{"id":"momento","text":"¿Qué momento se le dificulta?","type":"single","options":["Autonomía","Grupal"]},` +
		`{"id":"tipo","text":"¿Qué observás?","type":"multiple","options":["Pasiva","Activa"]}]}]`

	result := extractQuestions(input)

	require.Len(t, result, 3)
	assert.Equal(t, "edad", result[0].ID)
	assert.Equal(t, "open", result[0].Type)
	assert.Empty(t, result[0].Options, "una pregunta abierta no debe tener opciones")
	assert.Equal(t, "single", result[1].Type)
	assert.Len(t, result[1].Options, 2)
	assert.Equal(t, "multiple", result[2].Type)
}

func Test_extractQuestions_ReturnsNilForNoMarker(t *testing.T) {
	assert.Nil(t, extractQuestions("Probá anticipar la consigna en pasos cortos."))
}

func Test_extractQuestions_ReturnsNilForMalformedJSON(t *testing.T) {
	assert.Nil(t, extractQuestions("[QUESTIONS_JSON:{not valid}]"))
}

func Test_extractQuestions_DropsInvalidTypesAndEmptyText(t *testing.T) {
	// Tipo desconocido y pregunta sin texto se descartan; queda solo la válida.
	input := `[QUESTIONS_JSON:{"questions":[` +
		`{"id":"a","text":"","type":"open"},` +
		`{"id":"b","text":"¿Sí?","type":"checkbox"},` +
		`{"id":"c","text":"¿Qué edad?","type":"open"}]}]`

	result := extractQuestions(input)

	require.Len(t, result, 1)
	assert.Equal(t, "c", result[0].ID)
}

func Test_extractQuestions_OpenTypeIgnoresOptions(t *testing.T) {
	// Una "open" con opciones (el modelo se equivocó): se limpian las opciones.
	input := `[QUESTIONS_JSON:{"questions":[{"id":"x","text":"¿Edad?","type":"open","options":["a","b"]}]}]`

	result := extractQuestions(input)

	require.Len(t, result, 1)
	assert.Empty(t, result[0].Options)
}

func Test_stripAdaptationBlock_RemovesQuestionsBlock(t *testing.T) {
	in := `Para ayudarte necesito saber: [QUESTIONS_JSON:{"questions":[{"id":"e","text":"¿Edad?","type":"open"}]}]`

	out := stripAdaptationBlock(in)

	assert.NotContains(t, out, "[QUESTIONS_JSON:")
	assert.Equal(t, "Para ayudarte necesito saber:", out)
}

func Test_buildAssistSystemPrompt_ContainsQuestionsBlockFormat(t *testing.T) {
	prompt := buildAssistSystemPrompt(nil, nil, nil, nil, false)

	assert.Contains(t, prompt, "PREGUNTAS AL DOCENTE (bloque estructurado)")
	assert.Contains(t, prompt, "[QUESTIONS_JSON:")
	assert.Contains(t, prompt, `"type":"open"`)
	assert.Contains(t, prompt, `"type":"single"`)
	assert.Contains(t, prompt, `"type":"multiple"`)
}

func Test_buildAssistSystemPrompt_ContainsRepreguntaGate(t *testing.T) {
	devices := []entities.Device{{ID: 1, Name: "Organizador visual"}}

	prompt := buildAssistSystemPrompt(devices, nil, nil, nil, false)

	assert.Contains(t, prompt, "ANTES DE PROPONER")
	assert.Contains(t, prompt, "CÓMO PREGUNTÁS")
	assert.Contains(t, prompt, "HASTA 4 opciones")
	assert.Contains(t, prompt, `"Otro"`)
	assert.Contains(t, prompt, "Abierta")
	assert.Contains(t, prompt, "De opción única")
	assert.Contains(t, prompt, "De opción múltiple")
}

func Test_buildAssistSystemPrompt_ContainsFirstProposalAndWarmClose(t *testing.T) {
	prompt := buildAssistSystemPrompt(nil, nil, nil, nil, false)

	assert.Contains(t, prompt, "PROPONÉ, NO INTERROGUES")
	assert.Contains(t, prompt, "PRIMERA propuesta concreta")
	assert.Contains(t, prompt, "¿Continuamos?")
}

func Test_buildAssistSystemPrompt_ContainsStepsFormat(t *testing.T) {
	prompt := buildAssistSystemPrompt(nil, nil, nil, nil, false)

	assert.Contains(t, prompt, "FORMATO DE RESPUESTA CON PASOS")
	assert.Contains(t, prompt, "[STEPS]")
	assert.Contains(t, prompt, "[/STEPS]")
	assert.Contains(t, prompt, "LINEAMIENTOS")
}

func Test_buildAssistSystemPrompt_AgenticInjectsFundamentos(t *testing.T) {
	promptSinAgentic := buildAssistSystemPrompt(nil, nil, nil, nil, false)
	promptConAgentic := buildAssistSystemPrompt(nil, nil, nil, nil, true)

	// search_content_hibrido solo existe en fundamentosRAG (agentic=true)
	assert.NotContains(t, promptSinAgentic, "search_content_hibrido")
	assert.Contains(t, promptConAgentic, "search_content_hibrido")
	assert.Contains(t, promptConAgentic, "[CONTENT_ID:")
	// Reforzamos buscar ANTES de preguntar para preguntar mejor.
	assert.Contains(t, promptConAgentic, "BUSCÁ ANTES DE PREGUNTAR")
}

func Test_buildAssistSystemPrompt_AsksForAgeNotGrade(t *testing.T) {
	prompt := buildAssistSystemPrompt(nil, nil, nil, nil, false)

	assert.Contains(t, prompt, "¿Qué edad tiene?")
	assert.NotContains(t, prompt, "edad o grado")
}

func Test_buildAssistSystemPrompt_AutoSavesResourceWithProposal(t *testing.T) {
	prompt := buildAssistSystemPrompt(nil, nil, nil, nil, false)

	// El recurso se guarda solo junto con la propuesta; ya NO se pide permiso.
	assert.Contains(t, prompt, "GENERAR Y GUARDAR EL RECURSO")
	assert.Contains(t, prompt, "se guarda solo")
	assert.NotContains(t, prompt, "¿Querés que la guarde como recurso?")
	assert.NotContains(t, prompt, "turno POSTERIOR")
}

func Test_buildAssistSystemPrompt_ListsConversationResourcesForUpdate(t *testing.T) {
	resources := []entities.Adaptation{{ID: 42, Title: "Fragmentar consignas", Status: "en_curso"}}

	prompt := buildAssistSystemPrompt(nil, nil, nil, resources, false)

	assert.Contains(t, prompt, "RECURSOS YA GUARDADOS EN ESTA CONVERSACIÓN")
	assert.Contains(t, prompt, "[REC_ID:42]")
	assert.Contains(t, prompt, "Fragmentar consignas")
	// Instrucción de update-vs-create por id.
	assert.Contains(t, prompt, `campo "id"`)
}

func Test_extractAdaptationJSON_ParsesID(t *testing.T) {
	withID := `[ADAPTATION_JSON:{"id":42,"title":"t","type":"estrategia_aula","strategy":"s","steps":[{"orden":1,"texto":"x"}]}]`
	withoutID := `[ADAPTATION_JSON:{"title":"t","type":"estrategia_aula","strategy":"s"}]`

	got := extractAdaptationJSON(withID)
	require.NotNil(t, got)
	require.NotNil(t, got.ID)
	assert.Equal(t, int64(42), *got.ID)

	none := extractAdaptationJSON(withoutID)
	require.NotNil(t, none)
	assert.Nil(t, none.ID)
}

func Test_sanitizeVisibleText(t *testing.T) {
	// Marker mal formado (nombre, no id): se desenvuelve al texto interno.
	assert.Equal(t, "en grupales a María se mueve", sanitizeVisibleText("en grupales a [STUDENT_ID:María] se mueve"))
	// Marker numérico válido: se conserva para que el FE lo renderice como chip.
	assert.Equal(t, "trabajá con [STUDENT_ID:7] así", sanitizeVisibleText("trabajá con [STUDENT_ID:7] así"))
	// Llave huérfana de un bloque JSON mal cerrado: se quita.
	assert.Equal(t, "y la hacemos más a medida.", sanitizeVisibleText("y la hacemos más a medida. }"))
	// DEVICE_ID mal formado también se desenvuelve.
	assert.Equal(t, "usá el Cojín dinámico", sanitizeVisibleText("usá el [DEVICE_ID:Cojín dinámico]"))
}

func Test_aliziaPersona_ReinforcesNoDiagnosis(t *testing.T) {
	assert.Contains(t, aliziaPersona, "No diagnosticás ni insinuás un diagnóstico")
	assert.Contains(t, aliziaPersona, "No abrís con empatía en abstracto ni con soluciones genéricas")
}

// El alumno nunca se enmarca como una molestia a contener/reubicar: el objetivo es la
// participación, no "que moleste menos". Encuadre inclusivo, no de gestión de la disrupción.
func Test_aliziaPersona_RejectsNuisanceFraming(t *testing.T) {
	assert.Contains(t, aliziaPersona, `que moleste menos`)
	assert.Contains(t, aliziaPersona, "participar y aprender")
	assert.Contains(t, aliziaPersona, "no de evitar que")
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
