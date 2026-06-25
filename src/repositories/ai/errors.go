package ai

import (
	"errors"
	"fmt"
)

// ErrEmptyResponse indicates the upstream returned a payload with no completion
// choices, so there is no content to return.
var ErrEmptyResponse = errors.New("azure generate: empty response")

// UpstreamStatusError reports a non-200 HTTP status returned by Azure OpenAI.
// It carries the status code and raw body so callers can branch on the code via
// errors.As instead of matching substrings of the message.
type UpstreamStatusError struct {
	StatusCode int
	Body       string
}

func (e *UpstreamStatusError) Error() string {
	return fmt.Sprintf("azure openai error (status %d): %s", e.StatusCode, e.Body)
}
