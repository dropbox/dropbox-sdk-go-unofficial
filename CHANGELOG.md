# Changelog

Notable changes to this SDK are documented here. Historical entries are
summarized from the repository history.

## Unreleased

### Added

- Added the `dropbox/contenthash` package for computing Dropbox API
  `content_hash` values.
- Added automatic `content_hash` population for seekable file upload readers
  on upload routes whose arguments include `ContentHash`.

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
