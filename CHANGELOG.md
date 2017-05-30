# CHANGELOG

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

## [x.y.z] - YYYY-MM-DD (Unreleased)

### Improvements

- The `grafton test` command now validates the credential names it receives.
- The main `grafton` package now exports a function `ValidCredentialName` for
  testing whether or not a given credential name is valid.
- Dependencies for bootstrapping the build are now vendored using glide.

## [0.8.0] - 2017-05-14

### Improvements

- The grafton client now returns `grafton.Error` if the received response
  conforms to the responses expected in the `provider.yaml` swagger
  specification.
- The main `grafton` package now exports a function `IsFatal` which returns
  whether or not an error is considered to be fatal to a provision, plan
  change, or deprovisioning flow.

## [0.7.0] - 2017-04-28

### Improvements

- The main `grafton` package now exports a function `CreateSsoURL` for deriving
  an SSO URL.

## [0.6.6] - 2017-04-27

### Fixed

- A backwards incompatible change was accidently introduced in `v0.6.5`, `/v1`
  would always be appended to the url, instead, `/v1` will only be appended if
  the provided url does not end with `/v1`.

## [0.6.5] - 2017-04-27

### Fixed

- If a trailing slash is provided for the base url, Grafton will no longer
  generate urls with duplicate `/` between path segments.
- Update README to include working example using `grafton test`
- `/v1` will always be prepended to the base url given to Grafton

## [0.6.4] - 2017-04-20

### Fixed

- A resource can now be queried *during* the provisioning flow instead of after
  it's completed.

## [0.6.3] - 2017-04-19

### Fixed

- A failing ErrorCase will not cause the rest of the test runs to abort early.
- Make the error messages for SSO error cases match the status code they're
  looking for.
- Don't require a message on 204 responses, which won't have a body.

## [0.6.2] - 2017-04-13

### Added

- Initial release from our public git repository.
- Detect, display, and fail on invalid status code responses during grafton
  test.

### Removed

- The credentials flag on grafton test, which was not used, has been removed.

### Fixed

- Correct the RunsInside logic, so dependent features are not run when an outer
  feature fails.
- Display cleaner error messages when grafton fails to connect to your provider
  implementation.
- Handle and report on missing `message` fields in provider responses.
