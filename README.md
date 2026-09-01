# Heyrafiki Go SDK

Server-side Go client for the Heyrafiki API.

## What this SDK is for

Use this SDK to call the Heyrafiki API from a Go service. It provides typed
requests, predictable errors and idempotency support for retried writes. The
platform applies access rules; clinical and financial decisions remain with the
accountable people and organizations.

The client covers every operation published in the version 1 OpenAPI contract:
Practitioners, Bookings, Sessions, eligibility, Coverage, pre-authorization,
Claims, Claim valuation, remittance and Webhooks.

## First request

```go
package main

import (
	"context"
	"log"
	"os"

heyrafiki "github.com/heyrafiki/rafiki-go"
)

func main() {
	client, err := heyrafiki.NewClient(os.Getenv("HEYRAFIKI_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	practitioners, err := client.Practitioners.List(
		context.Background(),
		&heyrafiki.ListOptions{Limit: 5},
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("found %d Practitioners", len(practitioners.Data))
}
```

Keep API keys on the server. Sandbox keys return synthetic data.

Build and test the client directly from its public source:

```bash
git clone https://github.com/heyrafiki/rafiki-go.git
cd rafiki-go
go test -race ./...
```

## Idempotent writes

The API contract requires a caller-owned idempotency key for replay-safe writes.

```go
booking, err := client.Bookings.Create(ctx, heyrafiki.BookingInput{
	PractitionerID: "prc_2481",
	StartsAt:       "2026-08-12T07:00:00Z",
	EndsAt:         "2026-08-12T08:00:00Z",
	Format:         heyrafiki.CareFormatOnline,
	PaymentSource:  heyrafiki.PaymentSourceCovered,
}, heyrafiki.WriteOptions{IdempotencyKey: "booking-20260812-001"})
```

Reads and writes with an idempotency key use bounded retries for transport
failures, `408`, `429`, `500`, `502`, `503` and `504` responses. Webhook
registration and test delivery are not retried because the contract does not
accept an idempotency key for those operations.

## Claim valuation

`Claims.Valuation` reproduces what a Claim was worth at a point in time. The
cutoff is explicit, and Heyrafiki includes only the facts whose business time and
knowledge time both fall on or before it, so the same cutoff returns the same
answer however much later you ask.

```go
valuation, err := client.Claims.Valuation(ctx, "clm_demo_001", heyrafiki.ClaimValuationOptions{
	ValuationAt: "2026-08-12T09:00:00Z",
})
if err != nil {
	log.Fatal(err)
}

for _, event := range valuation.Events {
	log.Printf("%d %s effective %s recorded %s",
		event.Sequence, event.Type, event.EffectiveAt, event.RecordedAt)
}
```

`EffectiveAt` is when the underlying fact took effect. `RecordedAt` is when
Heyrafiki persisted it. `Amount.Outstanding` is payer liability less the
settlement an authorized observation confirmed, and it is `nil` when the contract
reports it as unknown. `ValuationAt` must be an RFC 3339 timestamp; the SDK
rejects anything else rather than sending a cutoff the API cannot reproduce.

## Errors

```go
var apiError *heyrafiki.APIError
if errors.As(err, &apiError) {
	log.Printf("status=%d code=%s request_id=%s", apiError.StatusCode, apiError.Code, apiError.RequestID)
}
```

Use the request ID when tracing a failed call. The SDK never adds request or
response bodies to error strings.

## Compatibility

- Go 1.25 and 1.26
- Heyrafiki API version 1
- Bearer authentication by default; `x-api-key` is available through
  `WithAuthStyle(heyrafiki.AuthAPIKey)`
- Standard library only

The SDK tracks additive changes to the published contract. Breaking API changes
use a new major API version and a documented migration period.

## Contract provenance

This client was reviewed against
[`heyrafiki/contract`](https://github.com/heyrafiki/contract) contract `1.0.0`,
commit `62c32d1b99ddded0cfe0baf8ddc57bcbaa764167`, on 2026-08-28. The client is
handwritten; no generated source is included. See [`CONTRACT.md`](./CONTRACT.md).

## Develop

```bash
gofmt -w .
go test -race ./...
go vet ./...
go run golang.org/x/vuln/cmd/govulncheck@v1.1.4 ./...
```

## Resources

- [Documentation](https://docs.heyrafiki.space)
- [API contract](https://github.com/heyrafiki/contract)
- [Open insurance assurance benchmark](https://github.com/heyrafiki/proving-ground)
- [Webhooks](https://docs.heyrafiki.space/webhooks)
- [Security](./SECURITY.md)

## License

Licensed under the [Apache License 2.0](./LICENSE).
