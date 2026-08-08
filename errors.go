// SPDX-License-Identifier: Apache-2.0

package heyrafiki

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// APIError represents the documented Heyrafiki API error envelope.
type APIError struct {
	StatusCode int
	Code       string
	Message    string
	RequestID  string
	Docs       string
	RetryAfter time.Duration
	cause      error
}

func (err *APIError) Error() string {
	if err.RequestID != "" {
		return fmt.Sprintf("heyrafiki: %s (status %d, request %s)", err.Message, err.StatusCode, err.RequestID)
	}
	return fmt.Sprintf("heyrafiki: %s (status %d)", err.Message, err.StatusCode)
}

// Unwrap exposes JSON decoding failures without including response bodies.
func (err *APIError) Unwrap() error { return err.cause }

type errorEnvelope struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Docs    string `json:"docs"`
	} `json:"error"`
}

func parseAPIError(response *http.Response, payload []byte) *APIError {
	apiErr := &APIError{
		StatusCode: response.StatusCode,
		Code:       "request_failed",
		Message:    "The Heyrafiki API request failed.",
		RequestID:  response.Header.Get("X-Request-Id"),
		RetryAfter: parseRetryAfter(response.Header.Get("Retry-After"), time.Now()),
	}
	var envelope errorEnvelope
	if json.Unmarshal(payload, &envelope) == nil {
		if envelope.Error.Code != "" {
			apiErr.Code = envelope.Error.Code
		}
		if envelope.Error.Message != "" {
			apiErr.Message = envelope.Error.Message
		}
		apiErr.Docs = envelope.Error.Docs
	}
	return apiErr
}

func parseRetryAfter(value string, now time.Time) time.Duration {
	if seconds, err := strconv.Atoi(value); err == nil && seconds >= 0 {
		return time.Duration(seconds) * time.Second
	}
	if at, err := http.ParseTime(value); err == nil && at.After(now) {
		return at.Sub(now)
	}
	return 0
}
