# Changelog

## Overview

All notable changes to this project will be documented in this file.

The format is based on [Keep a
Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to
[Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Please [open an issue](https://github.com/atc0005/check-path/issues) for any
deviations that you spot; I'm still learning!.

## Types of changes

The following types of changes will be recorded in this file:

- `Added` for new features.
- `Changed` for changes in existing functionality.
- `Deprecated` for soon-to-be removed features.
- `Removed` for now removed features.
- `Fixed` for any bug fixes.
- `Security` in case of vulnerabilities.

## [Unreleased]

- placeholder

## [v0.1.2] - 2020-11-15

### Added

- Support ignoring paths (files, directories, subdirectories)
- Support minimum size checks in addition to the existing maximum size checks

### Changed

- Statically linked binary release
  - Built using Go 1.15.5
  - Windows
    - x86
    - x64
  - Linux
    - x86
    - x64

- Dependencies
  - `atc0005/go-nagios`
    - `v0.5.1` to `v0.5.2`

- Remove temporary workaround for swallowed panics
  - see `atc0005/go-nagios` `v0.5.2` release notes

### Fixed

- State change logic triggers when *reaching* thresholds in addition to when
  crossing them

- Fix doc comment breadcrumb URL

- Configuration validation used direct field access when getter methods were
  sufficient
  - may require further review in the future

- fail-fast logic appears to be applied regardless of flag use

- Documentation
  - update examples to reflect recent changes
  - expand "Known issues" section to better cover potentially unexpected
    behavior of combining `fail-fast` with other check options
  - explicitly note that permissions check support is not yet available (GH-6)

## [v0.1.1] - 2020-11-06

### Added

- `fail-fast` flag
  - allows toggling the `v0.1.0` behavior of quickly failing with
    indeterminate `WARNING` or `CRITICAL` state as soon as a non-`OK` state is
    detected
  - see README for more information

### Changed

- Statically linked binary release
  - Built using Go 1.15.4
  - Windows
    - x86
    - x64
  - Linux
    - x86
    - x64

- Dependencies
  - `actions/checkout`
    - `v2.3.3` to `v2.3.4`

### Fixed

- WARNING thresholds (may) trigger before CRITICAL thresholds, even if
  CRITICAL threshold would have a match
  - see new `fail-fast` flag, README for details

## [v0.1.0] - 2020-11-02

### Added

Initial release!

This release provides an early version of a Nagios plugin used to monitor
attributes of one or many specified paths. The intention is to provide a
multi-purpose or "Swiss Army Knife" tool that is capable of monitoring many
different attributes, though flexible enough to easily monitor just one.

- Statically linked binary release
  - Built using Go 1.15.3
  - Windows
    - x86
    - x64
  - Linux
    - x86
    - x64

#### Available now

Currently, monitored attributes include:

- `age`
  - `CRITICAL` and `WARNING` thresholds
- `size`
  - `CRITICAL` and `WARNING` thresholds
- `username`
  - `CRITICAL` or `WARNING` (as specified) if missing
  - **NOTE**: this check is not supported on Windows
- `group name`
  - `CRITICAL` or `WARNING` (as specified) if missing
  - **NOTE**: this check is not supported on Windows
- `exists`
  - `CRITICAL` or `WARNING` (as specified) if present

Optional support for ignoring missing files (does not apply to the `exists`
checks) and recursive evaluation is available, but disabled by default.

#### Coming "Soon"

- Permissions checks

[Unreleased]: https://github.com/atc0005/check-path/compare/v0.1.2...HEAD
[v0.1.2]: https://github.com/atc0005/check-path/releases/tag/v0.1.2
[v0.1.1]: https://github.com/atc0005/check-path/releases/tag/v0.1.1
[v0.1.0]: https://github.com/atc0005/check-path/releases/tag/v0.1.0
