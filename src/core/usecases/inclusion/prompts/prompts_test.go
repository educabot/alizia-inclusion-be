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

func TestAssistSystem_HasOutputFormatRules(t *testing.T) {
	// Required output format lives in the static frame.
	prompt := prompts.AssistSystem(nil, nil)

	assert.Contains(t, prompt, "1 a 3 acciones")
	assert.Contains(t, prompt, "3 niveles de diferenciación")
	assert.Contains(t, prompt, "rioplatense")
	assert.Contains(t, prompt, "EJEMPLO DE BUENA RESPUESTA", "incluye el few-shot estático")
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
