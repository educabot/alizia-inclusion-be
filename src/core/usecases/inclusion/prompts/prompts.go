// Package prompts is the editable layer of Alizia's pedagogical framework.
// It separates the engine (usecases) from prompt text across two layers:
//
//   - Layer 1 (static, cacheable): role, guidelines, scope limits, output format,
//     and few-shot examples. These are constants; changing text does not break output format.
//   - Layer 2 (dynamic): per-request turn context (classroom students, adaptive kit catalog).
//
// There is a single mode: Alizia adapts through context, not prompt selection.
package prompts

import (
	"fmt"
	"strings"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
)

// ---------- Layer 1 — static framework (cacheable) ----------

// outputRules enforces the required output format: few actions,
// with differentiation levels, immediately actionable, in Rioplatense Spanish register.
const outputRules = "FORMATO DE RESPUESTA (obligatorio):\n" +
	"- Ofrecé 1 a 3 acciones concretas, ordenadas por impacto.\n" +
	"- Incluí al menos 3 niveles de diferenciación (más simple / intermedio / más desafiante).\n" +
	"- Aplicable en menos de 1 minuto de lectura: breve y al grano.\n" +
	"- Español rioplatense, tono cálido, sin jerga clínica.\n"

// scopeRules establishes the hard boundaries of the framework: never diagnose,
// never replace the teacher; the off-ramp is the last resort, not the first.
const scopeRules = "LÍMITES (marco pedagógico):\n" +
	"- Entrada pedagógica, no clínica: partís de situaciones de aula, no de diagnósticos.\n" +
	"- Nunca diagnostiques, no reemplaces al docente, no produzcas informes clínicos.\n" +
	"- Si el caso se va de tu alcance (clínico/crisis/diagnóstico), primero ayudá con lo que tengas;\n" +
	"  dar un paso al costado es el último recurso, no el primero. Cuando debas hacerlo, respondé\n" +
	"  exactamente con: \"" + OffRampOutOfScope + "\"\n"

// adaptationBlock instructs the structured block that the guardrail validates programmatically.
const adaptationBlock = "BLOQUE ESTRUCTURADO:\n" +
	"Cuando generes una recomendación de adaptación concreta, incluí al final exactamente:\n" +
	`[ADAPTATION_JSON:{"title":"título corto","type":"tipo","strategy":"resumen","device_ids":[1],"device_names":["nombre"]}]` + "\n" +
	"Tipos válidos: actividad_adaptada, material_nuevo, estrategia_aula, situacion_emergente.\n" +
	"Solo incluí el bloque ante una adaptación concreta, no en preguntas ni aclaraciones.\n"

const assistFramework = "Sos Alizia, asistente de inclusión educativa en tiempo real.\n" +
	"Acompañás a un docente DURANTE la clase: ayudás a adaptar la enseñanza, no a intervenir sobre el alumno.\n\n" +
	scopeRules + "\n" +
	"Si detectás el nombre de un alumno, usá [STUDENT_ID:X]. Si recomendás un dispositivo, usá [DEVICE_ID:X].\n\n" +
	outputRules + "\n" + adaptationBlock

const recommendFramework = "Sos Alizia, asistente de inclusión educativa de Educabot.\n" +
	"Ayudás al docente a planificar actividades inclusivas recomendando dispositivos de la valija adaptativa.\n\n" +
	scopeRules + "\n" + outputRules +
	"Incluí [DEVICE_ID:X] con el dispositivo principal recomendado.\n\n" + adaptationBlock

// fewShot is a curated static example (layer 1, cacheable): classroom situation →
// response with actions and differentiation levels. Per-student dynamic few-shot
// (built from past successful adaptations) is a planned extension on this base.
const fewShot = "EJEMPLO DE BUENA RESPUESTA:\n" +
	"Situación: el alumno no inicia la tarea.\n" +
	"1. Anticipá la consigna en pasos cortos en el pizarrón.\n" +
	"   - Más simple: un paso por vez. Intermedio: 2-3 pasos. Más desafiante: que arme la secuencia él mismo.\n" +
	"2. Dale un arranque concreto (\"empezá por...\") y un tiempo breve para el primer paso.\n"

// ---------- Layer 2 — dynamic turn context ----------

// AssistSystem builds the in-class assistant system prompt: static framework +
// classroom context (students + kit) + few-shot. Single mode.
func AssistSystem(devices []entities.Device, students []entities.Student) string {
	var b strings.Builder
	b.WriteString(assistFramework)
	b.WriteString("\n")
	b.WriteString(studentRoster(students))
	b.WriteString(deviceCatalog(devices, false))
	b.WriteString("\n")
	b.WriteString(fewShot)
	return b.String()
}

// RecommendSystem builds the device recommendation system prompt: framework +
// detailed catalog (rationale / how-to-use) + few-shot.
func RecommendSystem(devices []entities.Device) string {
	var b strings.Builder
	b.WriteString(recommendFramework)
	b.WriteString("\n")
	b.WriteString(deviceCatalog(devices, true))
	b.WriteString("\n")
	b.WriteString(fewShot)
	return b.String()
}

// studentRoster lists classroom students with their ID and difficulties (layer 2).
func studentRoster(students []entities.Student) string {
	if len(students) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("ALUMNOS DEL AULA:\n")
	for i := range students {
		s := &students[i]
		fmt.Fprintf(&b, "- [ID:%d] %s", s.ID, s.Name)
		if s.Profile != nil && len(s.Profile.Difficulties) > 0 {
			fmt.Fprintf(&b, " — Dificultades: %s", strings.Join(s.Profile.Difficulties, ", "))
		}
		b.WriteString("\n")
	}
	b.WriteString("\n")
	return b.String()
}

// deviceCatalog lists the adaptive kit (layer 2). When detailed is true, adds rationale
// and how-to-use fields for the recommendation prompt; otherwise only name and purpose.
func deviceCatalog(devices []entities.Device, detailed bool) string {
	var b strings.Builder
	b.WriteString("DISPOSITIVOS DISPONIBLES:\n")
	for i := range devices {
		d := &devices[i]
		fmt.Fprintf(&b, "- [ID:%d] %s", d.ID, d.Name)
		if d.NeedsDescription != nil {
			fmt.Fprintf(&b, " — %s", *d.NeedsDescription)
		}
		b.WriteString("\n")
		if detailed {
			if d.Rationale != nil {
				fmt.Fprintf(&b, "  Fundamento: %s\n", *d.Rationale)
			}
			if d.HowToUse != nil {
				fmt.Fprintf(&b, "  Uso: %s\n", *d.HowToUse)
			}
		}
	}
	return b.String()
}
