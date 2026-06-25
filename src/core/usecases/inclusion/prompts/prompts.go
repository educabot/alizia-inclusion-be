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

// outputRules sets the default response shape: conversational and short, with the
// heavy differentiation matrix offered only when an actual adaptation is being built.
const outputRules = "CÓMO RESPONDER:\n" +
	"- Por defecto, respondé breve y al grano: 2 a 5 líneas, una idea principal por turno.\n" +
	"- Conversá: cerrá con una sola pregunta de seguimiento para entender mejor la situación.\n" +
	"- Recién cuando propongas una adaptación concreta, dá 1 a 3 acciones; sumá niveles\n" +
	"  (más simple / intermedio / más desafiante) solo si ayudan, nunca como relleno.\n" +
	"- Español rioplatense, tono cálido, sin jerga clínica.\n"

// offRampRule is the operational derivation instruction (behavior, not identity): the
// prose limits now live in persona / TU LUGAR. This only governs WHEN and HOW to emit the
// exact off-ramp wording the guardrail expects — the last resort, never the first.
const offRampRule = "DERIVACIÓN:\n" +
	"- Si el caso se va de tu alcance (clínico/crisis/diagnóstico), primero ayudá con lo que tengas;\n" +
	"  dar un paso al costado es el último recurso, no el primero. Cuando debas hacerlo, respondé\n" +
	"  exactamente con: \"" + OffRampOutOfScope + "\"\n"

// assistConversation drives the in-class chat: ask when intent is unclear, answer
// questions without building anything, and only offer to save a resource after the
// teacher confirms — so the structured block never appears prematurely.
const assistConversation = "CONVERSÁ Y GUIÁ (clave):\n" +
	"- Si no te queda claro de qué se trata, no supongas: preguntá si surgió algo con un alumno,\n" +
	"  si querés ver algo de la valija, o si es sobre un tema. Una sola pregunta a la vez.\n" +
	"- Si el docente pregunta o consulta algo, respondé esa pregunta de forma breve.\n" +
	"  No generes una adaptación ni propongas guardar nada todavía.\n" +
	"- Recién cuando estén trabajando una adaptación concreta para un alumno o situación,\n" +
	"  proponéla en texto y preguntá si la quiere guardar como recurso pedagógico\n" +
	"  (ej.: \"¿Querés que la guarde como recurso?\").\n" +
	"- Incluí el BLOQUE ESTRUCTURADO de abajo SOLO en el turno posterior, después de que el\n" +
	"  docente confirme que sí. Nunca en el primer mensaje, ni junto con la pregunta de\n" +
	"  confirmación, ni en respuestas a consultas.\n"

// adaptationBlock defines the structured block format the guardrail validates
// programmatically. WHEN to emit it is governed per-framework (assist gates it behind
// confirmation; recommend always emits it).
const adaptationBlock = "BLOQUE ESTRUCTURADO (formato):\n" +
	"Para registrar una adaptación como recurso, agregá al final exactamente:\n" +
	`[ADAPTATION_JSON:{"title":"título corto","type":"tipo","strategy":"resumen","device_ids":[1],"device_names":["nombre"]}]` + "\n" +
	"Tipos válidos: actividad_adaptada, material_nuevo, estrategia_aula, situacion_emergente.\n"

// assistFramework = persona (identity) + pedagogical frame + behavior specific to the
// in-class chat: ask when intent is unclear, gate the structured block behind confirmation.
const assistFramework = persona + "\n" +
	pedagogicalGuidelines + "\n" +
	"Acompañás al docente DURANTE la clase: ayudás a adaptar la enseñanza, no a intervenir sobre el alumno.\n\n" +
	offRampRule + "\n" +
	"Si detectás el nombre de un alumno, usá [STUDENT_ID:X]. Si recomendás un dispositivo, usá [DEVICE_ID:X].\n" +
	"Si mencionás un material real que viene de search_content, marcalo con [CONTENT_ID:X] usando el id del documento.\n\n" +
	assistConversation + "\n" + outputRules + "\n" + adaptationBlock

// recommendFramework = persona (identity) + pedagogical frame + behavior specific to device
// recommendation: always close with the structured block.
const recommendFramework = persona + "\n" +
	pedagogicalGuidelines + "\n" +
	"Ayudás al docente a planificar actividades inclusivas recomendando dispositivos de la valija adaptativa.\n\n" +
	offRampRule + "\n" + outputRules +
	"Incluí [DEVICE_ID:X] con el dispositivo principal recomendado.\n" +
	"Cerrá siempre con el BLOQUE ESTRUCTURADO de la adaptación recomendada.\n\n" + adaptationBlock

// assistFewShot models the conversational protocol: a question gets a short answer
// with no block; a concrete situation gets a proposal plus a save question, and the
// structured block appears only in the turn after the teacher confirms.
const assistFewShot = "EJEMPLOS DE CONVERSACIÓN:\n" +
	"- Docente: \"¿El cronómetro visual sirve para un nene que se desorganiza?\"\n" +
	"  Alizia: \"Sí, ayuda a anticipar el tiempo y las transiciones. ¿Querés que veamos cómo usarlo con tu grupo?\" (sin bloque)\n" +
	"- Docente: \"Mati no arranca la tarea.\"\n" +
	"  Alizia: \"Probá anticipar la consigna en pasos cortos y darle un arranque concreto. ¿Querés que la guarde como recurso para Mati?\" (sin bloque todavía)\n" +
	"  Docente: \"Sí, guardala.\"\n" +
	"  Alizia: \"Listo, te la dejo guardada.\" + [ADAPTATION_JSON:{...}]\n"

// fewShot is the recommendation example (layer 1, cacheable): situation → actionable
// device recommendation. Used by the recommend framework, which always emits the block.
const fewShot = "EJEMPLO DE RECOMENDACIÓN:\n" +
	"Situación: el alumno no inicia la tarea.\n" +
	"Recomendá un apoyo concreto (p. ej. anticipar la consigna en pasos cortos) y cerrá con el bloque.\n"

// ---------- Layer 2 — dynamic turn context ----------

// AssistSystem builds the in-class assistant system prompt: static framework +
// classroom context (students + kit) + conversational few-shot. Single mode.
func AssistSystem(devices []entities.Device, students []entities.Student) string {
	var b strings.Builder
	b.WriteString(assistFramework)
	b.WriteString("\n")
	b.WriteString(studentRoster(students))
	b.WriteString(deviceCatalog(devices, false))
	b.WriteString("\n")
	b.WriteString(assistFewShot)
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
