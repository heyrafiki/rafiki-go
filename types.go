// SPDX-License-Identifier: Apache-2.0

package heyrafiki

// ListOptions controls paginated list operations.
type ListOptions struct {
	Limit int
}

// WriteOptions supplies the caller-owned idempotency key required by a write operation.
type WriteOptions struct {
	IdempotencyKey string
}

// CoverageBatchOptions supplies the replay controls required by Coverage batch ingestion.
type CoverageBatchOptions struct {
	IdempotencyKey    string
	ArtifactReference string
}

// Currency is an ISO 4217 currency accepted by the version 1 contract.
type Currency string

const (
	CurrencyKES Currency = "KES"
	CurrencyUSD Currency = "USD"
	CurrencyEUR Currency = "EUR"
	CurrencyGBP Currency = "GBP"
)

// CareFormat identifies how a Session takes place.
type CareFormat string

const (
	CareFormatOnline   CareFormat = "online"
	CareFormatInPerson CareFormat = "in_person"
	CareFormatPhone    CareFormat = "phone"
)

// PaymentSource identifies how a Session is funded.
type PaymentSource string

const (
	PaymentSourceSelfPay PaymentSource = "self_pay"
	PaymentSourceCovered PaymentSource = "covered"
)

type APIInformation struct {
	Object      string   `json:"object"`
	Version     string   `json:"version"`
	Environment string   `json:"environment"`
	Resources   []string `json:"resources"`
}

type PractitionerList struct {
	Object  string         `json:"object"`
	Data    []Practitioner `json:"data"`
	HasMore bool           `json:"has_more"`
}

type Practitioner struct {
	ID         string               `json:"id"`
	Object     string               `json:"object"`
	Name       string               `json:"name"`
	Profession string               `json:"profession"`
	Location   PractitionerLocation `json:"location"`
	SessionFee Money                `json:"session_fee"`
}

type PractitionerLocation struct {
	City    string `json:"city"`
	Country string `json:"country"`
}

// Money contains an amount in the currency's minor unit.
type Money struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}

type PractitionerAvailability struct {
	Object         string               `json:"object"`
	PractitionerID string               `json:"practitioner_id"`
	Timezone       string               `json:"timezone"`
	WeeklyHours    []AvailabilityWindow `json:"weekly_hours"`
}

type AvailabilityWindow struct {
	Weekday int          `json:"weekday"`
	Start   string       `json:"start"`
	End     string       `json:"end"`
	Formats []CareFormat `json:"formats"`
}

type BookingList struct {
	Object  string    `json:"object"`
	Data    []Booking `json:"data"`
	HasMore bool      `json:"has_more"`
}

type Booking struct {
	ID             string        `json:"id"`
	Object         string        `json:"object"`
	SessionID      string        `json:"session_id"`
	PractitionerID string        `json:"practitioner_id"`
	StartsAt       string        `json:"starts_at"`
	EndsAt         string        `json:"ends_at"`
	Timezone       string        `json:"timezone"`
	Format         CareFormat    `json:"format"`
	Status         string        `json:"status"`
	PaymentSource  PaymentSource `json:"payment_source"`
}

type BookingInput struct {
	PractitionerID string        `json:"practitioner_id"`
	StartsAt       string        `json:"starts_at"`
	EndsAt         string        `json:"ends_at"`
	Format         CareFormat    `json:"format"`
	PaymentSource  PaymentSource `json:"payment_source"`
}

type SessionList struct {
	Object  string    `json:"object"`
	Data    []Session `json:"data"`
	HasMore bool      `json:"has_more"`
}

type Session struct {
	ID             string        `json:"id"`
	Object         string        `json:"object"`
	PractitionerID string        `json:"practitioner_id"`
	StartsAt       string        `json:"starts_at"`
	EndsAt         string        `json:"ends_at"`
	Timezone       string        `json:"timezone"`
	Format         CareFormat    `json:"format"`
	Status         string        `json:"status"`
	PaymentSource  PaymentSource `json:"payment_source"`
}

type EligibilityCheckInput struct {
	MemberReference string   `json:"member_reference"`
	ServiceCode     string   `json:"service_code"`
	ScheduledAt     string   `json:"scheduled_at"`
	Amount          int64    `json:"amount"`
	Currency        Currency `json:"currency"`
}

type EligibilityCheck struct {
	ID                    string             `json:"id"`
	Object                string             `json:"object"`
	Status                string             `json:"status"`
	ReasonCodes           []string           `json:"reason_codes"`
	Service               EligibilityService `json:"service"`
	Amount                EligibilityAmount  `json:"amount"`
	AuthorizationRequired bool               `json:"authorization_required"`
	RemainingSessions     *int               `json:"remaining_sessions"`
	CoverageValidUntil    *string            `json:"coverage_valid_until"`
	CheckedAt             string             `json:"checked_at"`
}

type EligibilityService struct {
	Code        string `json:"code"`
	ScheduledAt string `json:"scheduled_at"`
}

type EligibilityAmount struct {
	Requested         int64  `json:"requested"`
	Currency          string `json:"currency"`
	MaximumPerSession *int64 `json:"maximum_per_session"`
}

type CoverageStatus string

const (
	CoverageStatusActive    CoverageStatus = "active"
	CoverageStatusPaused    CoverageStatus = "paused"
	CoverageStatusExhausted CoverageStatus = "exhausted"
	CoverageStatusExpired   CoverageStatus = "expired"
)

type CoverageObservationInput struct {
	SourceContractReference   string         `json:"source_contract_reference"`
	ExternalCoverageReference string         `json:"external_coverage_reference"`
	SourceVersion             string         `json:"source_version"`
	TenantReference           string         `json:"tenant_reference"`
	MemberReference           string         `json:"member_reference"`
	PlanName                  string         `json:"plan_name"`
	ServiceCode               string         `json:"service_code"`
	Status                    CoverageStatus `json:"status"`
	Currency                  Currency       `json:"currency"`
	AmountLimit               int64          `json:"amount_limit"`
	RemainingSessions         int            `json:"remaining_sessions"`
	AuthorizationRequired     bool           `json:"authorization_required"`
	CoordinationPriority      *int           `json:"coordination_priority"`
	ValidFrom                 string         `json:"valid_from"`
	ValidUntil                string         `json:"valid_until"`
	ObservedAt                string         `json:"observed_at"`
	EvidenceReferences        []string       `json:"evidence_references"`
}

type CoverageObservation struct {
	ID                        string              `json:"id"`
	Object                    string              `json:"object"`
	CoverageID                string              `json:"coverage_id"`
	Source                    string              `json:"source"`
	SourceContractReference   string              `json:"source_contract_reference"`
	ExternalCoverageReference string              `json:"external_coverage_reference"`
	SourceVersion             string              `json:"source_version"`
	SnapshotVersion           int                 `json:"snapshot_version"`
	Status                    CoverageStatus      `json:"status"`
	ServiceCode               string              `json:"service_code"`
	AmountLimit               CoverageAmountLimit `json:"amount_limit"`
	RemainingSessions         int                 `json:"remaining_sessions"`
	AuthorizationRequired     bool                `json:"authorization_required"`
	CoordinationPriority      *int                `json:"coordination_priority"`
	ValidFrom                 string              `json:"valid_from"`
	ValidUntil                string              `json:"valid_until"`
	ObservedAt                string              `json:"observed_at"`
}

type CoverageAmountLimit struct {
	Currency Currency `json:"currency"`
	Value    int64    `json:"value"`
}

type CoverageBatchRecordInput struct {
	CoverageReference     string         `json:"coverage_reference"`
	RecordVersion         string         `json:"record_version"`
	TenantReference       string         `json:"tenant_reference"`
	MemberReference       string         `json:"member_reference"`
	PlanName              string         `json:"plan_name"`
	ServiceCode           string         `json:"service_code"`
	Status                CoverageStatus `json:"status"`
	Currency              Currency       `json:"currency"`
	AmountLimit           int64          `json:"amount_limit"`
	RemainingSessions     int            `json:"remaining_sessions"`
	AuthorizationRequired bool           `json:"authorization_required"`
	CoordinationPriority  *int           `json:"coordination_priority"`
	ValidFrom             string         `json:"valid_from"`
	ValidUntil            string         `json:"valid_until"`
}

type CoverageBatchInput struct {
	ContractVersion         string                     `json:"contract_version"`
	BatchReference          string                     `json:"batch_reference"`
	BatchVersion            string                     `json:"batch_version"`
	SourceContractReference string                     `json:"source_contract_reference"`
	GeneratedAt             string                     `json:"generated_at"`
	Records                 []CoverageBatchRecordInput `json:"records"`
}

type CoverageBatchResult struct {
	Object         string                `json:"object"`
	BatchReference string                `json:"batch_reference"`
	ArtifactSHA256 string                `json:"artifact_sha256"`
	Total          int                   `json:"total"`
	Recorded       int                   `json:"recorded"`
	Replayed       int                   `json:"replayed"`
	Observations   []CoverageObservation `json:"observations"`
}

type PreauthorizationInput struct {
	EligibilityCheckID string `json:"eligibility_check_id"`
	BookingID          string `json:"booking_id"`
}

// PreauthorizationDecisionInput is implemented by the approved and denied decision variants.
type PreauthorizationDecisionInput interface {
	preauthorizationDecisionInput()
}

type ApprovedPreauthorizationDecisionInput struct {
	ApprovedAmount     int64    `json:"approved_amount"`
	ValidUntil         string   `json:"valid_until"`
	ReasonCodes        []string `json:"reason_codes"`
	PolicyReference    string   `json:"policy_reference"`
	PolicyVersion      string   `json:"policy_version"`
	EvidenceReferences []string `json:"evidence_references"`
}

func (ApprovedPreauthorizationDecisionInput) preauthorizationDecisionInput() {}

type DeniedPreauthorizationDecisionInput struct {
	ReasonCodes        []string `json:"reason_codes"`
	PolicyReference    string   `json:"policy_reference"`
	PolicyVersion      string   `json:"policy_version"`
	EvidenceReferences []string `json:"evidence_references"`
}

func (DeniedPreauthorizationDecisionInput) preauthorizationDecisionInput() {}

type Preauthorization struct {
	ID                 string                    `json:"id"`
	Object             string                    `json:"object"`
	EligibilityCheckID string                    `json:"eligibility_check_id"`
	BookingID          string                    `json:"booking_id"`
	Status             string                    `json:"status"`
	ReasonCodes        []string                  `json:"reason_codes"`
	Amount             PreauthorizationAmount    `json:"amount"`
	ValidUntil         *string                   `json:"valid_until"`
	CreatedAt          string                    `json:"created_at"`
	Decision           *PreauthorizationDecision `json:"decision"`
}

type PreauthorizationAmount struct {
	Requested int64  `json:"requested"`
	Approved  *int64 `json:"approved"`
	Currency  string `json:"currency"`
}

type PreauthorizationDecision struct {
	ID                 string          `json:"id"`
	Object             string          `json:"object"`
	Version            int             `json:"version"`
	Outcome            string          `json:"outcome"`
	ReasonCodes        []string        `json:"reason_codes"`
	Policy             PolicyReference `json:"policy"`
	AuthorityReference string          `json:"authority_reference"`
	EvidenceReferences []string        `json:"evidence_references"`
	DecidedAt          string          `json:"decided_at"`
}

type PolicyReference struct {
	Reference string `json:"reference"`
	Version   string `json:"version"`
}

type ClaimList struct {
	Object  string  `json:"object"`
	Data    []Claim `json:"data"`
	HasMore bool    `json:"has_more"`
}

type Claim struct {
	ID                     string                    `json:"id"`
	Object                 string                    `json:"object"`
	Status                 string                    `json:"status"`
	ProviderClaimReference *string                   `json:"provider_claim_reference"`
	SubmissionVersion      int                       `json:"submission_version"`
	Amount                 ClaimAmount               `json:"amount"`
	ServicePeriod          ServicePeriod             `json:"service_period"`
	Lines                  []ClaimLine               `json:"lines"`
	InformationRequests    []ClaimInformationRequest `json:"information_requests"`
	Adjudication           *ClaimAdjudication        `json:"adjudication"`
	SubmittedAt            *string                   `json:"submitted_at"`
	UpdatedAt              string                    `json:"updated_at"`
}

type ClaimAmount struct {
	Currency string `json:"currency"`
	Billed   int64  `json:"billed"`
	Approved *int64 `json:"approved"`
	Remitted *int64 `json:"remitted"`
	Settled  *int64 `json:"settled"`
}

type ServicePeriod struct {
	StartsAt string `json:"starts_at"`
	EndsAt   string `json:"ends_at"`
}

type ClaimLine struct {
	LineNumber        int     `json:"line_number"`
	CodeSystem        string  `json:"code_system"`
	CodeSystemVersion *string `json:"code_system_version"`
	ServiceCode       string  `json:"service_code"`
	Units             float64 `json:"units"`
	Amount            int64   `json:"amount"`
}

type ClaimInput struct {
	EligibilityCheckID     string           `json:"eligibility_check_id"`
	PreauthorizationID     *string          `json:"preauthorization_id,omitempty"`
	SessionID              string           `json:"session_id"`
	ProviderClaimReference string           `json:"provider_claim_reference"`
	EvidenceRefs           []string         `json:"evidence_refs,omitempty"`
	Lines                  []ClaimLineInput `json:"lines"`
}

type ClaimLineInput struct {
	CodeSystem        string  `json:"code_system"`
	CodeSystemVersion *string `json:"code_system_version,omitempty"`
	ServiceCode       string  `json:"service_code"`
	Units             float64 `json:"units"`
	Amount            int64   `json:"amount"`
}

type ClaimInformationRequest struct {
	ID                     string   `json:"id"`
	Object                 string   `json:"object"`
	ReasonCode             string   `json:"reason_code"`
	RequestedEvidenceTypes []string `json:"requested_evidence_types"`
	Status                 string   `json:"status"`
	DueAt                  *string  `json:"due_at"`
	CreatedAt              string   `json:"created_at"`
	ResolvedAt             *string  `json:"resolved_at"`
}

type ClaimInformationRequestInput struct {
	ReasonCode             string   `json:"reason_code"`
	RequestedEvidenceTypes []string `json:"requested_evidence_types"`
	DueAt                  *string  `json:"due_at,omitempty"`
}

type ClaimEvidenceInput struct {
	InformationRequestID string   `json:"information_request_id"`
	EvidenceRefs         []string `json:"evidence_refs"`
}

type ClaimAdjudicationLine struct {
	LineNumber  int                      `json:"line_number"`
	Amount      ClaimAdjudicationAmounts `json:"amount"`
	ReasonCodes []string                 `json:"reason_codes"`
}

type ClaimAdjudicationAmounts struct {
	Billed                int64 `json:"billed"`
	Allowed               int64 `json:"allowed"`
	Payer                 int64 `json:"payer"`
	PatientResponsibility int64 `json:"patient_responsibility"`
	Adjustment            int64 `json:"adjustment"`
}

type ClaimAdjudicationInput struct {
	PolicyReference string                  `json:"policy_reference"`
	PolicyVersion   string                  `json:"policy_version"`
	ReasonCodes     []string                `json:"reason_codes"`
	EvidenceRefs    []string                `json:"evidence_refs,omitempty"`
	Lines           []ClaimAdjudicationLine `json:"lines"`
}

type ClaimAdjudication struct {
	ID          string                  `json:"id"`
	Object      string                  `json:"object"`
	Version     int                     `json:"version"`
	Decision    string                  `json:"decision"`
	Amount      ClaimAdjudicationTotal  `json:"amount"`
	ReasonCodes []string                `json:"reason_codes"`
	Policy      PolicyReference         `json:"policy"`
	Lines       []ClaimAdjudicationLine `json:"lines"`
	DecidedAt   string                  `json:"decided_at"`
}

type ClaimAdjudicationTotal struct {
	Currency              string `json:"currency"`
	Billed                int64  `json:"billed"`
	Payer                 int64  `json:"payer"`
	PatientResponsibility int64  `json:"patient_responsibility"`
	Adjustment            int64  `json:"adjustment"`
}

type RemittanceList struct {
	Object  string       `json:"object"`
	Data    []Remittance `json:"data"`
	HasMore bool         `json:"has_more"`
}

type Remittance struct {
	ID             string                 `json:"id"`
	Object         string                 `json:"object"`
	Status         string                 `json:"status"`
	PayerReference string                 `json:"payer_reference"`
	Amount         RemittanceAmount       `json:"amount"`
	Allocations    []RemittanceAllocation `json:"allocations"`
	ReceivedAt     string                 `json:"received_at"`
	ReconciledAt   *string                `json:"reconciled_at"`
}

type RemittanceAmount struct {
	Currency string `json:"currency"`
	Paid     int64  `json:"paid"`
}

type RemittanceAllocation struct {
	ClaimID     string   `json:"claim_id"`
	PaidAmount  int64    `json:"paid_amount"`
	ReasonCodes []string `json:"reason_codes"`
}

type RemittanceInput struct {
	PayerReference string                      `json:"payer_reference"`
	Currency       Currency                    `json:"currency"`
	ReceivedAt     string                      `json:"received_at"`
	Allocations    []RemittanceAllocationInput `json:"allocations"`
}

type RemittanceAllocationInput struct {
	ClaimID     string   `json:"claim_id"`
	PaidAmount  int64    `json:"paid_amount"`
	ReasonCodes []string `json:"reason_codes,omitempty"`
}

type WebhookEvent string

const (
	WebhookSandboxPing               WebhookEvent = "sandbox.ping"
	WebhookPreauthorizationRequested WebhookEvent = "preauthorization.requested"
	WebhookPreauthorizationApproved  WebhookEvent = "preauthorization.approved"
	WebhookPreauthorizationDenied    WebhookEvent = "preauthorization.denied"
	WebhookPreauthorizationExpired   WebhookEvent = "preauthorization.expired"
	WebhookClaimSubmitted            WebhookEvent = "claim.submitted"
	WebhookClaimInformationRequested WebhookEvent = "claim.information_requested"
	WebhookClaimResubmitted          WebhookEvent = "claim.resubmitted"
	WebhookClaimApproved             WebhookEvent = "claim.approved"
	WebhookClaimPartiallyApproved    WebhookEvent = "claim.partially_approved"
	WebhookClaimDenied               WebhookEvent = "claim.denied"
	WebhookClaimSettled              WebhookEvent = "claim.settled"
	WebhookRemittanceReconciled      WebhookEvent = "remittance.reconciled"
)

type WebhookEndpointList struct {
	Object  string            `json:"object"`
	Data    []WebhookEndpoint `json:"data"`
	HasMore bool              `json:"has_more"`
}

type WebhookEndpoint struct {
	ID         string         `json:"id"`
	Object     string         `json:"object"`
	URL        string         `json:"url"`
	Events     []WebhookEvent `json:"events"`
	Status     string         `json:"status"`
	CreatedAt  string         `json:"created_at"`
	DisabledAt *string        `json:"disabled_at"`
}

type WebhookEndpointWithSecret struct {
	ID            string         `json:"id"`
	Object        string         `json:"object"`
	URL           string         `json:"url"`
	Events        []WebhookEvent `json:"events"`
	Status        string         `json:"status"`
	CreatedAt     string         `json:"created_at"`
	DisabledAt    *string        `json:"disabled_at"`
	SigningSecret string         `json:"signing_secret"`
}

type WebhookEndpointInput struct {
	URL    string         `json:"url"`
	Events []WebhookEvent `json:"events"`
}

type WebhookDelivery struct {
	ID             string `json:"id"`
	Object         string `json:"object"`
	Delivered      bool   `json:"delivered"`
	Attempts       int    `json:"attempts"`
	ResponseStatus *int   `json:"response_status"`
}
