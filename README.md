<!-- omit in toc -->
# check-path

Go-based tooling to check/verify filesystem paths as part of a Nagios service
check.

[![Latest Release](https://img.shields.io/github/release/atc0005/check-path.svg?style=flat-square)](https://github.com/atc0005/check-path/releases/latest)
[![GoDoc](https://godoc.org/github.com/atc0005/check-path?status.svg)](https://godoc.org/github.com/atc0005/check-path)
[![Validate Codebase](https://github.com/atc0005/check-path/workflows/Validate%20Codebase/badge.svg)](https://github.com/atc0005/check-path/actions?query=workflow%3A%22Validate+Codebase%22)
[![Validate Docs](https://github.com/atc0005/check-path/workflows/Validate%20Docs/badge.svg)](https://github.com/atc0005/check-path/actions?query=workflow%3A%22Validate+Docs%22)
[![Lint and Build using Makefile](https://github.com/atc0005/check-path/workflows/Lint%20and%20Build%20using%20Makefile/badge.svg)](https://github.com/atc0005/check-path/actions?query=workflow%3A%22Lint+and+Build+using+Makefile%22)
[![Quick Validation](https://github.com/atc0005/check-path/workflows/Quick%20Validation/badge.svg)](https://github.com/atc0005/check-path/actions?query=workflow%3A%22Quick+Validation%22)

<!-- omit in toc -->
## Table of Contents

- [Project home](#project-home)
- [Overview](#overview)
- [Features](#features)
- [Changelog](#changelog)
- [Requirements](#requirements)
  - [Building source code](#building-source-code)
  - [Running](#running)
- [Installation](#installation)
- [Configuration](#configuration)
  - [Precedence](#precedence)
  - [Command-line Arguments](#command-line-arguments)
  - [Environment Variables](#environment-variables)
- [Examples](#examples)
- [License](#license)
- [Related projects](#related-projects)
- [References](#references)

## Project home

See [our GitHub repo][repo-url] for the latest code, to file an issue or
submit improvements for review and potential inclusion into the project.

## Overview

This repo provides a Nagios plugin named `check_path` that may be used to
verify the ownership, group, permissions, age or size of specific files or
directories.

While `check_path` is intended to be a Swiss Army knife-able to check for
multiple things, or not, as needed-it probably should not be used to check
*all the things* in a single service check.

Combining too many items together may make the service check harder for others
to follow when responding to an alert. Additionally, checking too many
attributes could slow down the service check enough that it causes the check
to fail from a console-enforced timeout.

For example, instead of checking permissions, owner, group, size and age all
on the same path from the same service check, it is probably a better fit to
check permissions and ownership from one, and maybe size and age together, or
even age and size separated into individual service checks. This way you cover
the spectrum while keeping each service check specific *enough* that any
alerts generated from check results have a common theme.

## Features

- Permissions checks
  - *coming soon*
- Existence checks
  - `CRITICAL` or `WARNING` (as specified) if present
- Age checks
  - `CRITICAL` and `WARNING` thresholds
- Size checks
  - `CRITICAL` and `WARNING` thresholds
- Username checks
  - `CRITICAL` or `WARNING` (as specified) if missing
- Group Name checks
  - `CRITICAL` or `WARNING` (as specified) if missing
- Optional recursive evaluation toggle
- Optional "missing OK" toggle for all checks aside from the "existence"
  checks

## Changelog

See the [`CHANGELOG.md`](CHANGELOG.md) file for the changes associated with
each release of this application.

## Requirements

The following is a loose guideline. Other combinations of Go and operating
systems for building and running tools from this repo may work, but have not
been tested.

### Building source code

- Go 1.14+
- GCC
  - if building with custom options (as the provided `Makefile` does)
- `make`
  - if using the provided `Makefile`

### Running

- Windows 7, Server 2008R2 or later
  - per official [Go install notes][go-docs-install]
- Windows 10 Version 1909
  - tested
- Ubuntu Linux 16.04, 18.04

## Installation

1. [Download][go-docs-download] Go
1. [Install][go-docs-install] Go
   - NOTE: Pay special attention to the remarks about `$HOME/.profile`
1. Clone the repo
   1. `cd /tmp`
   1. `git clone https://github.com/atc0005/check-path`
   1. `cd check-path`
1. Install dependencies (optional)
   - for Ubuntu Linux
     - `sudo apt-get install make gcc`
   - for CentOS Linux
     - `sudo yum install make gcc`
   - for Windows
     - Emulated environments (*easier*)
       - Skip all of this and build using the default `go build` command in
         Windows (see below for use of the `-mod=vendor` flag)
       - build using Windows Subsystem for Linux Ubuntu environment and just
         copy out the Windows binaries from that environment
       - If already running a Docker environment, use a container with the Go
         tool-chain already installed
       - If already familiar with LXD, create a container and follow the
         installation steps given previously to install required dependencies
     - Native tooling (*harder*)
       - see the StackOverflow Question `32127524` link in the
         [References](references.md) section for potential options for
         installing `make` on Windows
       - see the mingw-w64 project homepage link in the
         [References](references.md) section for options for installing `gcc`
         and related packages on Windows
1. Build binaries
   - for the current operating system, explicitly using bundled dependencies
         in top-level `vendor` folder
     - `go build -mod=vendor ./cmd/check_path/`
   - for all supported platforms (where `make` is installed)
      - `make all`
   - for use on Windows
      - `make windows`
   - for use on Linux
     - `make linux`
1. Copy the newly compiled binary from the applicable `/tmp` subdirectory path
   (based on the clone instructions in this section) below and deploy where
   needed.
   - if using `Makefile`
     - look in `/tmp/check-path/release_assets/check_path/`
   - if using `go build`
     - look in `/tmp/check-path/`

## Configuration

### Precedence

The priority order is:

1. Command line flags (highest priority)
1. Environment variables
1. Default settings (lowest priority)

In general, command-line options are the primary way of configuring settings
for this application, but environment variables are also a supported
alternative. Most plugin settings require that a value be specified by the
sysadmin, though some (e.g., logging) have useful defaults.

### Command-line Arguments

- Flags marked as **`required`** must be set via CLI flag or environment
  variable.
- Flags *not* marked as required are for settings where a useful default is
  already defined.
- `critical` and `warning` threshold values are *required* for `age` and
  `size` checks.
- For `username` and `group-name` checks, only one of `critical` or `warning`
  may be specified; specifying both is a configuration error.

| Option                        | Required | Default      | Repeat | Possible                                                                | Description                                                                                                                   |
| ----------------------------- | -------- | ------------ | ------ | ----------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------- |
| `h`, `help`                   | No       | `false`      | No     | `h`, `help`                                                             | Show Help text along with the list of supported flags.                                                                        |
| `emit-branding`               | No       | `false`      | No     | `true`, `false`                                                         | Toggles emission of branding details with plugin status details. This output is disabled by default.                          |
| `log-level`                   | No       | `info`       | No     | `disabled`, `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace` | Log message priority filter. Log messages with a lower level are ignored.                                                     |
| `paths`                       | Yes      | *empty list* | No     | *one or more valid files and directories*                               | List of comma or space-separated paths to process.                                                                            |
| `recursive`                   | No       | `false`      | No     | `true`, `false`                                                         | Perform recursive search into subdirectories.                                                                                 |
| `missing-ok`                  | No       | `false`      | No     | `true`, `false`                                                         | Whether a missing path is considered `OK`. Incompatible with `exists-critical` or `exists-warning` options.                   |
| `age-critical`                | No       | `0`          | No     | `2+` (*minimum 1 greater than warning*)                                 | Assert that age for specified paths is less than specified age in days, otherwise consider state to be `CRITICAL`.            |
| `age-warning`                 | No       | `0`          | No     | `1+` (*minimum of 1*)                                                   | Assert that age for specified paths is less than specified age in days, otherwise consider state to be `WARNING`.             |
| `size-critical`               | No       | `0`          | No     | `2+` (*minimum 1 greater than warning*)                                 | Assert that size for specified paths is less than specified size in bytes, otherwise consider state to be `CRITICAL`.         |
| `size-warning`                | No       | `0`          | No     | `1+` (*minimum of 1*)                                                   | Assert that size for specified paths is less than specified size in bytes, otherwise consider state to be `WARNING`.          |
| `exists-critical`             | No       | `false`      | No     | `true`, `false`                                                         | Assert that specified paths are missing, otherwise consider state to be `CRITICAL`.                                           |
| `exists-warning`              | No       | `false`      | No     | `true`, `false`                                                         | Assert that specified paths are missing, otherwise consider state to be `WARNING`.                                            |
| `username-missing-critical`   | No       | `false`      | No     | *valid username*   (**not supported on Windows**)                       | Assert that specified owner/username is present on all content in specified paths, otherwise consider state to be `CRITICAL`. |
| `username-missing-warning`    | No       | `false`      | No     | *valid username*   (**not supported on Windows**)                       | Assert that specified owner/username is present on all content in specified paths, otherwise consider state to be `WARNING`.  |
| `group-name-missing-critical` | No       | `false`      | No     | *valid group name* (**not supported on Windows**)                       | Assert that specified group name is present on all content in specified paths, otherwise consider state to be `CRITICAL`.     |
| `group-name-missing-warning`  | No       | `false`      | No     | *valid group name* (**not supported on Windows**)                       | Assert that specified group name is present on all content in specified paths, otherwise consider state to be `WARNING`.      |

### Environment Variables

If used, command-line arguments override the equivalent environment variables
listed below. See the [Command-line Arguments](#command-line-arguments) table
for more information.

| Flag Name                     | Environment Variable Name                | Notes | Example (mostly using default values)                     |
| ----------------------------- | ---------------------------------------- | ----- | --------------------------------------------------------- |
| `emit-branding`               | `CHECK_PATH_EMIT_BRANDING`               |       | `CHECK_PATH_EMIT_BRANDING="false"`                        |
| `log-level`                   | `CHECK_PATH_LOG_LEVEL`                   |       | `CHECK_PATH_LOG_LEVEL="info"`                             |
| `paths`                       | `CHECK_PATH_PATHS_LIST`                  |       | `CHECK_PATH_PATHS_LIST="/var/log/apache2 /var/log/samba"` |
| `recursive`                   | `CHECK_PATH_RECURSE`                     |       | `CHECK_PATH_RECURSE=false`                                |
| `missing-ok`                  | `CHECK_PATH_MISSING_OK`                  |       | `CHECK_PATH_MISSING_OK=false`                             |
| `age-critical`                | `CHECK_PATH_AGE_CRITICAL`                |       | `CHECK_PATH_AGE_CRITICAL="2"`                             |
| `age-warning`                 | `CHECK_PATH_AGE_WARNING`                 |       | `CHECK_PATH_AGE_WARNING="1"`                              |
| `size-critical`               | `CHECK_PATH_SIZE_CRITICAL`               |       | `CHECK_PATH_SIZE_CRITICAL="2"`                            |
| `size-warning`                | `CHECK_PATH_SIZE_WARNING`                |       | `CHECK_PATH_SIZE_WARNING="1"`                             |
| `exists-critical`             | `CHECK_PATH_EXISTS_CRITICAL`             |       | `CHECK_PATH_EXISTS_CRITICAL="true"`                       |
| `exists-warning`              | `CHECK_PATH_EXISTS_WARNING`              |       | `CHECK_PATH_EXISTS_WARNING="true"`                        |
| `username-missing-critical`   | `CHECK_PATH_USERNAME_MISSING_CRITICAL`   |       | `CHECK_PATH_USERNAME_MISSING_CRITICAL="ubuntu"`           |
| `username-missing-warning`    | `CHECK_PATH_USERNAME_MISSING_WARNING`    |       | `CHECK_PATH_USERNAME_MISSING_WARNING="ubuntu"`            |
| `group-name-missing-critical` | `CHECK_PATH_GROUP_NAME_MISSING_CRITICAL` |       | `CHECK_PATH_GROUP_NAME_MISSING_CRITICAL="adm"`            |
| `group-name-missing-warning`  | `CHECK_PATH_GROUP_NAME_MISSING_WARNING`  |       | `CHECK_PATH_GROUP_NAME_MISSING_WARNING="adm"`             |

## Examples

- placeholder

## License

```license
MIT License

Copyright (c) 2020 Adam Chalkley

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

## Related projects

- <https://github.com/atc0005/elbow>
- <https://github.com/atc0005/go-nagios>
- <https://github.com/atc0005/check-mail>
- <https://github.com/atc0005/check-cert>

## References

- <https://github.com/rs/zerolog>
- <https://github.com/alexflint/go-arg>
- <https://github.com/phayes/permbits>
- <https://github.com/atc0005/go-nagios>
- <https://nagios-plugins.org/doc/guidelines.html>

<!-- Footnotes here  -->

[repo-url]: <https://github.com/atc0005/check-path>  "This project's GitHub repo"

[go-docs-download]: <https://golang.org/dl>  "Download Go"

[go-docs-install]: <https://golang.org/doc/install>  "Install Go"

<!-- []: PLACEHOLDER "DESCRIPTION_HERE" -->
