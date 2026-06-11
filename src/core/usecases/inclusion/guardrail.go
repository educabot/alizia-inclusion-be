package inclusion

import (
	"fmt"
	"strconv"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
)

// Code-level guardrail (HU-6, §6.7). Before surfacing a response to the teacher,
// we enforce hard limits in code — not solely via prompt — ensuring that any
// DEVICE_ID mentioned exists in the org's device catalog and that an embedded
// ADAPTATION_JSON only references real devices. A response that fails validation
// must never reach the teacher: the caller retries or falls through to the off-ramp.

// GuardrailResult reports the code-validation verdict for a generated response.
// Violations lists each hard-limit breach in actionable log language; empty when Valid is true.
type GuardrailResult struct {
	Valid      bool
	Violations []string
}

// validateAnswer checks a model response against code-verifiable hard limits.
// validDeviceIDs is the org's device-catalog set (already loaded by the caller for
// prompt assembly). It does not re-validate ADAPTATION_JSON shape — extractAdaptationJSON
// already discards malformed blocks; here we only verify that referenced device_ids are real.
func validateAnswer(content string, validDeviceIDs map[int64]bool) GuardrailResult {
	var violations []string

	// 1) Inline DEVICE_ID tokens ([DEVICE_ID:X]) must exist in the catalog.
	for _, id := range extractDeviceIDs(content) {
		if !validDeviceIDs[id] {
			violations = append(violations, fmt.Sprintf("DEVICE_ID %d no existe en el catálogo", id))
		}
	}

	// 2) device_ids inside ADAPTATION_JSON must exist in the catalog.
	if adaptation := extractAdaptationJSON(content); adaptation != nil {
		for _, id := range adaptation.DeviceIDs {
			if !validDeviceIDs[id] {
				violations = append(violations, fmt.Sprintf("ADAPTATION_JSON referencia DEVICE_ID %d inexistente", id))
			}
		}
	}

	return GuardrailResult{Valid: len(violations) == 0, Violations: violations}
}

// deviceCatalogSet builds the valid device-ID set from the loaded catalog,
// so validateAnswer can check membership without an extra DB query.
func deviceCatalogSet(devices []entities.Device) map[int64]bool {
	set := make(map[int64]bool, len(devices))
	for i := range devices {
		set[devices[i].ID] = true
	}
	return set
}

// extractDeviceIDs returns all DEVICE_ID values referenced in the text,
// deduplicated. Unlike extractDeviceID, which returns only the first match.
func extractDeviceIDs(content string) []int64 {
	matches := deviceIDRegex.FindAllStringSubmatch(content, -1)
	if len(matches) == 0 {
		return nil
	}
	seen := make(map[int64]struct{}, len(matches))
	out := make([]int64, 0, len(matches))
	for _, m := range matches {
		if len(m) < 2 {
			continue
		}
		id, err := strconv.ParseInt(m[1], 10, 64)
		if err != nil {
			continue
		}
		if _, dup := seen[id]; dup {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}
