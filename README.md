# Heyrafiki Go SDK

Server-side Go client for the Heyrafiki API.

> Source Preview. The API contract and this SDK are in beta. Use Sandbox projects.

The client covers every operation published in the version 1 OpenAPI contract:
Practitioners, Bookings, Sessions, eligibility, Coverage, pre-authorization,
Claims, remittance and Webhooks.

## First request

```go
package main

import (
	"context"
	"log"
	"os"

	heyrafiki "github.com/heyrafiki/heyrafiki-go"
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

This source preview is not a tagged Go module release. Clone it locally for
evaluation; a release tag and module availability will be announced separately.

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

This source preview was reviewed against
[`heyrafiki/openapi`](https://github.com/heyrafiki/openapi) contract `1.0.0`,
commit `327e0a70de92771d3930380ce552c00e0ed8fc52`, on 2026-08-09. The client is
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
- [API contract](https://github.com/heyrafiki/openapi)
- [Webhooks](https://docs.heyrafiki.space/webhooks)
- [Security](./SECURITY.md)

## License

Licensed under the [Apache License 2.0](./LICENSE).
