package inclusion

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
)

func ptr(s string) *string { return &s }

func Test_extractDeviceID_ExtractsValidID(t *testing.T) {
	input := "Use this [DEVICE_ID:42] for the student"

	result := extractDeviceID(input)

	require.NotNil(t, result)
	assert.Equal(t, int64(42), *result)
}

func Test_extractDeviceID_ReturnsNilForNoMatch(t *testing.T) {
	input := "no device here"

	result := extractDeviceID(input)

	assert.Nil(t, result)
}

func Test_extractDeviceID_ReturnsNilForInvalidNumber(t *testing.T) {
	input := "[DEVICE_ID:abc]"

	result := extractDeviceID(input)

	assert.Nil(t, result)
}

func Test_extractStudentID_ExtractsValidID(t *testing.T) {
	input := "Student [STUDENT_ID:7] needs help"

	result := extractStudentID(input)

	require.NotNil(t, result)
	assert.Equal(t, int64(7), *result)
}

func Test_extractStudentID_ReturnsNilForNoMatch(t *testing.T) {
	input := "no student"

	result := extractStudentID(input)

	assert.Nil(t, result)
}

func Test_extractAdaptationJSON_ExtractsValidJSON(t *testing.T) {
	input := `text [ADAPTATION_JSON:{"title":"Test","type":"actividad_adaptada","strategy":"do thing","device_ids":[1,2],"device_names":["A","B"]}] more text`

	result := extractAdaptationJSON(input)

	require.NotNil(t, result)
	assert.Equal(t, "Test", result.Title)
	assert.Equal(t, "actividad_adaptada", result.Type)
	assert.Len(t, result.DeviceIDs, 2)
}

func Test_extractAdaptationJSON_ReturnsNilForNoMatch(t *testing.T) {
	input := "no json here"

	result := extractAdaptationJSON(input)

	assert.Nil(t, result)
}

func Test_extractAdaptationJSON_ReturnsNilForMalformedJSON(t *testing.T) {
	input := "[ADAPTATION_JSON:{invalid json}]"

	result := extractAdaptationJSON(input)

	assert.Nil(t, result)
}

func Test_buildRecommendSystemPrompt_ContainsDeviceInfo(t *testing.T) {
	devices := []entities.Device{
		{ID: 1, Name: "Timer Visual", NeedsDescription: ptr("Para alumnos con distraccion")},
	}

	prompt := buildRecommendSystemPrompt(devices)

	assert.True(t, strings.Contains(prompt, "Timer Visual"), "prompt should contain device name")
	assert.True(t, strings.Contains(prompt, "[ID:1]"), "prompt should contain device ID")
	assert.True(t, strings.Contains(prompt, "ADAPTATION_JSON"), "prompt should contain format instructions")
	assert.True(t, strings.Contains(prompt, "Para alumnos con distraccion"), "prompt should contain device needs description")
}
