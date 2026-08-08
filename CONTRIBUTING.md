# Contributing

Issues and pull requests are welcome.

## Before opening a change

- Keep behavior within the published [OpenAPI contract](https://github.com/heyrafiki/openapi).
- Use synthetic data in tests, examples and issue reports.
- Never include personal, health, payment or production credential data.
- Add or update tests for every behavior change.
- Record consumer-visible changes in `CHANGELOG.md`.

## Checks

Use Go 1.25 or newer.

```bash
gofmt -w .
go test -race ./...
go vet ./...
go run golang.org/x/vuln/cmd/govulncheck@v1.1.4 ./...
```

The formatted diff must remain clean after the checks.

## Pull requests

Describe the contract surface changed, compatibility impact and checks run.
Maintainer review is required before merge. Security reports follow
[`SECURITY.md`](./SECURITY.md), not public issues.

By contributing, you agree that your contribution is licensed under Apache-2.0.
The organization [Code of Conduct](https://github.com/heyrafiki/.github/blob/main/CODE_OF_CONDUCT.md)
applies.
