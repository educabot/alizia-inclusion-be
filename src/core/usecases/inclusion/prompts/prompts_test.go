package prompts_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/usecases/inclusion/prompts"
)

func ptr(s string) *string { return &s }

func TestRecommendSystem_ContainsDeviceInfo(t *testing.T) {
	devices := []entities.Device{
		{ID: 1, Name: "Timer Visual", NeedsDescription: ptr("Para alumnos con distracción"),
			Rationale: ptr("Estructura el tiempo"), HowToUse: ptr("Mostralo al inicio")},
	}

	prompt := prompts.RecommendSystem(devices)

	assert.Contains(t, prompt, "Timer Visual")
	assert.Contains(t, prompt, "[ID:1]")
	assert.Contains(t, prompt, "ADAPTATION_JSON")
	assert.Contains(t, prompt, "Para alumnos con distracción")
	assert.Contains(t, prompt, "Estructura el tiempo", "el catálogo detallado incluye el fundamento")
}

func TestAssistSystem_HasConversationalRules(t *testing.T) {
	// The assist frame is conversation-first: short by default, asks when intent is
	// unclear, and gates the save offer behind explicit confirmation.
	prompt := prompts.AssistSystem(nil, nil)

	assert.Contains(t, prompt, "breve y al grano")
	assert.Contains(t, prompt, "una sola pregunta")
	assert.Contains(t, prompt, "rioplatense")
	assert.Contains(t, prompt, "guarde como recurso", "ofrece guardar solo tras confirmación")
	assert.Contains(t, prompt, "EJEMPLOS DE CONVERSACIÓN", "incluye el few-shot conversacional")
}

func TestAssistSystem_DoesNotMandateDifferentiationMatrix(t *testing.T) {
	// The old prompt forced a 3-level matrix on every turn; the new one offers it
	// only when an adaptation is actually being built.
	prompt := prompts.AssistSystem(nil, nil)

	assert.NotContains(t, prompt, "obligatorio")
	assert.NotContains(t, prompt, "al menos 3 niveles")
}

func TestAssistSystem_ListsStudentsAndDevices(t *testing.T) {
	students := []entities.Student{
		{ID: 5, Name: "Ana", Profile: &entities.StudentProfile{Difficulties: []string{"se_distrae"}}},
	}
	devices := []entities.Device{{ID: 9, Name: "Auriculares"}}

	prompt := prompts.AssistSystem(devices, students)

	assert.Contains(t, prompt, "[ID:5] Ana")
	assert.Contains(t, prompt, "se_distrae")
	assert.Contains(t, prompt, "[ID:9] Auriculares")
}

func TestAssistSystem_OmitsRosterWhenNoStudents(t *testing.T) {
	prompt := prompts.AssistSystem([]entities.Device{{ID: 1, Name: "X"}}, nil)

	assert.False(t, strings.Contains(prompt, "ALUMNOS DEL AULA"))
}

func TestAssistSystem_EmbedsOutOfScopeOffRamp(t *testing.T) {
	// The frame supplies the model with the exact off-ramp wording.
	prompt := prompts.AssistSystem(nil, nil)

	assert.Contains(t, prompt, prompts.OffRampOutOfScope)
}

func TestOffRamp_WordingDoesNotDiagnose(t *testing.T) {
	// Off-ramp must redirect, never diagnose.
	assert.NotEmpty(t, prompts.OffRampInvalidOutput)
	assert.Contains(t, prompts.OffRampOutOfScope, "derivar")
	assert.Contains(t, prompts.OffRampOutOfScope, "No puedo")
}
