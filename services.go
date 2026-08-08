// SPDX-License-Identifier: Apache-2.0

package heyrafiki

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// APIService reads API metadata.
type APIService struct{ client *Client }

// PractitionerService discovers Practitioners and availability.
type PractitionerService struct{ client *Client }

// BookingService creates and reads Bookings.
type BookingService struct{ client *Client }

// SessionService reads Sessions.
type SessionService struct{ client *Client }

// EligibilityCheckService creates and reads Benefit eligibility checks.
type EligibilityCheckService struct{ client *Client }

// CoverageService records payer Coverage observations.
type CoverageService struct{ client *Client }

// CoverageBatchService records replay-safe payer Coverage batches.
type CoverageBatchService struct{ client *Client }

// PreauthorizationService creates, reads and decides pre-authorizations.
type PreauthorizationService struct{ client *Client }

// ClaimService creates and manages Claims.
type ClaimService struct{ client *Client }

// RemittanceService creates and reads remittance advice.
type RemittanceService struct{ client *Client }

// WebhookEndpointService manages Webhook endpoints and test deliveries.
type WebhookEndpointService struct{ client *Client }

// Retrieve returns API version and environment information.
func (service *APIService) Retrieve(ctx context.Context) (*APIInformation, error) {
	return doJSON[APIInformation](ctx, service.client, http.MethodGet, "/", nil, requestOptions{retryable: true})
}

// List returns visible Practitioners.
func (service *PractitionerService) List(ctx context.Context, options *ListOptions) (*PractitionerList, error) {
	query, err := validateListOptions(options)
	if err != nil {
		return nil, err
	}
	return doJSON[PractitionerList](ctx, service.client, http.MethodGet, "/practitioners"+query, nil, requestOptions{retryable: true})
}

// Retrieve returns one visible Practitioner.
func (service *PractitionerService) Retrieve(ctx context.Context, practitionerID string) (*Practitioner, error) {
	if err := validateID("Practitioner ID", practitionerID); err != nil {
		return nil, err
	}
	return doJSON[Practitioner](ctx, service.client, http.MethodGet, "/practitioners/"+url.PathEscape(practitionerID), nil, requestOptions{retryable: true})
}

// Availability returns a Practitioner's weekly availability.
func (service *PractitionerService) Availability(ctx context.Context, practitionerID string) (*PractitionerAvailability, error) {
	if err := validateID("Practitioner ID", practitionerID); err != nil {
		return nil, err
	}
	return doJSON[PractitionerAvailability](ctx, service.client, http.MethodGet, "/practitioners/"+url.PathEscape(practitionerID)+"/availability", nil, requestOptions{retryable: true})
}

// List returns Bookings visible to the project.
func (service *BookingService) List(ctx context.Context, options *ListOptions) (*BookingList, error) {
	query, err := validateListOptions(options)
	if err != nil {
		return nil, err
	}
	return doJSON[BookingList](ctx, service.client, http.MethodGet, "/bookings"+query, nil, requestOptions{retryable: true})
}

// Create creates a Booking with a caller-owned idempotency key.
func (service *BookingService) Create(ctx context.Context, input BookingInput, options WriteOptions) (*Booking, error) {
	if err := validateIdempotencyKey(options.IdempotencyKey); err != nil {
		return nil, err
	}
	return doJSON[Booking](ctx, service.client, http.MethodPost, "/bookings", input, requestOptions{idempotencyKey: options.IdempotencyKey, retryable: true})
}

// Retrieve returns one Booking visible to the project.
func (service *BookingService) Retrieve(ctx context.Context, bookingID string) (*Booking, error) {
	if err := validateID("Booking ID", bookingID); err != nil {
		return nil, err
	}
	return doJSON[Booking](ctx, service.client, http.MethodGet, "/bookings/"+url.PathEscape(bookingID), nil, requestOptions{retryable: true})
}

// List returns Sessions visible to the project.
func (service *SessionService) List(ctx context.Context, options *ListOptions) (*SessionList, error) {
	query, err := validateListOptions(options)
	if err != nil {
		return nil, err
	}
	return doJSON[SessionList](ctx, service.client, http.MethodGet, "/sessions"+query, nil, requestOptions{retryable: true})
}

// Retrieve returns one Session visible to the project.
func (service *SessionService) Retrieve(ctx context.Context, sessionID string) (*Session, error) {
	if err := validateID("Session ID", sessionID); err != nil {
		return nil, err
	}
	return doJSON[Session](ctx, service.client, http.MethodGet, "/sessions/"+url.PathEscape(sessionID), nil, requestOptions{retryable: true})
}

// Create checks Benefit eligibility with a caller-owned idempotency key.
func (service *EligibilityCheckService) Create(ctx context.Context, input EligibilityCheckInput, options WriteOptions) (*EligibilityCheck, error) {
	if err := validateIdempotencyKey(options.IdempotencyKey); err != nil {
		return nil, err
	}
	return doJSON[EligibilityCheck](ctx, service.client, http.MethodPost, "/eligibility_checks", input, requestOptions{idempotencyKey: options.IdempotencyKey, retryable: true})
}

// Retrieve returns one eligibility check visible to the project.
func (service *EligibilityCheckService) Retrieve(ctx context.Context, eligibilityCheckID string) (*EligibilityCheck, error) {
	if err := validateID("Eligibility Check ID", eligibilityCheckID); err != nil {
		return nil, err
	}
	return doJSON[EligibilityCheck](ctx, service.client, http.MethodGet, "/eligibility_checks/"+url.PathEscape(eligibilityCheckID), nil, requestOptions{retryable: true})
}

// Record stores a payer Coverage observation with a caller-owned idempotency key.
func (service *CoverageService) Record(ctx context.Context, input CoverageObservationInput, options WriteOptions) (*CoverageObservation, error) {
	if err := validateIdempotencyKey(options.IdempotencyKey); err != nil {
		return nil, err
	}
	return doJSON[CoverageObservation](ctx, service.client, http.MethodPost, "/coverages", input, requestOptions{idempotencyKey: options.IdempotencyKey, retryable: true})
}

// Record stores a payer Coverage batch with replay-safe artifact controls.
func (service *CoverageBatchService) Record(ctx context.Context, input CoverageBatchInput, options CoverageBatchOptions) (*CoverageBatchResult, error) {
	if err := validateIdempotencyKey(options.IdempotencyKey); err != nil {
		return nil, err
	}
	if err := validateArtifactReference(options.ArtifactReference); err != nil {
		return nil, err
	}
	return doJSON[CoverageBatchResult](ctx, service.client, http.MethodPost, "/coverage_batches", input, requestOptions{
		idempotencyKey: options.IdempotencyKey, artifactReference: options.ArtifactReference, retryable: true,
	})
}

// Create requests pre-authorization with a caller-owned idempotency key.
func (service *PreauthorizationService) Create(ctx context.Context, input PreauthorizationInput, options WriteOptions) (*Preauthorization, error) {
	if err := validateIdempotencyKey(options.IdempotencyKey); err != nil {
		return nil, err
	}
	return doJSON[Preauthorization](ctx, service.client, http.MethodPost, "/preauthorizations", input, requestOptions{idempotencyKey: options.IdempotencyKey, retryable: true})
}

// Retrieve returns one pre-authorization visible to the project.
func (service *PreauthorizationService) Retrieve(ctx context.Context, preauthorizationID string) (*Preauthorization, error) {
	if err := validateID("Preauthorization ID", preauthorizationID); err != nil {
		return nil, err
	}
	return doJSON[Preauthorization](ctx, service.client, http.MethodGet, "/preauthorizations/"+url.PathEscape(preauthorizationID), nil, requestOptions{retryable: true})
}

// Decide records an approved or denied pre-authorization decision.
func (service *PreauthorizationService) Decide(ctx context.Context, preauthorizationID string, input PreauthorizationDecisionInput, options WriteOptions) (*Preauthorization, error) {
	if err := validateID("Preauthorization ID", preauthorizationID); err != nil {
		return nil, err
	}
	if err := validateIdempotencyKey(options.IdempotencyKey); err != nil {
		return nil, err
	}
	payload, err := preauthorizationDecisionPayload(input)
	if err != nil {
		return nil, err
	}
	return doJSON[Preauthorization](ctx, service.client, http.MethodPost, "/preauthorizations/"+url.PathEscape(preauthorizationID)+"/decisions", payload, requestOptions{idempotencyKey: options.IdempotencyKey, retryable: true})
}

func preauthorizationDecisionPayload(input PreauthorizationDecisionInput) (any, error) {
	switch value := input.(type) {
	case ApprovedPreauthorizationDecisionInput:
		return struct {
			Outcome string `json:"outcome"`
			ApprovedPreauthorizationDecisionInput
		}{Outcome: "approved", ApprovedPreauthorizationDecisionInput: value}, nil
	case *ApprovedPreauthorizationDecisionInput:
		if value == nil {
			return nil, errors.New("heyrafiki: preauthorization decision is required")
		}
		return preauthorizationDecisionPayload(*value)
	case DeniedPreauthorizationDecisionInput:
		return struct {
			Outcome string `json:"outcome"`
			DeniedPreauthorizationDecisionInput
		}{Outcome: "denied", DeniedPreauthorizationDecisionInput: value}, nil
	case *DeniedPreauthorizationDecisionInput:
		if value == nil {
			return nil, errors.New("heyrafiki: preauthorization decision is required")
		}
		return preauthorizationDecisionPayload(*value)
	default:
		return nil, errors.New("heyrafiki: preauthorization decision is required")
	}
}

// List returns Claims visible to the project.
func (service *ClaimService) List(ctx context.Context, options *ListOptions) (*ClaimList, error) {
	query, err := validateListOptions(options)
	if err != nil {
		return nil, err
	}
	return doJSON[ClaimList](ctx, service.client, http.MethodGet, "/claims"+query, nil, requestOptions{retryable: true})
}

// Create submits a Claim with a caller-owned idempotency key.
func (service *ClaimService) Create(ctx context.Context, input ClaimInput, options WriteOptions) (*Claim, error) {
	if err := validateIdempotencyKey(options.IdempotencyKey); err != nil {
		return nil, err
	}
	return doJSON[Claim](ctx, service.client, http.MethodPost, "/claims", input, requestOptions{idempotencyKey: options.IdempotencyKey, retryable: true})
}

// Retrieve returns one Claim visible to the project.
func (service *ClaimService) Retrieve(ctx context.Context, claimID string) (*Claim, error) {
	if err := validateID("Claim ID", claimID); err != nil {
		return nil, err
	}
	return doJSON[Claim](ctx, service.client, http.MethodGet, "/claims/"+url.PathEscape(claimID), nil, requestOptions{retryable: true})
}

// RequestInformation records a payer request for Claim evidence.
func (service *ClaimService) RequestInformation(ctx context.Context, claimID string, input ClaimInformationRequestInput, options WriteOptions) (*Claim, error) {
	if err := validateID("Claim ID", claimID); err != nil {
		return nil, err
	}
	if err := validateIdempotencyKey(options.IdempotencyKey); err != nil {
		return nil, err
	}
	return doJSON[Claim](ctx, service.client, http.MethodPost, "/claims/"+url.PathEscape(claimID)+"/information_requests", input, requestOptions{idempotencyKey: options.IdempotencyKey, retryable: true})
}

// SubmitEvidence attaches evidence references to a Claim information request.
func (service *ClaimService) SubmitEvidence(ctx context.Context, claimID string, input ClaimEvidenceInput, options WriteOptions) (*Claim, error) {
	if err := validateID("Claim ID", claimID); err != nil {
		return nil, err
	}
	if err := validateIdempotencyKey(options.IdempotencyKey); err != nil {
		return nil, err
	}
	return doJSON[Claim](ctx, service.client, http.MethodPost, "/claims/"+url.PathEscape(claimID)+"/evidence", input, requestOptions{idempotencyKey: options.IdempotencyKey, retryable: true})
}

// Adjudicate records a Claim adjudication decision.
func (service *ClaimService) Adjudicate(ctx context.Context, claimID string, input ClaimAdjudicationInput, options WriteOptions) (*Claim, error) {
	if err := validateID("Claim ID", claimID); err != nil {
		return nil, err
	}
	if err := validateIdempotencyKey(options.IdempotencyKey); err != nil {
		return nil, err
	}
	return doJSON[Claim](ctx, service.client, http.MethodPost, "/claims/"+url.PathEscape(claimID)+"/adjudications", input, requestOptions{idempotencyKey: options.IdempotencyKey, retryable: true})
}

// List returns remittances visible to the project.
func (service *RemittanceService) List(ctx context.Context, options *ListOptions) (*RemittanceList, error) {
	query, err := validateListOptions(options)
	if err != nil {
		return nil, err
	}
	return doJSON[RemittanceList](ctx, service.client, http.MethodGet, "/remittances"+query, nil, requestOptions{retryable: true})
}

// Create records remittance advice with a caller-owned idempotency key.
func (service *RemittanceService) Create(ctx context.Context, input RemittanceInput, options WriteOptions) (*Remittance, error) {
	if err := validateIdempotencyKey(options.IdempotencyKey); err != nil {
		return nil, err
	}
	return doJSON[Remittance](ctx, service.client, http.MethodPost, "/remittances", input, requestOptions{idempotencyKey: options.IdempotencyKey, retryable: true})
}

// Retrieve returns one remittance visible to the project.
func (service *RemittanceService) Retrieve(ctx context.Context, remittanceID string) (*Remittance, error) {
	if err := validateID("Remittance ID", remittanceID); err != nil {
		return nil, err
	}
	return doJSON[Remittance](ctx, service.client, http.MethodGet, "/remittances/"+url.PathEscape(remittanceID), nil, requestOptions{retryable: true})
}

// List returns Webhook endpoints visible to the project.
func (service *WebhookEndpointService) List(ctx context.Context, options *ListOptions) (*WebhookEndpointList, error) {
	query, err := validateListOptions(options)
	if err != nil {
		return nil, err
	}
	return doJSON[WebhookEndpointList](ctx, service.client, http.MethodGet, "/webhook_endpoints"+query, nil, requestOptions{retryable: true})
}

// Create registers a Webhook endpoint and returns its signing secret once.
func (service *WebhookEndpointService) Create(ctx context.Context, input WebhookEndpointInput) (*WebhookEndpointWithSecret, error) {
	return doJSON[WebhookEndpointWithSecret](ctx, service.client, http.MethodPost, "/webhook_endpoints", input, requestOptions{})
}

// Retrieve returns one Webhook endpoint visible to the project.
func (service *WebhookEndpointService) Retrieve(ctx context.Context, endpointID string) (*WebhookEndpoint, error) {
	if err := validateID("Webhook Endpoint ID", endpointID); err != nil {
		return nil, err
	}
	return doJSON[WebhookEndpoint](ctx, service.client, http.MethodGet, "/webhook_endpoints/"+url.PathEscape(endpointID), nil, requestOptions{retryable: true})
}

// Disable disables a Webhook endpoint.
func (service *WebhookEndpointService) Disable(ctx context.Context, endpointID string) (*WebhookEndpoint, error) {
	if err := validateID("Webhook Endpoint ID", endpointID); err != nil {
		return nil, err
	}
	return doJSON[WebhookEndpoint](ctx, service.client, http.MethodDelete, "/webhook_endpoints/"+url.PathEscape(endpointID), nil, requestOptions{retryable: true})
}

// SendTest sends a contract-defined test event to a Webhook endpoint.
func (service *WebhookEndpointService) SendTest(ctx context.Context, endpointID string) (*WebhookDelivery, error) {
	if err := validateID("Webhook Endpoint ID", endpointID); err != nil {
		return nil, err
	}
	return doJSON[WebhookDelivery](ctx, service.client, http.MethodPost, "/webhook_endpoints/"+url.PathEscape(endpointID)+"/test", nil, requestOptions{})
}
