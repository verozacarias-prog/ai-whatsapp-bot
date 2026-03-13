package main

import (
	"errors"
	"log"
	"time"

	"github.com/openai/openai-go/v3"
)

func isRetryable(err error) bool {
	var apierr *openai.Error
	if !errors.As(err, &apierr) {
		return false
	}

	// quota agotada — no reintentar
	if apierr.StatusCode == 429 && apierr.Code == "insufficient_quota" {
		log.Printf("OpenAI quota exceeded — check billing: %v", err)
		return false
	}

	// rate limit por RPM o errores de servidor — reintentar
	return apierr.StatusCode == 429 || apierr.StatusCode == 500 || apierr.StatusCode == 503
}

func withRetry(maxAttempts int, fn func() error) error {
	var err error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		err = fn()
		if err == nil {
			return nil
		}

		if !isRetryable(err) {
			return err
		}

		waitSeconds := time.Duration(1<<(attempt-1)) * time.Second
		log.Printf("attempt %d/%d failed, retrying in %v: %v", attempt, maxAttempts, waitSeconds, err)
		time.Sleep(waitSeconds)
	}
	return err
}
