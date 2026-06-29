package inclusion

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeGrade(t *testing.T) {
	cases := []struct {
		in      string
		name    string
		grade   string
		section string
	}{
		{"3ro A", "3ro A", "3ro", "A"},
		{"tercero B", "3ro B", "3ro", "B"},
		{"3° A", "3ro A", "3ro", "A"},
		{"3 A", "3ro A", "3ro", "A"},
		{"3roA", "3ro A", "3ro", "A"},
		{"primero", "1ro", "1ro", ""},
		{"séptimo", "7mo", "7mo", ""},
		{"primer grado a", "1ro A", "1ro", "A"},
		{"  Cuarto   c ", "4to C", "4to", "C"},
		{"sala roja", "", "", ""},
		{"", "", "", ""},
		{"9no", "", "", ""}, // fuera de rango 1..7
	}
	for _, c := range cases {
		t.Run(c.in, func(t *testing.T) {
			name, grade, section := normalizeGrade(c.in)
			assert.Equal(t, c.name, name, "name")
			assert.Equal(t, c.grade, grade, "grade")
			assert.Equal(t, c.section, section, "section")
		})
	}
}

func TestNormalizeName(t *testing.T) {
	assert.Equal(t, "lucas perez", normalizeName("  Lucas   Pérez "))
	assert.Equal(t, "martin nunez", normalizeName("Martín Núñez"))
	assert.Equal(t, normalizeName("JOSÉ"), normalizeName("jose"))
}
