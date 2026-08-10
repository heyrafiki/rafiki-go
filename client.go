// SPDX-License-Identifier: Apache-2.0

package heyrafiki

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var artifactReferencePattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9:._/-]*$`)

const (
	defaultBaseURL      = "https://api.heyrafiki.space/v1"
	defaultUserAgent    = "rafiki-go/0.1.0-beta.1"
	defaultResponseSize = int64(4 << 20)
)

// AuthStyle selects one of the authentication schemes published by the API contract.
type AuthStyle string

const (
	// AuthBearer sends the credential in the Authorization header.
	AuthBearer AuthStyle = "bearer"
	// AuthAPIKey sends the credential in the x-api-key header.
	AuthAPIKey AuthStyle = "api_key"
)

// Client is a concurrency-safe client for the Heyrafiki API.
type Client struct {
	apiKey          string
	authStyle       AuthStyle
	baseURL         string
	httpClient      *http.Client
	retryPolicy     RetryPolicy
	userAgent       string
	maxResponseSize int64

	API               *APIService
	Practitioners     *PractitionerService
	Bookings          *BookingService
	Sessions          *SessionService
	EligibilityChecks *EligibilityCheckService
	Coverages         *CoverageService
	CoverageBatches   *CoverageBatchService
	Preauthorizations *PreauthorizationService
	Claims            *ClaimService
	Remittances       *RemittanceService
	WebhookEndpoints  *WebhookEndpointService
}

// Option configures a Client.
type Option func(*clientConfig) error

type clientConfig struct {
	baseURL         string
	httpClient      *http.Client
	authStyle       AuthStyle
	retryPolicy     RetryPolicy
	userAgent       string
	maxResponseSize int64
}

// NewClient creates a server-side Heyrafiki API client.
func NewClient(apiKey string, options ...Option) (*Client, error) {
	if strings.TrimSpace(apiKey) == "" {
		return nil, errors.New("heyrafiki: API key is required")
	}
	if strings.ContainsAny(apiKey, "\r\n") {
		return nil, errors.New("heyrafiki: API key contains invalid characters")
	}

	config := clientConfig{
		baseURL:         defaultBaseURL,
		httpClient:      &http.Client{Timeout: 30 * time.Second},
		authStyle:       AuthBearer,
		retryPolicy:     DefaultRetryPolicy(),
		userAgent:       defaultUserAgent,
		maxResponseSize: defaultResponseSize,
	}
	for _, option := range options {
		if option == nil {
			return nil, errors.New("heyrafiki: client option cannot be nil")
		}
		if err := option(&config); err != nil {
			return nil, err
		}
	}

	client := &Client{
		apiKey:          apiKey,
		authStyle:       config.authStyle,
		baseURL:         config.baseURL,
		httpClient:      config.httpClient,
		retryPolicy:     config.retryPolicy,
		userAgent:       config.userAgent,
		maxResponseSize: config.maxResponseSize,
	}
	client.API = &APIService{client: client}
	client.Practitioners = &PractitionerService{client: client}
	client.Bookings = &BookingService{client: client}
	client.Sessions = &SessionService{client: client}
	client.EligibilityChecks = &EligibilityCheckService{client: client}
	client.Coverages = &CoverageService{client: client}
	client.CoverageBatches = &CoverageBatchService{client: client}
	client.Preauthorizations = &PreauthorizationService{client: client}
	client.Claims = &ClaimService{client: client}
	client.Remittances = &RemittanceService{client: client}
	client.WebhookEndpoints = &WebhookEndpointService{client: client}
	return client, nil
}

// WithBaseURL replaces the API base URL. HTTPS is required except for loopback test servers.
func WithBaseURL(rawURL string) Option {
	return func(config *clientConfig) error {
		parsed, err := url.Parse(rawURL)
		if err != nil || !parsed.IsAbs() || parsed.Host == "" {
			return errors.New("heyrafiki: base URL must be absolute")
		}
		if parsed.User != nil || parsed.RawQuery != "" || parsed.Fragment != "" {
			return errors.New("heyrafiki: base URL cannot contain credentials, a query, or a fragment")
		}
		host := parsed.Hostname()
		if parsed.Scheme != "https" && !(parsed.Scheme == "http" && isLoopbackHost(host)) {
			return errors.New("heyrafiki: base URL must use HTTPS")
		}
		config.baseURL = strings.TrimRight(parsed.String(), "/")
		return nil
	}
}

func isLoopbackHost(host string) bool {
	if strings.EqualFold(host, "localhost") {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

// WithHTTPClient supplies the HTTP client used for all requests.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(config *clientConfig) error {
		if httpClient == nil {
			return errors.New("heyrafiki: HTTP client cannot be nil")
		}
		config.httpClient = httpClient
		return nil
	}
}

// WithAuthStyle selects bearer or x-api-key authentication.
func WithAuthStyle(style AuthStyle) Option {
	return func(config *clientConfig) error {
		if style != AuthBearer && style != AuthAPIKey {
			return fmt.Errorf("heyrafiki: unsupported authentication style %q", style)
		}
		config.authStyle = style
		return nil
	}
}

// WithRetryPolicy replaces the default retry policy.
func WithRetryPolicy(policy RetryPolicy) Option {
	return func(config *clientConfig) error {
		if err := policy.validate(); err != nil {
			return err
		}
		config.retryPolicy = policy
		return nil
	}
}

// WithUserAgent appends a caller identifier to the SDK user agent.
func WithUserAgent(identifier string) Option {
	return func(config *clientConfig) error {
		identifier = strings.TrimSpace(identifier)
		if identifier == "" || strings.ContainsAny(identifier, "\r\n") {
			return errors.New("heyrafiki: user agent identifier is invalid")
		}
		config.userAgent = defaultUserAgent + " " + identifier
		return nil
	}
}

// WithMaxResponseSize limits response bodies read into memory.
func WithMaxResponseSize(bytes int64) Option {
	return func(config *clientConfig) error {
		if bytes < 1 {
			return errors.New("heyrafiki: maximum response size must be positive")
		}
		config.maxResponseSize = bytes
		return nil
	}
}

type requestOptions struct {
	idempotencyKey    string
	artifactReference string
	retryable         bool
}

func doJSON[T any](ctx context.Context, client *Client, method, path string, body any, options requestOptions) (*T, error) {
	if ctx == nil {
		return nil, errors.New("heyrafiki: context cannot be nil")
	}
	var encoded []byte
	var err error
	if body != nil {
		encoded, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("heyrafiki: encode request: %w", err)
		}
	}

	attempts := 1
	if options.retryable {
		attempts = client.retryPolicy.MaxAttempts
	}
	for attempt := 1; attempt <= attempts; attempt++ {
		request, requestErr := http.NewRequestWithContext(ctx, method, client.baseURL+path, bytes.NewReader(encoded))
		if requestErr != nil {
			return nil, fmt.Errorf("heyrafiki: create request: %w", requestErr)
		}
		request.Header.Set("Accept", "application/json")
		request.Header.Set("User-Agent", client.userAgent)
		if body != nil {
			request.Header.Set("Content-Type", "application/json")
		}
		if options.idempotencyKey != "" {
			request.Header.Set("Idempotency-Key", options.idempotencyKey)
		}
		if options.artifactReference != "" {
			request.Header.Set("X-Heyrafiki-Artifact-Reference", options.artifactReference)
		}
		client.authenticate(request)

		response, requestErr := client.httpClient.Do(request)
		if requestErr != nil {
			if attempt < attempts && retryableTransportError(requestErr) {
				if err := client.waitForRetry(ctx, attempt, 0); err != nil {
					return nil, err
				}
				continue
			}
			return nil, fmt.Errorf("heyrafiki: request failed: %w", requestErr)
		}

		payload, readErr := readResponse(response, client.maxResponseSize)
		if readErr != nil {
			return nil, readErr
		}
		if response.StatusCode < 200 || response.StatusCode >= 300 {
			apiErr := parseAPIError(response, payload)
			if attempt < attempts && retryableStatus(response.StatusCode) {
				if err := client.waitForRetry(ctx, attempt, apiErr.RetryAfter); err != nil {
					return nil, err
				}
				continue
			}
			return nil, apiErr
		}
		if len(payload) == 0 {
			return nil, &APIError{
				StatusCode: response.StatusCode,
				Code:       "invalid_response",
				Message:    "The Heyrafiki API returned an empty response.",
				RequestID:  response.Header.Get("X-Request-Id"),
			}
		}
		var result T
		if err := json.Unmarshal(payload, &result); err != nil {
			return nil, &APIError{
				StatusCode: response.StatusCode,
				Code:       "invalid_response",
				Message:    "The Heyrafiki API returned an invalid JSON response.",
				RequestID:  response.Header.Get("X-Request-Id"),
				cause:      err,
			}
		}
		return &result, nil
	}
	return nil, errors.New("heyrafiki: retry loop exhausted")
}

func (client *Client) authenticate(request *http.Request) {
	if client.authStyle == AuthAPIKey {
		request.Header.Set("X-Api-Key", client.apiKey)
		return
	}
	request.Header.Set("Authorization", "Bearer "+client.apiKey)
}

func readResponse(response *http.Response, limit int64) ([]byte, error) {
	defer response.Body.Close()
	payload, err := io.ReadAll(io.LimitReader(response.Body, limit+1))
	if err != nil {
		return nil, fmt.Errorf("heyrafiki: read response: %w", err)
	}
	if int64(len(payload)) > limit {
		return nil, &APIError{
			StatusCode: response.StatusCode,
			Code:       "response_too_large",
			Message:    "The Heyrafiki API response exceeded the configured limit.",
			RequestID:  response.Header.Get("X-Request-Id"),
		}
	}
	return payload, nil
}

func validateListOptions(options *ListOptions) (string, error) {
	if options == nil || options.Limit == 0 {
		return "", nil
	}
	if options.Limit < 1 || options.Limit > 100 {
		return "", errors.New("heyrafiki: list limit must be between 1 and 100")
	}
	return "?limit=" + strconv.Itoa(options.Limit), nil
}

func validateID(label, value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("heyrafiki: %s is required", label)
	}
	return nil
}

func validateIdempotencyKey(value string) error {
	if len(value) < 8 || len(value) > 255 || !visibleASCII(value) {
		return errors.New("heyrafiki: idempotency key must contain 8 to 255 visible ASCII characters")
	}
	return nil
}

func validateArtifactReference(value string) error {
	if len(value) > 190 || !artifactReferencePattern.MatchString(value) {
		return errors.New("heyrafiki: artifact reference does not match the published contract")
	}
	return nil
}

func visibleASCII(value string) bool {
	for index := 0; index < len(value); index++ {
		if value[index] < 0x21 || value[index] > 0x7e {
			return false
		}
	}
	return true
}
