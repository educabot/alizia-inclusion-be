package inclusion

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_resolveChatDimension(t *testing.T) {
	sid := int64(5)
	cases := []struct {
		name      string
		dimension string
		studentID *int64
		want      string
	}{
		{"explícita alumno gana sobre inferencia", "alumno", nil, DimensionStudent},
		{"explícita valija", "valija", &sid, DimensionToolkit},
		{"explícita tema", "tema", nil, DimensionTopic},
		{"normaliza mayúsculas y espacios", "  Alumno ", nil, DimensionStudent},
		{"sin dimensión infiere alumno por StudentID", "", &sid, DimensionStudent},
		{"sin dimensión ni alumno queda vacío", "", nil, ""},
		{"dimensión inválida cae a inferencia", "xxx", &sid, DimensionStudent},
		{"dimensión inválida sin alumno queda vacío", "xxx", nil, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, resolveChatDimension(tc.dimension, tc.studentID))
		})
	}
}
