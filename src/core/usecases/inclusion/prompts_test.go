package inclusion

import (
	"strings"
	"testing"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
)

func ptr(s string) *string { return &s }

func Test_extractDeviceID_ExtractsValidID(t *testing.T) {
	input := "Use this [DEVICE_ID:42] for the student"

	result := extractDeviceID(input)

	if result == nil || *result != 42 {
		t.Errorf("expected 42, got %v", result)
	}
}

func Test_extractDeviceID_ReturnsNilForNoMatch(t *testing.T) {
	input := "no device here"

	result := extractDeviceID(input)

	if result != nil {
		t.Errorf("expected nil, got %v", *result)
	}
}

func Test_extractDeviceID_ReturnsNilForInvalidNumber(t *testing.T) {
	input := "[DEVICE_ID:abc]"

	result := extractDeviceID(input)

	if result != nil {
		t.Errorf("expected nil, got %v", *result)
	}
}

func Test_extractStudentID_ExtractsValidID(t *testing.T) {
	input := "Student [STUDENT_ID:7] needs help"

	result := extractStudentID(input)

	if result == nil || *result != 7 {
		t.Errorf("expected 7, got %v", result)
	}
}

func Test_extractStudentID_ReturnsNilForNoMatch(t *testing.T) {
	input := "no student"

	result := extractStudentID(input)

	if result != nil {
		t.Errorf("expected nil")
	}
}

func Test_extractAdaptationJSON_ExtractsValidJSON(t *testing.T) {
	input := `text [ADAPTATION_JSON:{"title":"Test","type":"actividad_adaptada","strategy":"do thing","device_ids":[1,2],"device_names":["A","B"]}] more text`

	result := extractAdaptationJSON(input)

	if result == nil {
		t.Fatal("expected non-nil")
	}
	if result.Title != "Test" {
		t.Errorf("title = %q, want Test", result.Title)
	}
	if result.Type != "actividad_adaptada" {
		t.Errorf("type = %q", result.Type)
	}
	if len(result.DeviceIDs) != 2 {
		t.Errorf("device_ids len = %d, want 2", len(result.DeviceIDs))
	}
}

func Test_extractAdaptationJSON_ReturnsNilForNoMatch(t *testing.T) {
	input := "no json here"

	result := extractAdaptationJSON(input)

	if result != nil {
		t.Error("expected nil")
	}
}

func Test_extractAdaptationJSON_ReturnsNilForMalformedJSON(t *testing.T) {
	input := "[ADAPTATION_JSON:{invalid json}]"

	result := extractAdaptationJSON(input)

	if result != nil {
		t.Error("expected nil")
	}
}

func Test_buildRecommendSystemPrompt_ContainsDeviceInfo(t *testing.T) {
	devices := []entities.Device{
		{ID: 1, Name: "Timer Visual", NeedsDescription: ptr("Para alumnos con distraccion")},
	}

	prompt := buildRecommendSystemPrompt(devices)

	if !strings.Contains(prompt, "Timer Visual") {
		t.Error("prompt should contain device name")
	}
	if !strings.Contains(prompt, "[ID:1]") {
		t.Error("prompt should contain device ID")
	}
	if !strings.Contains(prompt, "ADAPTATION_JSON") {
		t.Error("prompt should contain format instructions")
	}
	if !strings.Contains(prompt, "Para alumnos con distraccion") {
		t.Error("prompt should contain device needs description")
	}
}
