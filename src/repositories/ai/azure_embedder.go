package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

// AzureEmbedder implementa providers.Embedder contra el endpoint de embeddings de
// Azure OpenAI. Usa su propio recurso (endpoint/key/deployment) porque puede no ser
// el mismo que el de chat. Reusa azureError de azure_client.go (mismo paquete).
type AzureEmbedder struct {
	endpoint   string
	apiKey     string
	deployment string
	apiVersion string
	dimensions int
	httpClient *http.Client
}

func NewAzureEmbedder(endpoint, apiKey, deployment, apiVersion string, dimensions int) providers.Embedder {
	return &AzureEmbedder{
		endpoint:   strings.TrimRight(endpoint, "/"),
		apiKey:     apiKey,
		deployment: deployment,
		apiVersion: apiVersion,
		dimensions: dimensions,
		httpClient: &http.Client{},
	}
}

func (e *AzureEmbedder) EmbedQuery(ctx context.Context, text string) ([]float32, error) {
	url := fmt.Sprintf("%s/openai/deployments/%s/embeddings?api-version=%s", e.endpoint, e.deployment, e.apiVersion)

	body := map[string]any{"input": text}
	if e.dimensions > 0 {
		body["dimensions"] = e.dimensions
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal embedding request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("create embedding request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", e.apiKey)

	start := time.Now()
	resp, err := e.httpClient.Do(req)
	if err != nil {
		slog.ErrorContext(ctx, "embed.request", "deployment", e.deployment, "chars", len([]rune(text)), "error", err.Error())
		return nil, fmt.Errorf("do embedding request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read embedding response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("azure embeddings error (status %d): %s", resp.StatusCode, string(raw))
	}

	var parsed struct {
		Data []struct {
			Embedding []float32 `json:"embedding"`
		} `json:"data"`
		Error *azureError `json:"error,omitempty"`
	}
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, fmt.Errorf("unmarshal embedding response: %w", err)
	}
	if parsed.Error != nil {
		return nil, fmt.Errorf("azure embeddings: %s (code: %s)", parsed.Error.Message, parsed.Error.Code)
	}
	if len(parsed.Data) == 0 || len(parsed.Data[0].Embedding) == 0 {
		return nil, fmt.Errorf("azure embeddings: empty embedding in response")
	}
	slog.InfoContext(ctx, "embed.request",
		"deployment", e.deployment,
		"chars", len([]rune(text)),
		"dims", len(parsed.Data[0].Embedding),
		"duration_ms", time.Since(start).Milliseconds(),
	)
	return parsed.Data[0].Embedding, nil
}
