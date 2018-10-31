# CHANGELOG

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

## [Unreleased]

## [0.14.2] - 2018-10-31

### Added

- releases are handled by [promulgate](https://www.github.com/manifoldco/promulgate)

### Added

## [0.14.1] - 2018-10-25

### Added

- `verbose` mode now correctly shows all HTTP header fields and body

## [0.14.0] - 2018-10-09

### Added

- `grafton serve` runs a local version of the Connector API
- Plan features can be passed during provision with `--features` flag as JSON
- Update go-manifold
- Add mini-marketplace to Grafton as a new command, with provision, deprovision and SSO
- Support latest go versions

### Fixed

- Provisioning the same resource test passes if the provider chooses to respond with a callback.

## [0.13.0] - 2018-04-13

### Added

- Resource features: the resource provisioning endpoint now receives additional data describing the
  selected values for customizable features when provisioning the resource.
- The change plan endpoint also receives new feature data for resource plan changes
- The acceptance tests have been updated to allow for feature data to be supplied as a JSON object string through
  `plan-features/PLAN_FEATURES` and `new-plan-features/NEW_PLAN_FEATURES`
- Resource measures spec: allows Manifold to pull information about a resource
  features usage
- Resource measures acceptance test
- Add `label` field to resource

### Fixed

- The Grafton mock connector API now supports a JSON message body to POST /oauth/tokens, as described in the provider documentation

## [0.12.0] - 2017-11-23

### Added

- `credentials list` cmd for providers to list their OAuth Credentials
- `credentials delete` cmd for providers to delete their OAuth Credentials

## [0.11.0] - 2017-11-03

### Added

- `credentials rotate` cmd for providers to rotate their OAuth Credentials

## [0.10.3] - 2017-10-27

### Fixed

- Rename dependency from `Sirupsen/logrus` to `sirupsen/logrus`

## [0.10.2] - 2017-07-03

### Fixed

- Grafton no longer errors if a `message` is not provided with a 201 Created
  response to a credential provisioning request.

## [0.10.1] - 2017-06-26

### Fixed

- Grafton no longer returns an error if a 204 Response is returned from a Plan
  Change request.

## [0.10.0] - 2017-06-22

### Added

- Introduced the cleanup acceptance test set, which tests to ensure any
  half-created resources can be cleaned up.
- Added a teardown step to the resize acceptance tests, ensuring a resource can
  be resized back to it's original plan.
- The Grafton Client now returns an `ErrMissingMsg` on a `200`, `201`, or `202`
  request missing a `message` property.
- The Grafton client now differentiates between an error and a message returned
  from the provider, a message will be returned if a provider sends a valid
  message despite whether or not it was an error response.
- The Grafton client now takes a `logrus.Entry`, when provided Grafton will
  push error logs to the provided logger.

## [0.9.0] - 2017-06-01

### Added

- The `grafton test` command now validates the credential names it receives.
- The main `grafton` package now exports a function `ValidCredentialName` for
  testing whether or not a given credential name is valid.
- Dependencies for bootstrapping the build are now vendored using glide.

### Fixed

- A bug has been fixed which resulted in all of the error acceptance tests
  failing despite the response returned from the provider.

## [0.8.0] - 2017-05-14

### Added

- The grafton client now returns `grafton.Error` if the received response
  conforms to the responses expected in the `provider.yaml` swagger
  specification.
- The main `grafton` package now exports a function `IsFatal` which returns
  whether or not an error is considered to be fatal to a provision, plan
  change, or deprovisioning flow.

## [0.7.0] - 2017-04-28

### Added

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
