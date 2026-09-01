# Contract alignment

| Item | Value |
| --- | --- |
| Status | Contract aligned |
| Contract | Heyrafiki API `1.0.0` |
| Source | `heyrafiki/contract@62c32d1b99ddded0cfe0baf8ddc57bcbaa764167` |
| Document SHA-256 | `cf2747ba4282bc79e81f69aca66937de78f067266e3dd5651edecdf531863dda` |
| Reviewed | 2026-08-28 |

The SDK implements only operations present in the published OpenAPI document.
Its 31 service methods map to 31 HTTP method and path pairs in that contract.
The operation-surface test fixes those pairs in one reviewable table.

Types are handwritten from the contract. There is no generated code. A contract
update requires:

1. reviewing the OpenAPI diff and compatibility impact;
2. updating types and service methods only where the contract changed;
3. updating the operation-surface and behavior tests;
4. recording the new contract version, commit and document SHA-256 here; and
5. adding a changelog entry.

The contract repository remains the authority. This file records provenance; it
does not replace or fork the contract.
