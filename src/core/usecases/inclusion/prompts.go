package inclusion

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
)

type GeneratedAdaptation struct {
	Title       string   `json:"title"`
	Type        string   `json:"type"`
	Strategy    string   `json:"strategy"`
	DeviceIDs   []int64  `json:"device_ids"`
	DeviceNames []string `json:"device_names"`
	// StudentID is the student the adaptation is for, taken from the [STUDENT_ID:X]
	// tag of the same turn. Surfaced here so the frontend can build the save request
	// from a single object instead of parsing the response text.
	StudentID *int64 `json:"student_id,omitempty"`
}

func buildRecommendUserPrompt(student *entities.Student, req RecommendDeviceRequest) string {
	var b strings.Builder

	fmt.Fprintf(&b, "Asignatura: %s\n", req.Subject)
	if req.Objective != "" {
		fmt.Fprintf(&b, "Objetivo de la clase: %s\n", req.Objective)
	}
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

var deviceIDRegex = regexp.MustCompile(`\[DEVICE_ID:(\d+)\]`)
var studentIDRegex = regexp.MustCompile(`\[STUDENT_ID:(\d+)\]`)
var contentIDRegex = regexp.MustCompile(`\[CONTENT_ID:(\d+)\]`)
var adaptationJSONRegex = regexp.MustCompile(`\[ADAPTATION_JSON:(\{.+\})\]`)

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

// extractContentIDs returns the pedagogical content ids the model cited via
// [CONTENT_ID:X] tags, in order and without duplicates. Lets the frontend deep-link
// a material chip to the specific document instead of the materials list.
func extractContentIDs(content string) []int64 {
	matches := contentIDRegex.FindAllStringSubmatch(content, -1)
	if len(matches) == 0 {
		return nil
	}
	seen := make(map[int64]struct{}, len(matches))
	ids := make([]int64, 0, len(matches))
	for _, m := range matches {
		id, err := strconv.ParseInt(m[1], 10, 64)
		if err != nil {
			continue
		}
		if _, dup := seen[id]; dup {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	return ids
}

func extractAdaptationJSON(content string) *GeneratedAdaptation {
	matches := adaptationJSONRegex.FindStringSubmatch(content)
	if len(matches) < 2 {
		return nil
	}
	var adaptation GeneratedAdaptation
	if err := json.Unmarshal([]byte(matches[1]), &adaptation); err != nil {
		return nil
	}
	return &adaptation
}
