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

## [v0.1.10] - 2022-05-06

### Overview

- Dependency updates
- built using Go 1.17.9
  - Statically linked
  - Windows (x86, x64)
  - Linux (x86, x64)

### Changed

- Dependencies
  - `Go`
    - `1.17.7` to `1.17.9`

## [v0.1.9] - 2022-03-03

### Overview

- Dependency updates
- CI / linting improvements
- built using Go 1.17.7
  - Statically linked
  - Windows (x86, x64)
  - Linux (x86, x64)

### Changed

- Dependencies
  - `Go`
    - `1.17.6` to `1.17.7`
  - `alexflint/go-arg`
    - `v1.4.2` to `v1.4.3`
  - `actions/checkout`
    - `v2.4.0` to `v3`
  - `actions/setup-node`
    - `v2.5.1` to `v3`

- (GH-110) Expand linting GitHub Actions Workflow to include `oldstable`,
  `unstable` container images
- (GH-111) Switch Docker image source from Docker Hub to GitHub Container
  Registry (GHCR)

### Fixed

- (GH-113) var-declaration: should omit type string from declaration of var
  (revive)

## [v0.1.8] - 2022-01-25

### Overview

- Dependency updates
- built using Go 1.17.6
  - Statically linked
  - Windows (x86, x64)
  - Linux (x86, x64)

### Changed

- Dependencies
  - `Go`
    - `1.16.12` to `1.17.6`
    - (GH-106) Update go.mod file, canary Dockerfile to reflect current
      dependencies
  - `atc0005/go-nagios`
    - `v0.8.1` to `v0.8.2`

## [v0.1.7] - 2021-12-29

### Overview

- Dependency updates
- built using Go 1.16.12
  - Statically linked
  - Windows (x86, x64)
  - Linux (x86, x64)

### Changed

- Dependencies
  - `Go`
    - `1.16.10` to `1.16.12`
  - `rs/zerolog`
    - `v1.26.0` to `v1.26.1`
  - `actions/setup-node`
    - `v2.4.1` to `v2.5.1`

## [v0.1.6] - 2021-11-10

### Overview

- Dependency updates
- built using Go 1.16.10
  - Statically linked
  - Windows (x86, x64)
  - Linux (x86, x64)

### Changed

- Dependencies
  - `Go`
    - `1.16.7` to `1.16.10`
  - `atc0005/go-nagios`
    - `v0.7.0` to `v0.8.1`
  - `rs/zerolog`
    - `v1.23.0` to `v1.26.0`
  - `actions/checkout`
    - `v2.3.4` to `v2.4.0`
  - `actions/setup-node`
    - `v2.4.0` to `v2.4.1`

- (GH-82) Lock Go version to the latest "oldstable" series

### Fixed

- (GH-87) Update build tags for Go 1.17 compatibility

## [v0.1.5] - 2021-08-09

### Overview

- Dependency updates
- built using Go 1.16.7
  - Statically linked
  - Windows (x86, x64)
  - Linux (x86, x64)

### Changed

- Dependencies
  - `Go`
    - `1.16.6` to `1.16.7`
  - `actions/setup-node`
    - updated from `v2.2.0` to `v2.4.0`

## [v0.1.4] - 2021-07-19

### Overview

- Dependency updates
- Minor fixes
- Built using Go 1.16.6
  - **Statically linked**
  - Windows (x86, x64)
  - Linux (x86, x64)

### Added

- Add "canary" Dockerfile to track stable Go releases, serve as a reminder to
  generate fresh binaries

### Changed

- Dependencies
  - `Go`
    - `1.15.8` to `1.16.6`
  - `atc0005/go-nagios`
    - `v0.6.0` to `v0.6.1`
  - `rs/zerolog`
    - `v1.21.0` to `v1.23.0`
  - `alexflint/go-arg`
    - `v1.3.0` to `v1.4.2`
  - `actions/setup-node`
    - `v2.1.5` to `v2.2.0`
    - update `node-version` value to always use latest LTS version instead of
      hard-coded version

### Fixed

- Documentation
  - Incorrect flag name

## [v0.1.3] - 2021-02-21

### Overview

- Dependency updates
- built using Go 1.15.8

### Changed

- Swap out GoDoc badge for pkg.go.dev badge

- dependencies
  - `go.mod` Go version
    - updated from `1.14` to `1.15`
  - built using Go 1.15.8
    - **Statically linked**
    - Windows (x86, x64)
    - Linux (x86, x64)
  - `atc0005/go-nagios`
    - updated from `v0.5.2` to `v0.6.0`
  - `actions/setup-node`
    - `v2.1.2` to `v2.1.4`

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

[Unreleased]: https://github.com/atc0005/check-path/compare/v0.1.10...HEAD
[v0.1.10]: https://github.com/atc0005/check-path/releases/tag/v0.1.10
[v0.1.9]: https://github.com/atc0005/check-path/releases/tag/v0.1.9
[v0.1.8]: https://github.com/atc0005/check-path/releases/tag/v0.1.8
[v0.1.7]: https://github.com/atc0005/check-path/releases/tag/v0.1.7
[v0.1.6]: https://github.com/atc0005/check-path/releases/tag/v0.1.6
[v0.1.5]: https://github.com/atc0005/check-path/releases/tag/v0.1.5
[v0.1.4]: https://github.com/atc0005/check-path/releases/tag/v0.1.4
[v0.1.3]: https://github.com/atc0005/check-path/releases/tag/v0.1.3
[v0.1.2]: https://github.com/atc0005/check-path/releases/tag/v0.1.2
[v0.1.1]: https://github.com/atc0005/check-path/releases/tag/v0.1.1
[v0.1.0]: https://github.com/atc0005/check-path/releases/tag/v0.1.0
