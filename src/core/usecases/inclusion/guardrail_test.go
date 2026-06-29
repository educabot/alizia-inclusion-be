package inclusion

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_crossedClinicalLine_Disparos(t *testing.T) {
	cases := []struct {
		name   string
		text   string
		reason string
	}{
		{"diagnostico afirmado", "Por lo que contás, el diagnóstico es TDAH y hay que trabajarlo así.", "diagnosis"},
		{"caso de autismo", "Se trata de un caso de autismo, claramente.", "diagnosis"},
		{"padece condicion", "El nene padece dislexia, eso explica todo.", "diagnosis"},
		{"certeza diagnostica", "Definitivamente tiene un trastorno del espectro.", "diagnosis"},
		{"medicacion", "Lo mejor es medicarlo para que se concentre en clase.", "prescription"},
		{"tratamiento", "El tratamiento es una dosis diaria por la mañana.", "prescription"},
		{"receta", "Conviene recetarle algo que baje la ansiedad.", "prescription"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tripped, reason := crossedClinicalLine(tc.text)
			assert.True(t, tripped, "debería disparar: %q", tc.text)
			assert.Equal(t, tc.reason, reason)
		})
	}
}

func Test_crossedClinicalLine_NoDispara(t *testing.T) {
	cases := []struct {
		name string
		text string
	}{
		{"derivacion mencionando condicion", "No puedo decirte si tiene autismo; eso lo ve mejor el equipo de orientación. Mientras tanto, en el aula podemos anticipar la consigna."},
		{"condicion condicional", "Si tiene un diagnóstico de dislexia, podría ayudar ofrecer la consigna también en audio."},
		{"barrera observable pedagogica", "Veo que le cuesta sostener la atención en tareas largas; probá dividir la actividad en pasos cortos."},
		{"mencion sin acto diagnostico", "El TDAH suele asociarse a dificultades de autorregulación; desde DUA conviene dar apoyos visuales."},
		{"deriva por medicacion al medico", "La medicación la define el médico tratante; yo te acompaño con lo del aula."},
		{"texto pedagogico limpio", "Probá anticipar la consigna en pasos cortos y dar más tiempo. ¿Querés que lo guarde como recurso?"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tripped, _ := crossedClinicalLine(tc.text)
			assert.False(t, tripped, "NO debería disparar: %q", tc.text)
		})
	}
}
