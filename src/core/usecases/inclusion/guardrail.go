package inclusion

import (
	"fmt"
	"strconv"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
)

// Guardrail por código (HU-6, §6.7). Antes de mostrar una respuesta al docente,
// validamos por código —no sólo confiando en el prompt— que no cruce límites
// duros: que los DEVICE_ID que menciona existan en el catálogo de la valija y que
// un ADAPTATION_JSON embebido sólo referencie dispositivos reales. Una respuesta
// que falla la validación NUNCA debe llegar al docente: el caller reintenta o cae
// al off-ramp.

// offRampMessage es el mensaje seguro por defecto cuando una respuesta no pasa el
// guardrail y no se pudo reparar. El wording definitivo y su edición se afinan en
// T-6.3 (off-ramp); acá vive como constante editable para no mostrar salida inválida.
const offRampMessage = "Perdón, no pude armar una recomendación válida con la valija en este momento. " +
	"¿Probamos de nuevo describiéndome la situación del alumno con otras palabras?"

// GuardrailResult reporta el veredicto de la validación por código de una
// respuesta generada. Violations describe, en lenguaje accionable para el log,
// qué límite se cruzó (vacío cuando Valid es true).
type GuardrailResult struct {
	Valid      bool
	Violations []string
}

// validateAnswer chequea una respuesta del modelo contra los límites verificables
// por código. validDeviceIDs es el set de ids de dispositivos del catálogo de la
// org (lo que el caller ya cargó para armar el prompt). No re-valida la forma del
// ADAPTATION_JSON: extractAdaptationJSON ya descarta bloques malformados; acá sólo
// nos importa que los device_ids referenciados sean reales.
func validateAnswer(content string, validDeviceIDs map[int64]bool) GuardrailResult {
	var violations []string

	// 1) DEVICE_ID sueltos en el texto ([DEVICE_ID:X]) deben existir en catálogo.
	for _, id := range extractDeviceIDs(content) {
		if !validDeviceIDs[id] {
			violations = append(violations, fmt.Sprintf("DEVICE_ID %d no existe en el catálogo", id))
		}
	}

	// 2) device_ids dentro del ADAPTATION_JSON deben existir en catálogo.
	if adaptation := extractAdaptationJSON(content); adaptation != nil {
		for _, id := range adaptation.DeviceIDs {
			if !validDeviceIDs[id] {
				violations = append(violations, fmt.Sprintf("ADAPTATION_JSON referencia DEVICE_ID %d inexistente", id))
			}
		}
	}

	return GuardrailResult{Valid: len(violations) == 0, Violations: violations}
}

// deviceCatalogSet arma el set de ids de dispositivos válidos a partir del
// catálogo cargado, para alimentar validateAnswer sin re-consultar la DB.
func deviceCatalogSet(devices []entities.Device) map[int64]bool {
	set := make(map[int64]bool, len(devices))
	for i := range devices {
		set[devices[i].ID] = true
	}
	return set
}

// extractDeviceIDs devuelve TODOS los DEVICE_ID referenciados en el texto (a
// diferencia de extractDeviceID, que devuelve sólo el primero), deduplicados.
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
