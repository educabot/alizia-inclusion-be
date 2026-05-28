package ai_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ai "github.com/educabot/alizia-inclusion-be/src/repositories/ai"

	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

func TestStubClient_Generate_ReturnsStubContent(t *testing.T) {
	client := ai.NewStubClient()
	got, err := client.Generate(context.Background(), providers.GenerateParams{
		UserPrompt: "test prompt",
	})

	require.NoError(t, err)
	assert.NotEmpty(t, got)
	assert.Contains(t, got, "[stub]")
}

func TestStubClient_Chat_ReturnsStubResponse(t *testing.T) {
	client := ai.NewStubClient()
	resp, err := client.Chat(context.Background(), []providers.ChatMessage{
		{Role: "user", Content: "hello"},
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.NotEmpty(t, resp.Content)
}

func TestStubClient_ChatWithTools_ReturnsStubResponse(t *testing.T) {
	client := ai.NewStubClient()
	resp, err := client.ChatWithTools(
		context.Background(),
		[]providers.ChatMessage{{Role: "user", Content: "use tool"}},
		[]providers.ToolDefinition{{Name: "noop", Description: "does nothing"}},
	)

	require.NoError(t, err)
	require.NotNil(t, resp)
}
