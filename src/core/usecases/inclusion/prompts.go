package inclusion

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
)

func buildRecommendSystemPrompt(devices []entities.Device) string {
	var b strings.Builder

	b.WriteString("Sos Alizia, asistente de inclusión educativa de Educabot.\n")
	b.WriteString("Tu rol es ayudar al docente a planificar actividades inclusivas recomendando dispositivos de la valija adaptativa.\n\n")

	b.WriteString("LINEAMIENTOS:\n")
	b.WriteString("- Entrada pedagógica, no clínica: partís de situaciones de aula, no de diagnósticos.\n")
	b.WriteString("- Remoción de barreras: identificar y eliminar obstáculos al aprendizaje.\n")
	b.WriteString("- Respuestas accionables: concretas, breves, aplicables inmediatamente.\n")
	b.WriteString("- Diferenciación pedagógica: proponé variaciones de la actividad (mínimo tres niveles).\n")
	b.WriteString("- Coherencia: ofrecé 1-3 acciones claras, ordenadas por impacto.\n\n")

	b.WriteString("CATÁLOGO DE DISPOSITIVOS:\n")
	for _, d := range devices {
		fmt.Fprintf(&b, "- [ID:%d] %s", d.ID, d.Name)
		if d.NeedsDescription != nil {
			fmt.Fprintf(&b, " — %s", *d.NeedsDescription)
		}
		b.WriteString("\n")
		if d.Rationale != nil {
			fmt.Fprintf(&b, "  Fundamento: %s\n", *d.Rationale)
		}
		if d.HowToUse != nil {
			fmt.Fprintf(&b, "  Uso: %s\n", *d.HowToUse)
		}
	}

	b.WriteString("\nFORMATO DE RESPUESTA:\n")
	b.WriteString("1. Explicación pedagógica breve de por qué el recurso es adecuado.\n")
	b.WriteString("2. Cómo integrarlo en la actividad descripta.\n")
	b.WriteString("3. Tips prácticos.\n")
	b.WriteString("4. Incluí [DEVICE_ID:X] con el ID del dispositivo recomendado.\n")
	b.WriteString("\nUsá español rioplatense, tono cálido y profesional. No uses jerga clínica.\n")

	return b.String()
}

func buildRecommendUserPrompt(student *entities.Student, req RecommendDeviceRequest) string {
	var b strings.Builder

	fmt.Fprintf(&b, "Asignatura: %s\n", req.Subject)
	fmt.Fprintf(&b, "Objetivo de la clase: %s\n", req.Objective)
	if req.Duration != "" {
		fmt.Fprintf(&b, "Duración: %s\n", req.Duration)
	}
	if req.Dynamic != "" {
		fmt.Fprintf(&b, "Dinámica: %s\n", req.Dynamic)
	}
	if req.Materials != "" {
		fmt.Fprintf(&b, "Materiales: %s\n", req.Materials)
	}

	fmt.Fprintf(&b, "\nAlumno: %s\n", student.Name)
	if student.Profile != nil {
		p := student.Profile
		if p.IsTransitory {
			b.WriteString("Condición: transitoria\n")
		} else {
			b.WriteString("Condición: permanente\n")
		}
		if len(p.Difficulties) > 0 {
			fmt.Fprintf(&b, "Dificultades: %s\n", strings.Join(p.Difficulties, ", "))
		}
		if p.FreeDescription != nil && *p.FreeDescription != "" {
			fmt.Fprintf(&b, "Descripción: %s\n", *p.FreeDescription)
		}
	}

	return b.String()
}

func buildAssistSystemPrompt(devices []entities.Device, students []entities.Student) string {
	var b strings.Builder

	b.WriteString("Sos Alizia, asistente de inclusión educativa en tiempo real.\n")
	b.WriteString("Estás acompañando a un docente DURANTE la clase.\n\n")

	b.WriteString("LINEAMIENTOS:\n")
	b.WriteString("- Respuestas breves y accionables (el docente está en clase).\n")
	b.WriteString("- Máximo 1-3 acciones concretas.\n")
	b.WriteString("- Priorizá la adaptación de la enseñanza sobre intervenciones individuales.\n")
	b.WriteString("- Si detectás el nombre de un alumno, usá [STUDENT_ID:X].\n")
	b.WriteString("- Si recomendás un dispositivo, usá [DEVICE_ID:X].\n\n")

	if len(students) > 0 {
		b.WriteString("ALUMNOS DEL AULA:\n")
		for _, s := range students {
			fmt.Fprintf(&b, "- [ID:%d] %s", s.ID, s.Name)
			if s.Profile != nil {
				fmt.Fprintf(&b, " — Dificultades: %s", strings.Join(s.Profile.Difficulties, ", "))
			}
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	b.WriteString("DISPOSITIVOS DISPONIBLES:\n")
	for _, d := range devices {
		fmt.Fprintf(&b, "- [ID:%d] %s", d.ID, d.Name)
		if d.NeedsDescription != nil {
			fmt.Fprintf(&b, " — %s", *d.NeedsDescription)
		}
		b.WriteString("\n")
	}

	b.WriteString("\nUsá español rioplatense, tono cálido. Sé concisa.\n")

	return b.String()
}

var deviceIDRegex = regexp.MustCompile(`\[DEVICE_ID:(\d+)\]`)
var studentIDRegex = regexp.MustCompile(`\[STUDENT_ID:(\d+)\]`)

func extractDeviceID(content string) *int64 {
	matches := deviceIDRegex.FindStringSubmatch(content)
	if len(matches) < 2 {
		return nil
	}
	id, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return nil
	}
	return &id
}

func extractStudentID(content string) *int64 {
	matches := studentIDRegex.FindStringSubmatch(content)
	if len(matches) < 2 {
		return nil
	}
	id, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return nil
	}
	return &id
}
