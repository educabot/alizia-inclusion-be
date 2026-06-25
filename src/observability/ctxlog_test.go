package observability

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContextHandler_InjectsCorrelationIDs(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	logger := slog.New(NewContextHandler(slog.NewJSONHandler(&buf, nil)))
	org := uuid.New()
	ctx := WithUser(WithOrg(WithRequestID(context.Background(), "req-123"), org), int64(7))

	// Act
	logger.InfoContext(ctx, "hello")

	// Assert
	var rec map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &rec))
	assert.Equal(t, "req-123", rec["request_id"])
	assert.Equal(t, org.String(), rec["org_id"])
	assert.Equal(t, float64(7), rec["user_id"])
}

func TestContextHandler_SkipsMissingIDs(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(NewContextHandler(slog.NewJSONHandler(&buf, nil)))

	logger.InfoContext(context.Background(), "hello")

	var rec map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &rec))
	_, hasReq := rec["request_id"]
	assert.False(t, hasReq)
}

func TestText_VerboseTogglesFullTextVsLength(t *testing.T) {
	t.Cleanup(func() { SetVerbose(true) })

	SetVerbose(true)
	a := Text("prompt", "hola mundo")
	assert.Equal(t, "prompt", a.Key)
	assert.Equal(t, "hola mundo", a.Value.String())

	SetVerbose(false)
	a = Text("prompt", "hola mundo")
	assert.Equal(t, "prompt_len", a.Key)
	assert.Equal(t, int64(10), a.Value.Int64())
}
