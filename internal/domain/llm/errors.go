package llm

import "errors"

var (
	ErrProviderNotFound   = errors.New("llm provider not registered or supported")
	ErrInvalidAPIKey      = errors.New("invalid or missing API key")
	ErrRateLimitExceeded  = errors.New("rate limit exceeded for provider")
	ErrContextTooLong     = errors.New("context length exceeded model token limit")
	ErrProviderTimeout    = errors.New("upstream provider request timed out")
	ErrStreamingFailed    = errors.New("streaming connection interrupted")
)
