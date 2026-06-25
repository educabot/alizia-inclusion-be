package inclusion

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
)

func TestValidateAnswer_AllowsAnswerWithoutDeviceRefs(t *testing.T) {
	// Arrange
	valid := map[int64]bool{1: true}

	// Act
	got := validateAnswer("Probá con pausas sensoriales y consignas cortas.", valid)

	// Assert
	assert.True(t, got.Valid)
	assert.Empty(t, got.Violations)
}

func TestValidateAnswer_AllowsExistingDeviceID(t *testing.T) {
	valid := map[int64]bool{7: true, 8: true}

	got := validateAnswer("Te sugiero el tablero [DEVICE_ID:7] para anticipar la rutina.", valid)

	assert.True(t, got.Valid)
}

func TestValidateAnswer_RejectsUnknownDeviceID(t *testing.T) {
	valid := map[int64]bool{7: true}

	got := validateAnswer("Usá [DEVICE_ID:99] que no existe.", valid)

	assert.False(t, got.Valid)
	assert.Len(t, got.Violations, 1)
	assert.Contains(t, got.Violations[0], "99")
}

func TestValidateAnswer_RejectsUnknownDeviceInAdaptationJSON(t *testing.T) {
	valid := map[int64]bool{1: true}
	content := `Listo: [ADAPTATION_JSON:{"title":"t","type":"actividad_adaptada","strategy":"s","device_ids":[1,42],"device_names":["a","b"]}]`

	got := validateAnswer(content, valid)

	assert.False(t, got.Valid)
	assert.Len(t, got.Violations, 1)
	assert.Contains(t, got.Violations[0], "42")
}

func TestValidateAnswer_AcceptsAdaptationJSONWithOnlyRealDevices(t *testing.T) {
	valid := map[int64]bool{1: true, 2: true}
	content := `[ADAPTATION_JSON:{"title":"t","type":"actividad_adaptada","strategy":"s","device_ids":[1,2],"device_names":["a","b"]}]`

	got := validateAnswer(content, valid)

	assert.True(t, got.Valid)
}

func TestValidateAnswer_ReportsBothLooseAndJSONViolations(t *testing.T) {
	valid := map[int64]bool{1: true}
	content := `Mirá [DEVICE_ID:50] y guardá [ADAPTATION_JSON:{"title":"t","type":"x","strategy":"s","device_ids":[60]}]`

	got := validateAnswer(content, valid)

	assert.False(t, got.Valid)
	assert.Len(t, got.Violations, 2)
}

func TestExtractDeviceIDs_DedupesAndParsesAll(t *testing.T) {
	got := extractDeviceIDs("usar [DEVICE_ID:3], y de nuevo [DEVICE_ID:3] más [DEVICE_ID:5]")

	assert.Equal(t, []int64{3, 5}, got)
}

func TestExtractDeviceIDs_NoneWhenAbsent(t *testing.T) {
	assert.Nil(t, extractDeviceIDs("sin dispositivos acá"))
}

func TestDeviceCatalogSet_BuildsIDSet(t *testing.T) {
	devices := []entities.Device{{ID: 4}, {ID: 9}}

	set := deviceCatalogSet(devices)

	assert.True(t, set[4])
	assert.True(t, set[9])
	assert.False(t, set[1])
}
