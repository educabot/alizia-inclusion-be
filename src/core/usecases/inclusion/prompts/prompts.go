// Package prompts es la capa editable del marco pedagógico de Alizia (HU-6, T-6.1).
// Separa el motor (usecases) del texto de los prompts, en dos capas:
//
//   - Capa 1 (estática, cacheable): rol, lineamientos, límites, formato de salida
//     y few-shot. Son constantes; cambiar el texto no rompe el formato de salida.
//   - Capa 2 (dinámica): el contexto del turno (alumnos del aula, catálogo de la
//     valija) que se arma por request.
//
// Hay un solo modo: no se elige entre prompts según un "mode"; Alizia se adapta
// por el contexto.
package prompts

import (
	"fmt"
	"strings"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
)

// ---------- Capa 1 — marco estático (cacheable) ----------

// outputRules fija el formato de salida exigido por HU-6 (T-6.6): pocas acciones,
// con niveles de diferenciación, accionables al instante y en tono rioplatense.
const outputRules = "FORMATO DE RESPUESTA (obligatorio):\n" +
	"- Ofrecé 1 a 3 acciones concretas, ordenadas por impacto.\n" +
	"- Incluí al menos 3 niveles de diferenciación (más simple / intermedio / más desafiante).\n" +
	"- Aplicable en menos de 1 minuto de lectura: breve y al grano.\n" +
	"- Español rioplatense, tono cálido, sin jerga clínica.\n"

// scopeRules fija los límites duros del marco (HU-6): nunca diagnostica ni
// reemplaza al docente; el paso al costado es el último recurso (off-ramp, T-6.3).
const scopeRules = "LÍMITES (marco pedagógico):\n" +
	"- Entrada pedagógica, no clínica: partís de situaciones de aula, no de diagnósticos.\n" +
	"- Nunca diagnostiques, no reemplaces al docente, no produzcas informes clínicos.\n" +
	"- Si el caso se va de tu alcance (clínico/crisis), primero ayudá con lo que tengas;\n" +
	"  dar un paso al costado es el último recurso, no el primero.\n"

// adaptationBlock instruye el bloque estructurado que el guardrail valida por código.
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

// fewShot es un ejemplo estático curado (capa 1 cacheable) de situación de aula →
// respuesta con acciones + niveles. El few-shot dinámico por alumno (a partir de
// adaptaciones previas que funcionaron) es una extensión futura sobre esta base.
const fewShot = "EJEMPLO DE BUENA RESPUESTA:\n" +
	"Situación: el alumno no inicia la tarea.\n" +
	"1. Anticipá la consigna en pasos cortos en el pizarrón.\n" +
	"   - Más simple: un paso por vez. Intermedio: 2-3 pasos. Más desafiante: que arme la secuencia él mismo.\n" +
	"2. Dale un arranque concreto (\"empezá por...\") y un tiempo breve para el primer paso.\n"

// ---------- Capa 2 — contexto dinámico del turno ----------

// AssistSystem arma el system prompt del asistente en clase: marco estático +
// contexto del aula (alumnos + valija) + few-shot. Un solo modo (HU-6, T-6.1).
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

// RecommendSystem arma el system prompt de recomendación de dispositivos: marco +
// catálogo detallado (fundamento / uso) + few-shot.
func RecommendSystem(devices []entities.Device) string {
	var b strings.Builder
	b.WriteString(recommendFramework)
	b.WriteString("\n")
	b.WriteString(deviceCatalog(devices, true))
	b.WriteString("\n")
	b.WriteString(fewShot)
	return b.String()
}

// studentRoster lista los alumnos del aula con su id y dificultades (capa 2).
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

// deviceCatalog lista la valija (capa 2). Con detailed agrega fundamento y uso,
// para el prompt de recomendación; sin él, solo nombre y para qué sirve.
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
