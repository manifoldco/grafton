# CHANGELOG

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

## [Unreleased]

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
