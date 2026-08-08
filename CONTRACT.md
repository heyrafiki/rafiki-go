# Contract alignment

| Item | Value |
| --- | --- |
| Status | Source Preview |
| Contract | Heyrafiki API `1.0.0` |
| Source | `heyrafiki/openapi@327e0a70de92771d3930380ce552c00e0ed8fc52` |
| Reviewed | 2026-08-09 |

The SDK implements only operations present in the published OpenAPI document.
Its 30 service methods map to 30 HTTP method and path pairs in that contract.
The operation-surface test fixes those pairs in one reviewable table.

Types are handwritten from the contract. There is no generated code. A contract
update requires:

1. reviewing the OpenAPI diff and compatibility impact;
2. updating types and service methods only where the contract changed;
3. updating the operation-surface and behavior tests;
4. recording the new contract version and commit here; and
5. adding a changelog entry.

The OpenAPI repository remains the authority. This file records provenance; it
does not replace or fork the contract.
