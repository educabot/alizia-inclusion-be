package inclusion

import (
	"strings"
	"testing"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
)

func ptr(s string) *string { return &s }

func Test_extractDeviceID(t *testing.T) {
	t.Run("extracts valid ID", func(t *testing.T) {
		// Arrange
		input := "Use this [DEVICE_ID:42] for the student"

		// Act
		result := extractDeviceID(input)

		// Assert
		if result == nil || *result != 42 {
			t.Errorf("expected 42, got %v", result)
		}
	})

	t.Run("returns nil for no match", func(t *testing.T) {
		// Arrange
		input := "no device here"

		// Act
		result := extractDeviceID(input)

		// Assert
		if result != nil {
			t.Errorf("expected nil, got %v", *result)
		}
	})

	t.Run("returns nil for invalid number", func(t *testing.T) {
		// Arrange
		input := "[DEVICE_ID:abc]"

		// Act
		result := extractDeviceID(input)

		// Assert
		if result != nil {
			t.Errorf("expected nil, got %v", *result)
		}
	})
}

func Test_extractStudentID(t *testing.T) {
	t.Run("extracts valid ID", func(t *testing.T) {
		// Arrange
		input := "Student [STUDENT_ID:7] needs help"

		// Act
		result := extractStudentID(input)

		// Assert
		if result == nil || *result != 7 {
			t.Errorf("expected 7, got %v", result)
		}
	})

	t.Run("returns nil for no match", func(t *testing.T) {
		// Arrange
		input := "no student"

		// Act
		result := extractStudentID(input)

		// Assert
		if result != nil {
			t.Errorf("expected nil")
		}
	})
}

func Test_extractAdaptationJSON(t *testing.T) {
	t.Run("extracts valid JSON", func(t *testing.T) {
		// Arrange
		input := `text [ADAPTATION_JSON:{"title":"Test","type":"actividad_adaptada","strategy":"do thing","device_ids":[1,2],"device_names":["A","B"]}] more text`

		// Act
		result := extractAdaptationJSON(input)

		// Assert
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
	})

	t.Run("returns nil for no match", func(t *testing.T) {
		// Arrange
		input := "no json here"

		// Act
		result := extractAdaptationJSON(input)

		// Assert
		if result != nil {
			t.Error("expected nil")
		}
	})

	t.Run("returns nil for malformed JSON", func(t *testing.T) {
		// Arrange
		input := "[ADAPTATION_JSON:{invalid json}]"

		// Act
		result := extractAdaptationJSON(input)

		// Assert
		if result != nil {
			t.Error("expected nil")
		}
	})
}

func Test_buildRecommendSystemPrompt(t *testing.T) {
	// Arrange
	devices := []entities.Device{
		{ID: 1, Name: "Timer Visual", NeedsDescription: ptr("Para alumnos con distraccion")},
	}

	// Act
	prompt := buildRecommendSystemPrompt(devices)

	// Assert
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
