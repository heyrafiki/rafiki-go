# Release process

The module is currently a Source Preview. Do not create a tag until every gate
below has passed.

## Release gates

1. Confirm the SDK matches the latest published version 1 OpenAPI contract.
2. Review the public boundary, examples, licence and repository files.
3. Run the full Go 1.25 and Go 1.26 CI matrix, race tests, vet and `govulncheck`.
4. Require CODEOWNERS approval and resolved review conversations.
5. Confirm branch protection, secret scanning, dependency review and private
   vulnerability reporting are enabled.
6. Update `CHANGELOG.md` and remove the Source Preview publication disclaimer.
7. Create a signed `v0.x.y-beta.n` tag from the reviewed commit.
8. Create a GitHub release with compatibility notes and the contract commit.
9. Verify module discovery and documentation before announcing availability.

Beta releases may add operations and fields in line with the version 1 contract.
Breaking SDK changes require a clear migration note. Stable compatibility begins
at `v1.0.0`.
