// SPDX-License-Identifier: Apache-2.0

package heyrafiki

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (function roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return function(request)
}

func testClient(t *testing.T, handler http.Handler) *Client {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	client, err := NewClient(
		"test_api_key",
		WithBaseURL(server.URL+"/v1"),
		WithRetryPolicy(RetryPolicy{MaxAttempts: 1}),
	)
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func TestPublishedOperationSurface(t *testing.T) {
	var mu sync.Mutex
	var received []string
	client := testClient(t, http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		mu.Lock()
		received = append(received, request.Method+" "+request.URL.Path)
		mu.Unlock()
		response.Header().Set("Content-Type", "application/json")
		_, _ = response.Write([]byte(`{}`))
	}))

	ctx := context.Background()
	write := WriteOptions{IdempotencyKey: "request-0001"}
	calls := []func() error{
		func() error { _, err := client.API.Retrieve(ctx); return err },
		func() error { _, err := client.Practitioners.List(ctx, nil); return err },
		func() error { _, err := client.Practitioners.Retrieve(ctx, "prc_1"); return err },
		func() error { _, err := client.Practitioners.Availability(ctx, "prc_1"); return err },
		func() error { _, err := client.Bookings.List(ctx, nil); return err },
		func() error { _, err := client.Bookings.Create(ctx, BookingInput{}, write); return err },
		func() error { _, err := client.Bookings.Retrieve(ctx, "bkg_1"); return err },
		func() error { _, err := client.Sessions.List(ctx, nil); return err },
		func() error { _, err := client.Sessions.Retrieve(ctx, "ses_1"); return err },
		func() error {
			_, err := client.EligibilityChecks.Create(ctx, EligibilityCheckInput{}, write)
			return err
		},
		func() error { _, err := client.EligibilityChecks.Retrieve(ctx, "elig_1"); return err },
		func() error { _, err := client.Coverages.Record(ctx, CoverageObservationInput{}, write); return err },
		func() error {
			_, err := client.CoverageBatches.Record(ctx, CoverageBatchInput{}, CoverageBatchOptions{
				IdempotencyKey: "request-0001", ArtifactReference: "payer:file:1",
			})
			return err
		},
		func() error {
			_, err := client.Preauthorizations.Create(ctx, PreauthorizationInput{}, write)
			return err
		},
		func() error { _, err := client.Preauthorizations.Retrieve(ctx, "preauth_1"); return err },
		func() error {
			_, err := client.Preauthorizations.Decide(ctx, "preauth_1", DeniedPreauthorizationDecisionInput{}, write)
			return err
		},
		func() error { _, err := client.Claims.List(ctx, nil); return err },
		func() error { _, err := client.Claims.Create(ctx, ClaimInput{}, write); return err },
		func() error { _, err := client.Claims.Retrieve(ctx, "clm_1"); return err },
		func() error {
			_, err := client.Claims.RequestInformation(ctx, "clm_1", ClaimInformationRequestInput{}, write)
			return err
		},
		func() error {
			_, err := client.Claims.SubmitEvidence(ctx, "clm_1", ClaimEvidenceInput{}, write)
			return err
		},
		func() error {
			_, err := client.Claims.Adjudicate(ctx, "clm_1", ClaimAdjudicationInput{}, write)
			return err
		},
		func() error { _, err := client.Remittances.List(ctx, nil); return err },
		func() error { _, err := client.Remittances.Create(ctx, RemittanceInput{}, write); return err },
		func() error { _, err := client.Remittances.Retrieve(ctx, "rem_1"); return err },
		func() error { _, err := client.WebhookEndpoints.List(ctx, nil); return err },
		func() error { _, err := client.WebhookEndpoints.Create(ctx, WebhookEndpointInput{}); return err },
		func() error { _, err := client.WebhookEndpoints.Retrieve(ctx, "whe_1"); return err },
		func() error { _, err := client.WebhookEndpoints.Disable(ctx, "whe_1"); return err },
		func() error { _, err := client.WebhookEndpoints.SendTest(ctx, "whe_1"); return err },
	}
	for index, call := range calls {
		if err := call(); err != nil {
			t.Fatalf("operation %d failed: %v", index, err)
		}
	}

	expected := []string{
		"GET /v1/",
		"GET /v1/practitioners",
		"GET /v1/practitioners/prc_1",
		"GET /v1/practitioners/prc_1/availability",
		"GET /v1/bookings",
		"POST /v1/bookings",
		"GET /v1/bookings/bkg_1",
		"GET /v1/sessions",
		"GET /v1/sessions/ses_1",
		"POST /v1/eligibility_checks",
		"GET /v1/eligibility_checks/elig_1",
		"POST /v1/coverages",
		"POST /v1/coverage_batches",
		"POST /v1/preauthorizations",
		"GET /v1/preauthorizations/preauth_1",
		"POST /v1/preauthorizations/preauth_1/decisions",
		"GET /v1/claims",
		"POST /v1/claims",
		"GET /v1/claims/clm_1",
		"POST /v1/claims/clm_1/information_requests",
		"POST /v1/claims/clm_1/evidence",
		"POST /v1/claims/clm_1/adjudications",
		"GET /v1/remittances",
		"POST /v1/remittances",
		"GET /v1/remittances/rem_1",
		"GET /v1/webhook_endpoints",
		"POST /v1/webhook_endpoints",
		"GET /v1/webhook_endpoints/whe_1",
		"DELETE /v1/webhook_endpoints/whe_1",
		"POST /v1/webhook_endpoints/whe_1/test",
	}
	if !reflect.DeepEqual(received, expected) {
		t.Fatalf("operation surface mismatch\nreceived: %#v\nexpected: %#v", received, expected)
	}
}

func TestAuthenticationHeaders(t *testing.T) {
	for _, test := range []struct {
		name       string
		option     Option
		headerName string
		header     string
	}{
		{name: "bearer", option: WithAuthStyle(AuthBearer), headerName: "Authorization", header: "Bearer test_api_key"},
		{name: "api key", option: WithAuthStyle(AuthAPIKey), headerName: "X-Api-Key", header: "test_api_key"},
	} {
		t.Run(test.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
				if got := request.Header.Get(test.headerName); got != test.header {
					t.Errorf("authentication header = %q, want %q", got, test.header)
				}
				_, _ = response.Write([]byte(`{"object":"api","version":"v1","environment":"sandbox","resources":[]}`))
			}))
			defer server.Close()
			client, err := NewClient("test_api_key", WithBaseURL(server.URL), test.option)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := client.API.Retrieve(context.Background()); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestIdempotentWriteRetriesAndPreservesBody(t *testing.T) {
	var attempts atomic.Int32
	var bodies []string
	var mu sync.Mutex
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		body, _ := io.ReadAll(request.Body)
		mu.Lock()
		bodies = append(bodies, string(body))
		mu.Unlock()
		if request.Header.Get("Idempotency-Key") != "booking-0001" {
			t.Error("missing idempotency key")
		}
		if attempts.Add(1) < 3 {
			response.WriteHeader(http.StatusServiceUnavailable)
			_, _ = response.Write([]byte(`{"error":{"code":"service_unavailable","message":"Retry shortly.","docs":"https://docs.heyrafiki.space/errors"}}`))
			return
		}
		_, _ = response.Write([]byte(`{"id":"bkg_1"}`))
	}))
	defer server.Close()
	client, err := NewClient(
		"test_api_key",
		WithBaseURL(server.URL),
		WithRetryPolicy(RetryPolicy{MaxAttempts: 3}),
	)
	if err != nil {
		t.Fatal(err)
	}
	booking, err := client.Bookings.Create(context.Background(), BookingInput{PractitionerID: "prc_1"}, WriteOptions{IdempotencyKey: "booking-0001"})
	if err != nil {
		t.Fatal(err)
	}
	if booking.ID != "bkg_1" || attempts.Load() != 3 {
		t.Fatalf("booking = %#v, attempts = %d", booking, attempts.Load())
	}
	if len(bodies) != 3 || bodies[0] != bodies[1] || bodies[1] != bodies[2] {
		t.Fatalf("request body changed across retries: %#v", bodies)
	}
}

func TestNonIdempotentWriteIsNotRetried(t *testing.T) {
	var attempts atomic.Int32
	client := testClient(t, http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		attempts.Add(1)
		response.WriteHeader(http.StatusServiceUnavailable)
		_, _ = response.Write([]byte(`{"error":{"code":"service_unavailable","message":"Retry shortly.","docs":"https://docs.heyrafiki.space/errors"}}`))
	}))
	_, err := client.WebhookEndpoints.Create(context.Background(), WebhookEndpointInput{})
	if err == nil || attempts.Load() != 1 {
		t.Fatalf("error = %v, attempts = %d", err, attempts.Load())
	}
}

func TestTransportRetryUsesConfiguredHTTPClient(t *testing.T) {
	var attempts atomic.Int32
	httpClient := &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		if request.Header.Get("User-Agent") != "rafiki-go/0.1.0-beta.1 integration-tests/1.0" {
			t.Errorf("user agent = %q", request.Header.Get("User-Agent"))
		}
		if attempts.Add(1) == 1 {
			return nil, &net.DNSError{Err: "temporary failure", IsTimeout: true}
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`{"object":"api","version":"v1","environment":"sandbox","resources":[]}`)),
			Request:    request,
		}, nil
	})}
	client, err := NewClient(
		"test_api_key",
		WithHTTPClient(httpClient),
		WithUserAgent("integration-tests/1.0"),
		WithRetryPolicy(RetryPolicy{MaxAttempts: 2}),
	)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := client.API.Retrieve(context.Background()); err != nil {
		t.Fatal(err)
	}
	if attempts.Load() != 2 {
		t.Fatalf("attempts = %d, want 2", attempts.Load())
	}
}

func TestRetryHonorsContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	client := testClient(t, http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		cancel()
		response.WriteHeader(http.StatusServiceUnavailable)
		_, _ = response.Write([]byte(`{"error":{"code":"service_unavailable","message":"Retry shortly.","docs":"https://docs.heyrafiki.space/errors"}}`))
	}))
	client.retryPolicy = RetryPolicy{MaxAttempts: 2, BaseDelay: time.Minute, MaxDelay: time.Minute}
	_, err := client.API.Retrieve(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v, want context cancellation", err)
	}
}

func TestAPIErrorPreservesEnvelopeWithoutBodyDisclosure(t *testing.T) {
	client := testClient(t, http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("X-Request-Id", "req_123")
		response.WriteHeader(http.StatusForbidden)
		_, _ = response.Write([]byte(`{"error":{"code":"permission_denied","message":"This key does not permit the operation.","docs":"https://docs.heyrafiki.space/errors"},"secret":"must-not-surface"}`))
	}))
	_, err := client.Claims.Retrieve(context.Background(), "clm_1")
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error type = %T", err)
	}
	if apiErr.StatusCode != 403 || apiErr.Code != "permission_denied" || apiErr.RequestID != "req_123" {
		t.Fatalf("API error = %#v", apiErr)
	}
	if strings.Contains(err.Error(), "must-not-surface") {
		t.Fatal("error disclosed response body")
	}
}

func TestPreauthorizationDecisionAddsContractDiscriminator(t *testing.T) {
	client := testClient(t, http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body["outcome"] != "approved" {
			t.Errorf("outcome = %v", body["outcome"])
		}
		_, _ = response.Write([]byte(`{}`))
	}))
	_, err := client.Preauthorizations.Decide(
		context.Background(),
		"preauth_1",
		ApprovedPreauthorizationDecisionInput{ApprovedAmount: 100},
		WriteOptions{IdempotencyKey: "decision-0001"},
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestValidationRejectsUnsafeConfiguration(t *testing.T) {
	tests := []struct {
		name string
		make func() error
	}{
		{name: "missing key", make: func() error { _, err := NewClient(" "); return err }},
		{name: "insecure remote base URL", make: func() error { _, err := NewClient("key", WithBaseURL("http://api.example.com/v1")); return err }},
		{name: "unknown auth", make: func() error { _, err := NewClient("key", WithAuthStyle("unknown")); return err }},
		{name: "invalid list limit", make: func() error {
			client, err := NewClient("key")
			if err != nil {
				return err
			}
			_, err = client.Claims.List(context.Background(), &ListOptions{Limit: 101})
			return err
		}},
		{name: "short idempotency key", make: func() error {
			client, err := NewClient("key")
			if err != nil {
				return err
			}
			_, err = client.Claims.Create(context.Background(), ClaimInput{}, WriteOptions{IdempotencyKey: "short"})
			return err
		}},
		{name: "unsafe idempotency key", make: func() error {
			client, err := NewClient("key")
			if err != nil {
				return err
			}
			_, err = client.Claims.Create(context.Background(), ClaimInput{}, WriteOptions{IdempotencyKey: "request key"})
			return err
		}},
		{name: "invalid artifact reference", make: func() error {
			client, err := NewClient("key")
			if err != nil {
				return err
			}
			_, err = client.CoverageBatches.Record(context.Background(), CoverageBatchInput{}, CoverageBatchOptions{
				IdempotencyKey: "request-0001", ArtifactReference: "payer file",
			})
			return err
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := test.make(); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func TestResponseSizeLimit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		_, _ = fmt.Fprint(response, `{"value":"123456789"}`)
	}))
	defer server.Close()
	client, err := NewClient("key", WithBaseURL(server.URL), WithMaxResponseSize(8))
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.API.Retrieve(context.Background())
	var apiErr *APIError
	if !errors.As(err, &apiErr) || apiErr.Code != "response_too_large" {
		t.Fatalf("error = %v", err)
	}
}

func TestRetryAfterParsing(t *testing.T) {
	now := time.Date(2026, 8, 9, 0, 0, 0, 0, time.UTC)
	if got := parseRetryAfter("3", now); got != 3*time.Second {
		t.Fatalf("delta retry-after = %s", got)
	}
	if got := parseRetryAfter(now.Add(2*time.Second).Format(http.TimeFormat), now); got != 2*time.Second {
		t.Fatalf("date retry-after = %s", got)
	}
	if got := parseRetryAfter("invalid", now); got != 0 {
		t.Fatalf("invalid retry-after = %s", got)
	}
}

func TestInvalidJSONResponseWrapsDecodeError(t *testing.T) {
	client := testClient(t, http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("X-Request-Id", "req_invalid")
		_, _ = response.Write([]byte(`{"invalid"`))
	}))
	_, err := client.API.Retrieve(context.Background())
	var apiErr *APIError
	if !errors.As(err, &apiErr) || apiErr.Code != "invalid_response" || apiErr.Unwrap() == nil {
		t.Fatalf("error = %#v", err)
	}
}
