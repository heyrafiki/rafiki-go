# Changelog

All notable changes to this module are recorded here.

## 0.1.0-beta.2

- Added `Claims.Valuation` for `GET /claims/{claim_id}/valuation`.
- Added the `ClaimValuation`, `ClaimValuationAmount` and `ClaimValuationEvent` types.
- Required an RFC 3339 `ValuationAt` cutoff and sent it unchanged on every retry.
- Realigned contract provenance with the revision that publishes Claim valuation.

## 0.1.0-beta.1

- Added typed clients for every operation in the Heyrafiki API version 1 contract.
- Added bearer and `x-api-key` authentication.
- Added bounded retries for reads and idempotent writes.
- Added typed API errors, request identifiers and response-size limits.
- Added contract-surface, retry, authentication and disclosure-boundary tests.
