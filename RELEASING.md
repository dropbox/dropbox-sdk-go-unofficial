# Releasing

Repository administrators can publish the Go SDK from the GitHub Actions UI.
The release workflow builds, vets, and tests the module before it creates the
version tag and GitHub Release. Creating the tag publishes this public Go
module; the workflow then explicitly requests the version from the public Go
module proxy, which also requests indexing by pkg.go.dev.

## Prepare the release

1. Merge the release changes into `master`. Update `CHANGELOG.md` and any SDK
   version metadata as part of the release pull request.
2. Choose a new stable semantic version such as `6.7.0`. The corresponding
   `v6.7.0` tag must not already exist, and the repository must contain the
   matching `v6/go.mod` module.
3. Confirm all required checks on the latest `master` commit have passed.

## Run the release workflow

1. Open **Actions**, select **Release**, and choose **Run workflow**.
2. Select the `master` branch.
3. Enter the version without a `v` prefix and run the workflow.

The workflow rejects non-administrators, non-`master` refs, stale `master`
commits, invalid versions, and existing tags. Releases are serialized so two
versions cannot be published concurrently. After the build and tests pass, it
creates `vX.Y.Z` at the selected commit and publishes a GitHub Release with
generated release notes.

The Release workflow calls the reusable Documentation workflow directly after
publication. This is intentional: releases created with `GITHUB_TOKEN` do not
normally trigger other workflows. Documentation also listens for
`release.published` so releases created by a user or GitHub App keep the
release-event publishing path.

If proxy or documentation indexing fails after the release is created, do not
move or recreate the tag. Open **Actions**, select **Documentation**, choose
**Run workflow**, and enter the same version (with or without the `v` prefix)
to retry indexing.
