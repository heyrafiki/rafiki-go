// SPDX-License-Identifier: Apache-2.0

package heyrafiki

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"
)

// RetryPolicy controls bounded retries for reads and writes protected by an idempotency key.
type RetryPolicy struct {
	MaxAttempts int
	BaseDelay   time.Duration
	MaxDelay    time.Duration
}

// DefaultRetryPolicy returns the SDK's bounded exponential-backoff policy.
func DefaultRetryPolicy() RetryPolicy {
	return RetryPolicy{
		MaxAttempts: 3,
		BaseDelay:   200 * time.Millisecond,
		MaxDelay:    2 * time.Second,
	}
}

func (policy RetryPolicy) validate() error {
	if policy.MaxAttempts < 1 || policy.MaxAttempts > 10 {
		return errors.New("heyrafiki: retry attempts must be between 1 and 10")
	}
	if policy.BaseDelay < 0 || policy.MaxDelay < 0 || policy.BaseDelay > policy.MaxDelay {
		return errors.New("heyrafiki: retry delays are invalid")
	}
	return nil
}

func retryableStatus(status int) bool {
	switch status {
	case http.StatusRequestTimeout,
		http.StatusTooManyRequests,
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:
		return true
	default:
		return false
	}
}

func retryableTransportError(err error) bool {
	var netErr net.Error
	return errors.As(err, &netErr)
}

func (client *Client) waitForRetry(ctx context.Context, attempt int, retryAfter time.Duration) error {
	delay := client.retryPolicy.BaseDelay
	for step := 1; step < attempt && delay < client.retryPolicy.MaxDelay; step++ {
		delay *= 2
		if delay > client.retryPolicy.MaxDelay {
			delay = client.retryPolicy.MaxDelay
		}
	}
	if retryAfter > 0 {
		delay = retryAfter
		if delay > client.retryPolicy.MaxDelay {
			delay = client.retryPolicy.MaxDelay
		}
	}
	if delay <= 0 {
		return nil
	}
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
