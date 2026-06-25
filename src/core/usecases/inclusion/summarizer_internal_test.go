package inclusion

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/educabot/alizia-inclusion-be/src/core/usecases/inclusion/prompts"
)

func TestSummarizerSystemPrompt_ReusesSharedRol(t *testing.T) {
	// The summarizer is the third identity surface: it must reuse the shared role
	// (prompts.RolAlizia) instead of hardcoding its own "Sos Alizia" line.
	assert.True(t, strings.HasPrefix(summarizerSystemPrompt, prompts.RolAlizia),
		"el summarizer debe arrancar con la identidad compartida RolAlizia")
}

func TestSummarizerSystemPrompt_KeepsJSONContract(t *testing.T) {
	// Reusing the role must not leak conversational voice rules that would break the
	// JSON-only output contract.
	assert.Contains(t, summarizerSystemPrompt, "EXCLUSIVAMENTE un JSON")
	assert.NotContains(t, summarizerSystemPrompt, "una sola pregunta")
}
