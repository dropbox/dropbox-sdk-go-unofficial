# Changelog

Notable changes to this SDK are documented here. Historical entries are
summarized from the repository history.

## Unreleased

### Added

- Added the `dropbox/oauth` package with PKCE authorization URL, code
  exchange, and refresh-token helpers.
- Added web redirect PKCE OAuth helpers with caller-owned CSRF validation.
- Added `dropbox.Config.TokenSource` for refreshable OAuth token sources.

## v6.3.0 - 2026-07-04

### Added

- Added the `dropbox/contenthash` package for computing Dropbox API
  `content_hash` values.
- Added automatic `content_hash` population for seekable file upload readers
  on upload routes whose arguments include `ContentHash`.
- Added `go vet` step to test workflow for platform-specific issue detection.
- Added `.golangci.yml` configuration with govet enable-all.
- Added godoc comment to `SDKInternalError`.

### Changed

- Added Go module caching (`cache-dependency-path`) to test and lint workflows.
- Added CodeQL scanning for GitHub Actions workflows.
- Fixed CodeQL Go analysis to use explicit module build.

### Security

- Added explicit least-privilege permissions to GitHub Actions workflows.

## v6.2.0 - 2026-06-30

### Added

- Added `context.Context` support to generated namespace clients.

### Changed

- Bumped `golang.org/x/oauth2` to a Go 1.23-compatible version.
- Updated GitHub Actions dependencies, including checkout, setup-go, CodeQL,
  and golangci-lint actions.
- Fixed lint issues reported by golangci-lint v9.

## v6.1.0 - 2026-06-29

### Fixed

- Fixed timestamp serialization format.

### Changed

- Updated the Dropbox API spec and regenerated the SDK.
- Updated README usage examples for polymorphic responses and timestamps.
- Replaced deprecated `io/ioutil` usage with `io` package equivalents.
- Bumped dependencies and CI configuration.

## v6.0.5 - 2022-10-03

### Changed

- Updated the Dropbox API spec and regenerated the SDK.
- Added generator support for map types.
- Updated GitHub Actions dependencies, including checkout, setup-go, CodeQL,
  and golangci-lint actions.
