<!-- omit in toc -->
# check-path

Go-based tooling to check/verify filesystem paths as part of a Nagios service
check.

[![Latest Release](https://img.shields.io/github/release/atc0005/check-path.svg?style=flat-square)](https://github.com/atc0005/check-path/releases/latest)
[![Go Reference](https://pkg.go.dev/badge/github.com/atc0005/check-path.svg)](https://pkg.go.dev/github.com/atc0005/check-path)
[![go.mod Go version](https://img.shields.io/github/go-mod/go-version/atc0005/check-path)](https://github.com/atc0005/check-path)
[![Lint and Build](https://github.com/atc0005/check-path/actions/workflows/lint-and-build.yml/badge.svg)](https://github.com/atc0005/check-path/actions/workflows/lint-and-build.yml)
[![Project Analysis](https://github.com/atc0005/check-path/actions/workflows/project-analysis.yml/badge.svg)](https://github.com/atc0005/check-path/actions/workflows/project-analysis.yml)

<!-- omit in toc -->
## Table of Contents

- [Project home](#project-home)
- [Overview](#overview)
- [Features](#features)
- [Known issues](#known-issues)
  - [Permissions checks](#permissions-checks)
  - [`fail-fast` option](#fail-fast-option)
    - [Indeterminate exit state](#indeterminate-exit-state)
    - [Minimum and Maximum size checks](#minimum-and-maximum-size-checks)
- [Changelog](#changelog)
- [Requirements](#requirements)
  - [Building source code](#building-source-code)
  - [Running](#running)
- [Installation](#installation)
  - [From source](#from-source)
  - [Using release binaries](#using-release-binaries)
- [Configuration](#configuration)
  - [Precedence](#precedence)
  - [Command-line Arguments](#command-line-arguments)
  - [Environment Variables](#environment-variables)
- [Examples](#examples)
  - [Help output](#help-output)
  - [Existence check](#existence-check)
    - [`CRITICAL`](#critical)
    - [`WARNING`](#warning)
  - [Age check](#age-check)
    - [`error examining path`](#error-examining-path)
    - [`CRITICAL`](#critical-1)
    - [`WARNING`](#warning-1)
  - [Size check](#size-check)
    - [Maximum file size](#maximum-file-size)
      - [`CRITICAL`](#critical-2)
      - [`WARNING`, enable `fail-fast` option](#warning-enable-fail-fast-option)
    - [Minimum Size check](#minimum-size-check)
      - [`OK`](#ok)
      - [`CRITICAL`](#critical-3)
      - [`WARNING`](#warning-2)
      - [`CRITICAL`, enable `fail-fast` behavior](#critical-enable-fail-fast-behavior)
  - [Username check](#username-check)
    - [`CRITICAL`](#critical-4)
    - [`WARNING`](#warning-3)
  - [Group Name check](#group-name-check)
    - [`OK`](#ok-1)
    - [`CRITICAL`](#critical-5)
    - [`WARNING`](#warning-4)
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

- Existence checks
  - `CRITICAL` or `WARNING` (as specified) if present
- Age checks
  - `CRITICAL` and `WARNING` thresholds
- Size checks
  - minimum `CRITICAL` and `WARNING` thresholds
    - e.g., "path required to be X size or larger"
  - maximum `CRITICAL` and `WARNING` thresholds
    - e.g., "path required to be X size or smaller"
- Username checks
  - `CRITICAL` or `WARNING` (as specified) if missing
  - **NOTE**: this check is not supported on Windows
- Group Name checks
  - `CRITICAL` or `WARNING` (as specified) if missing
  - **NOTE**: this check is not supported on Windows
- Optional recursive evaluation toggle
- Optional "missing OK" toggle for all checks aside from the "existence"
  checks
- Optional exclusion of specific paths from evaluation
  - NOTE: This does not apply to "existence" checks
- Optional "fail fast" behavior in an effort to avoid I/O churn over deep
  paths
  - see [Known issues](#known-issues) for potential issues with this option

## Known issues

### Permissions checks

Not included yet. Planned addition as part of the work for GH-6.

### `fail-fast` option

The `fail-fast` option allows us to evaluate the results of a check
*immediately* instead of waiting for each specified path to be completely
crawled (with or without recursion as specified). This optimization helps to
avoid unnecessary processing of content once a non-OK state is confirmed. This
benefit doesn't come without cost however; as outlined in the subsections
below, this setting may produce unexpected results if the sysadmin is not
aware of how this setting interacts with other specified options.

#### Indeterminate exit state

If using the "early exit" behavior provided by the `fail-fast` flag, this
plugin will exit ASAP once a non-OK state is determined, regardless of whether
the first non-OK state is `CRITICAL` or `WARNING`. Receiving a `WARNING` state
for a path with files/directories which warrant a `CRITICAL` state result
could be confusing to troubleshoot, which is why this behavior is not enabled
by default. If you enable this option, just be aware of the tradeoff.

#### Minimum and Maximum size checks

When the `fail-fast` option is combined with the logic behind the minimum size
(`size-min-critical` and `size-min-warning`) or maximum size
(`size-max-critical` and `size-max-warning`) checks, a state change is
triggered upon finding any **single** file that does not meet the specified
thresholds.

That said, the state change for crossing specified thresholds for maximum size
is likely to be expected behavior, whereas the state change for minimum size
thresholds (due to the use of `fail-fast`) may not be expected.

Scenario:

- you specify 10 MB (in bytes) as the minimum acceptable size value,
- this plugin (seemingly immediately) encounters a file smaller than that
  - e.g., 0 bytes, 500 KB, 9.99 MB

Result:

Depending on the specified threshold values, either a `WARNING` or `CRITICAL`
state change is triggered. This occurs even if the combined size of all files
in a specified path safely clears the minimum file threshold values.

Because of this behavior, this combination of options is probably best
reserved for monitoring paths receiving file uploads or automatically
generated files/reports. In those scenarios, very small (or zero byte) files
are usually indicators of a failed task, so having a state change for those
cases is both helpful and not surprising.

Outside of those scenarios, combining `size-min-critical` or
`size-min-warning` with `fail-fast` is likely to produce unexpected results.

## Changelog

See the [`CHANGELOG.md`](CHANGELOG.md) file for the changes associated with
each release of this application.

## Requirements

The following is a loose guideline. Other combinations of Go and operating
systems for building and running tools from this repo may work, but have not
been tested.

### Building source code

- Go
  - see this project's `go.mod` file for *preferred* version
  - this project tests against [officially supported Go
    releases][go-supported-releases]
    - the most recent stable release (aka, "stable")
    - the prior, but still supported release (aka, "oldstable")
- GCC
  - if building with custom options (as the provided `Makefile` does)
- `make`
  - if using the provided `Makefile`

### Running

- Windows 10
- Ubuntu Linux 18.04+

## Installation

### From source

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

**NOTE**: Depending on which `Makefile` recipe you use the generated binary
may be compressed and have an `xz` extension. If so, you should decompress the
binary first before deploying it (e.g., `xz -d check_path-linux-amd64.xz`).

### Using release binaries

1. Download the [latest
   release](https://github.com/atc0005/check-path/releases/latest) binaries
1. Decompress binaries
   - e.g., `xz -d check_path-linux-amd64.xz`
1. Deploy
   - Place `check_path` alongside your other Nagios plugins
     - e.g., `/usr/lib/nagios/plugins/` or `/usr/lib64/nagios/plugins/`

**NOTE**:

DEB and RPM packages are provided as an alternative to manually deploying
binaries.

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

| Option                        | Required | Default      | Repeat | Possible                                                                | Description                                                                                                                                                                                      |
| ----------------------------- | -------- | ------------ | ------ | ----------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `h`, `help`                   | No       | `false`      | No     | `h`, `help`                                                             | Show Help text along with the list of supported flags.                                                                                                                                           |
| `emit-branding`               | No       | `false`      | No     | `true`, `false`                                                         | Toggles emission of branding details with plugin status details. This output is disabled by default.                                                                                             |
| `log-level`                   | No       | `info`       | No     | `disabled`, `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace` | Log message priority filter. Log messages with a lower level are ignored.                                                                                                                        |
| `paths`                       | Yes      | *empty list* | No     | *one or more valid files and directories*                               | List of comma or space-separated paths to check.                                                                                                                                                 |
| `ignore`                      | No       | *empty list* | No     | *one or more valid files and directories*                               | List of comma or space-separated paths to ignore. Does not apply to existence checks.                                                                                                            |
| `recurse`                     | No       | `false`      | No     | `true`, `false`                                                         | Perform recursive search into subdirectories.                                                                                                                                                    |
| `missing-ok`                  | No       | `false`      | No     | `true`, `false`                                                         | Whether a missing path is considered `OK`. Incompatible with `exists-critical` or `exists-warning` options.                                                                                      |
| `fail-fast`                   | No       | `false`      | No     | `true`, `false`                                                         | Whether this plugin prioritizes speed of check results over always returning a `CRITICAL` state result before a `WARNING` state. This can be useful for processing large collections of content. |
| `age-critical`                | No       | `0`          | No     | `2+` (*minimum 1 greater than warning*)                                 | Assert that age for specified paths is less than or equal to the specified age in days, otherwise consider state to be `CRITICAL`.                                                               |
| `age-warning`                 | No       | `0`          | No     | `1+` (*minimum of 1*)                                                   | Assert that age for specified paths is less than or equal to the specified age in days, otherwise consider state to be `WARNING`.                                                                |
| `size-min-critical`           | No       | `0`          | No     | `1+` (*minimum of 1*)                                                   | Assert that size for specified paths is the specified size in bytes or greater, otherwise consider state to be `CRITICAL`.                                                                       |
| `size-min-warning`            | No       | `0`          | No     | `2+` (*minimum 1 larger than size-min-critical*)                        | Assert that size for specified paths is the specified size in bytes or greater, otherwise consider state to be `WARNING`.                                                                        |
| `size-max-critical`           | No       | `0`          | No     | `2+` (*minimum 1 greater than size-max-warning*)                        | Assert that size for specified paths is the specified size in bytes or less, otherwise consider state to be `CRITICAL`.                                                                          |
| `size-max-warning`            | No       | `0`          | No     | `1+` (*minimum of 1*)                                                   | Assert that size for specified paths is the specified size in bytes or less , otherwise consider state to be `WARNING`.                                                                          |
| `exists-critical`             | No       | `false`      | No     | `true`, `false`                                                         | Assert that specified paths are missing, otherwise consider state to be `CRITICAL`.                                                                                                              |
| `exists-warning`              | No       | `false`      | No     | `true`, `false`                                                         | Assert that specified paths are missing, otherwise consider state to be `WARNING`.                                                                                                               |
| `username-missing-critical`   | No       | `false`      | No     | *valid username*   (**not supported on Windows**)                       | Assert that specified owner/username is present on all content in specified paths, otherwise consider state to be `CRITICAL`.                                                                    |
| `username-missing-warning`    | No       | `false`      | No     | *valid username*   (**not supported on Windows**)                       | Assert that specified owner/username is present on all content in specified paths, otherwise consider state to be `WARNING`.                                                                     |
| `group-name-missing-critical` | No       | `false`      | No     | *valid group name* (**not supported on Windows**)                       | Assert that specified group name is present on all content in specified paths, otherwise consider state to be `CRITICAL`.                                                                        |
| `group-name-missing-warning`  | No       | `false`      | No     | *valid group name* (**not supported on Windows**)                       | Assert that specified group name is present on all content in specified paths, otherwise consider state to be `WARNING`.                                                                         |

### Environment Variables

If used, command-line arguments override the equivalent environment variables
listed below. See the [Command-line Arguments](#command-line-arguments) table
for more information.

| Flag Name                     | Environment Variable Name                | Notes | Example (mostly using default values)                        |
| ----------------------------- | ---------------------------------------- | ----- | ------------------------------------------------------------ |
| `emit-branding`               | `CHECK_PATH_EMIT_BRANDING`               |       | `CHECK_PATH_EMIT_BRANDING="false"`                           |
| `log-level`                   | `CHECK_PATH_LOG_LEVEL`                   |       | `CHECK_PATH_LOG_LEVEL="info"`                                |
| `paths`                       | `CHECK_PATH_PATHS_INCLUDE`               |       | `CHECK_PATH_PATHS_INCLUDE="/var/log/apache2 /var/log/samba"` |
| `ignore`                      | `CHECK_PATH_PATHS_IGNORE`                |       | `CHECK_PATH_PATHS_IGNORE="/var/log/apache2/access.log"`      |
| `recurse`                     | `CHECK_PATH_RECURSE`                     |       | `CHECK_PATH_RECURSE="false"`                                 |
| `missing-ok`                  | `CHECK_PATH_MISSING_OK`                  |       | `CHECK_PATH_MISSING_OK="false"`                              |
| `fail-fast`                   | `CHECK_PATH_FAIL_FAST`                   |       | `CHECK_PATH_FAIL_FAST="false"`                               |
| `age-critical`                | `CHECK_PATH_AGE_CRITICAL`                |       | `CHECK_PATH_AGE_CRITICAL="2"`                                |
| `age-warning`                 | `CHECK_PATH_AGE_WARNING`                 |       | `CHECK_PATH_AGE_WARNING="1"`                                 |
| `size-min-critical`           | `CHECK_PATH_SIZE_MIN_CRITICAL`           |       | `CHECK_PATH_SIZE_MIN_CRITICAL="2"`                           |
| `size-min-warning`            | `CHECK_PATH_SIZE_MIN_WARNING`            |       | `CHECK_PATH_SIZE_MIN_WARNING="1"`                            |
| `size-max-critical`           | `CHECK_PATH_SIZE_MAX_CRITICAL`           |       | `CHECK_PATH_SIZE_MAX_CRITICAL="2"`                           |
| `size-max-warning`            | `CHECK_PATH_SIZE_MAX_WARNING`            |       | `CHECK_PATH_SIZE_MAX_WARNING="1"`                            |
| `exists-critical`             | `CHECK_PATH_EXISTS_CRITICAL`             |       | `CHECK_PATH_EXISTS_CRITICAL="true"`                          |
| `exists-warning`              | `CHECK_PATH_EXISTS_WARNING`              |       | `CHECK_PATH_EXISTS_WARNING="true"`                           |
| `username-missing-critical`   | `CHECK_PATH_USERNAME_MISSING_CRITICAL`   |       | `CHECK_PATH_USERNAME_MISSING_CRITICAL="ubuntu"`              |
| `username-missing-warning`    | `CHECK_PATH_USERNAME_MISSING_WARNING`    |       | `CHECK_PATH_USERNAME_MISSING_WARNING="ubuntu"`               |
| `group-name-missing-critical` | `CHECK_PATH_GROUP_NAME_MISSING_CRITICAL` |       | `CHECK_PATH_GROUP_NAME_MISSING_CRITICAL="adm"`               |
| `group-name-missing-warning`  | `CHECK_PATH_GROUP_NAME_MISSING_WARNING`  |       | `CHECK_PATH_GROUP_NAME_MISSING_WARNING="adm"`                |

## Examples

### Help output

```ShellSession
$ ./release_assets/check_path/check_path-v0.1.1-28-gf18e040-linux-amd64 --help

Go-based tooling to check/verify filesystem paths as part of a Nagios service check
check-path v0.1.1-28-gf18e040 (https://github.com/atc0005/check-path)

Usage: check_path-v0.1.1-28-gf18e040-linux-amd64 [--log-level LOG-LEVEL] [--emit-branding] [--paths PATHS] [--ignore IGNORE] [--recurse] [--missing-ok] [--fail-fast] [--age-critical AGE-CRITICAL] [--age-warning AGE-WARNING] [--size-min-critical SIZE-MIN-CRITICAL] [--size-min-warning SIZE-MIN-WARNING] [--size-max-critical SIZE-MAX-CRITICAL] [--size-max-warning SIZE-MAX-WARNING] [--exists-critical] [--exists-warning] [--username-missing-critical USERNAME-MISSING-CRITICAL] [--username-missing-warning USERNAME-MISSING-WARNING] [--group-name-missing-critical GROUP-NAME-MISSING-CRITICAL] [--group-name-missing-warning GROUP-NAME-MISSING-WARNING]

Options:
  --log-level LOG-LEVEL
                         Maximum log level at which messages will be logged. Log messages below this threshold will be discarded.
  --emit-branding        Whether 'generated by' text is included at the bottom of application output. This output is included in the Nagios dashboard and notifications. This output may not mix well with branding output from other tools such as atc0005/send2teams which also insert their own branding output.
  --paths PATHS          List of comma or space-separated paths to check.
  --ignore IGNORE        List of comma or space-separated paths to ignore. Does not apply to existence checks.
  --recurse              Perform recursive search into subdirectories per provided path.
  --missing-ok           Whether a missing path is considered OK. Incompatible with exists-critical or exists-warning options.
  --fail-fast            Whether this plugin prioritizes speed of check results over always returning a CRITICAL state result before a WARNING state. This can be useful for processing large collections of content.
  --age-critical AGE-CRITICAL
                         Assert that age for specified paths is less than or equal to the specified age in days, otherwise consider state to be CRITICAL.
  --age-warning AGE-WARNING
                         Assert that age for specified paths is less than or equal to the specified age in days, otherwise consider state to be WARNING.
  --size-min-critical SIZE-MIN-CRITICAL
                         Assert that size for specified paths is the specified size in bytes or greater, otherwise consider state to be CRITICAL.
  --size-min-warning SIZE-MIN-WARNING
                         Assert that size for specified paths is the specified size in bytes or greater, otherwise consider state to be WARNING.
  --size-max-critical SIZE-MAX-CRITICAL
                         Assert that size for specified paths is the specified size in bytes or less, otherwise consider state to be CRITICAL.
  --size-max-warning SIZE-MAX-WARNING
                         Assert that size for specified paths is the specified size in bytes or less , otherwise consider state to be WARNING.
  --exists-critical      Assert that specified paths are missing, otherwise consider state to be CRITICAL.
  --exists-warning       Assert that specified paths are missing, otherwise consider state to be WARNING.
  --username-missing-critical USERNAME-MISSING-CRITICAL
                         Assert that specified owner/username is present on all content in specified paths, otherwise consider state to be CRITICAL.
  --username-missing-warning USERNAME-MISSING-WARNING
                         Assert that specified owner/username is present on all content in specified paths, otherwise consider state to be WARNING.
  --group-name-missing-critical GROUP-NAME-MISSING-CRITICAL
                         Assert that specified group name is present on all content in specified paths, otherwise consider state to be CRITICAL.
  --group-name-missing-warning GROUP-NAME-MISSING-WARNING
                         Assert that specified group name is present on all content in specified paths, otherwise consider state to be WARNING.
  --help, -h             display this help and exit
  --version              display version and exit
```

### Existence check

#### `CRITICAL`

```ShellSession
$ ./release_assets/check_path/check_path-v0.1.1-28-gf18e040-linux-amd64 --exists-critical --paths /tmp/go1.15.3.linux-amd64.tar.gz
CRITICAL: file "/tmp/go1.15.3.linux-amd64.tar.gz": path exists

**ERRORS**

* file "/tmp/go1.15.3.linux-amd64.tar.gz": path exists

**THRESHOLDS**

* CRITICAL: [Paths exist]
* WARNING: N/A

**DETAILED INFO**

* Paths to check: [/tmp/go1.15.3.linux-amd64.tar.gz]
* Paths to ignore: []
* Recursive search: false
* Fail-Fast: false
* Plugin: check-path v0.1.1-28-gf18e040 (https://github.com/atc0005/check-path)
* Path "/tmp/go1.15.3.linux-amd64.tar.gz"
** Last Modified: 2020-10-15 05:23:50.0738968 -0500 CDT
```

#### `WARNING`

```ShellSession
$ ./release_assets/check_path/check_path-v0.1.1-28-gf18e040-linux-amd64 --exists-warning --paths /tmp/go1.15.3.linux-amd64.tar.gz
WARNING: file "/tmp/go1.15.3.linux-amd64.tar.gz": path exists

**ERRORS**

* file "/tmp/go1.15.3.linux-amd64.tar.gz": path exists

**THRESHOLDS**

* CRITICAL: N/A
* WARNING: [Paths exist]

**DETAILED INFO**

* Paths to check: [/tmp/go1.15.3.linux-amd64.tar.gz]
* Paths to ignore: []
* Recursive search: false
* Fail-Fast: false
* Plugin: check-path v0.1.1-28-gf18e040 (https://github.com/atc0005/check-path)
* Path "/tmp/go1.15.3.linux-amd64.tar.gz"
** Last Modified: 2020-10-15 05:23:50.0738968 -0500 CDT
```

### Age check

#### `error examining path`

An example where `sudo` is needed to handle permission errors.

```ShellSession
$ ./release_assets/check_path/check_path-v0.1.1-28-gf18e040-linux-amd64 --age-warning 15 --age-critical 30 --paths /tmp/
{"level":"error","version":"v0.1.1-28-gf18e040","logging_level":"info","error":"error examining path \"/tmp\": open /tmp/tmp0dyy3wu9: permission denied","recursive":false,"path":"/tmp","caller":"github.com/atc0005/check-path/cmd/check_path/main.go:148","message":"error processing path"}
CRITICAL: Error processing path: /tmp

**ERRORS**

* error examining path "/tmp": open /tmp/tmp0dyy3wu9: permission denied

**THRESHOLDS**

* CRITICAL: [File age in days: 30]
* WARNING: [File age in days: 15]

**DETAILED INFO**

* Paths to check: [/tmp]
* Paths to ignore: []
* Recursive search: false
* Fail-Fast: false
* Plugin: check-path v0.1.1-28-gf18e040 (https://github.com/atc0005/check-path)
```

#### `CRITICAL`

```ShellSession
$ sudo ./release_assets/check_path/check_path-v0.1.1-28-gf18e040-linux-amd64 --age-warning 15 --age-critical 30 --paths /tmp/
{"level":"error","version":"v0.1.1-28-gf18e040","logging_level":"info","error":"old files found in path","critical_age_days":30,"warning_age_days":15,"age_check_enabled":true,"path":"/tmp","caller":"github.com/atc0005/check-path/cmd/check_path/age.go:53","message":"old files found"}
CRITICAL: file older than 30 days (65.06) found [path: "/tmp"]

**ERRORS**

* 29 files & directories evaluated: old files found in path

**THRESHOLDS**

* CRITICAL: [File age in days: 30]
* WARNING: [File age in days: 15]

**DETAILED INFO**

* Paths to check: [/tmp]
* Paths to ignore: []
* Recursive search: false
* Fail-Fast: false
* Plugin: check-path v0.1.1-28-gf18e040 (https://github.com/atc0005/check-path)
* File
** parent dir: "/tmp"
** name: "go1.15.2.linux-amd64.tar.gz"
** age: 65.06013975135417
```

#### `WARNING`

```ShellSession
$ sudo ./release_assets/check_path/check_path-v0.1.1-28-gf18e040-linux-amd64 --age-warning 30 --age-critical 90 --paths /tmp/
{"level":"error","version":"v0.1.1-28-gf18e040","logging_level":"info","error":"old files found in path","critical_age_days":90,"warning_age_days":30,"age_check_enabled":true,"path":"/tmp","caller":"github.com/atc0005/check-path/cmd/check_path/age.go:53","message":"old files found"}
WARNING: file older than 30 days (65.06) found [path: "/tmp"]

**ERRORS**

* 29 files & directories evaluated: old files found in path

**THRESHOLDS**

* CRITICAL: [File age in days: 90]
* WARNING: [File age in days: 30]

**DETAILED INFO**

* Paths to check: [/tmp]
* Paths to ignore: []
* Recursive search: false
* Fail-Fast: false
* Plugin: check-path v0.1.1-28-gf18e040 (https://github.com/atc0005/check-path)
* File
** parent dir: "/tmp"
** name: "go1.15.2.linux-amd64.tar.gz"
** age: 65.06036823248265
```

### Size check

#### Maximum file size

##### `CRITICAL`

Without recursion. This check looks only at the files in the immediate `/tmp`
directory, not subdirectories.

```ShellSession
$ sudo ./release_assets/check_path/check_path-v0.1.1-28-gf18e040-linux-amd64 --size-max-warning 100000000 --size-max-critical 121097000 --paths /tmp/
{"level":"error","version":"v0.1.1-28-gf18e040","logging_level":"info","error":"evaluated files in specified path too large (29 evaluated)","critical_size_max_bytes":121097000,"warning_size_max_bytes":100000000,"actual_size_bytes":594897382,"actual_size_hr":"567.3 MiB","size_max_check_enabled":true,"path":"/tmp","caller":"github.com/atc0005/check-path/cmd/check_path/size.go:90","message":"evaluated files in specified path too large (29 evaluated)"}
CRITICAL: maximum size threshold crossed; 567.3 MiB found in path "/tmp"

**ERRORS**

* evaluated files in specified path too large (29 evaluated)

**THRESHOLDS**

* CRITICAL: [Max File size (bytes: 121097000, Human: 115.5 MiB)]
* WARNING: [Max File size (bytes: 100000000, Human: 95.4 MiB)]

**DETAILED INFO**

* Paths to check: [/tmp]
* Paths to ignore: []
* Recursive search: false
* Fail-Fast: false
* Plugin: check-path v0.1.1-28-gf18e040 (https://github.com/atc0005/check-path)
* Size
** path: "/tmp"
** bytes: 594897382
** human-readable: 567.3 MiB
```

With recursion:

```ShellSession
$ sudo ./release_assets/check_path/check_path-v0.1.1-28-gf18e040-linux-amd64 --size-max-warning 100000000 --size-max-critical 121097000 --paths /tmp/ --recurse
{"level":"error","version":"v0.1.1-28-gf18e040","logging_level":"info","error":"evaluated files in specified path too large (4675 evaluated)","critical_size_max_bytes":121097000,"warning_size_max_bytes":100000000,"actual_size_bytes":854523940,"actual_size_hr":"814.9 MiB","size_max_check_enabled":true,"path":"/tmp","caller":"github.com/atc0005/check-path/cmd/check_path/size.go:90","message":"evaluated files in specified path too large (4675 evaluated)"}
CRITICAL: maximum size threshold crossed; 814.9 MiB found in path "/tmp"

**ERRORS**

* evaluated files in specified path too large (4675 evaluated)

**THRESHOLDS**

* CRITICAL: [Max File size (bytes: 121097000, Human: 115.5 MiB)]
* WARNING: [Max File size (bytes: 100000000, Human: 95.4 MiB)]

**DETAILED INFO**

* Paths to check: [/tmp]
* Paths to ignore: []
* Recursive search: true
* Fail-Fast: false
* Plugin: check-path v0.1.1-28-gf18e040 (https://github.com/atc0005/check-path)
* Size
** path: "/tmp"
** bytes: 854523940
** human-readable: 814.9 MiB
```

##### `WARNING`, enable `fail-fast` option

Combining the `fail-fast` option with `size-max-warning` and
`size-max-critical` triggers a state change upon finding any **single** file
that does not meet the specified thresholds.

See the [Known issues](#known-issues) section for more details.

```ShellSession
$ sudo ./release_assets/check_path/check_path-v0.1.1-28-gf18e040-linux-amd64 --size-max-warning 1000 --size-max-critical 121097000 --paths /tmp/ --fail-fast
{"level":"error","version":"v0.1.1-28-gf18e040","logging_level":"info","error":"evaluated files in specified path too large (4 evaluated)","critical_size_max_bytes":121097000,"warning_size_max_bytes":1000,"actual_size_bytes":3400457,"actual_size_hr":"3.2 MiB","size_max_check_enabled":true,"path":"/tmp","caller":"github.com/atc0005/check-path/cmd/check_path/size.go:90","message":"evaluated files in specified path too large (4 evaluated)"}
WARNING: maximum size threshold crossed; 3.2 MiB found in path "/tmp"

**ERRORS**

* evaluated files in specified path too large (4 evaluated)

**THRESHOLDS**

* CRITICAL: [Max File size (bytes: 121097000, Human: 115.5 MiB)]
* WARNING: [Max File size (bytes: 1000, Human: 1000 B)]

**DETAILED INFO**

* Paths to check: [/tmp]
* Paths to ignore: []
* Recursive search: false
* Fail-Fast: true
* Plugin: check-path v0.1.1-28-gf18e040 (https://github.com/atc0005/check-path)
* Size
** path: "/tmp"
** bytes: 3400457
** human-readable: 3.2 MiB
```

#### Minimum Size check

##### `OK`

Recursively assert that `/tmp` is greater than specified thresholds in bytes.

```ShellSession
$ sudo ./release_assets/check_path/check_path-v0.1.1-28-gf18e040-linux-amd64 --size-min-warning 150000000 --size-min-critical 121097000 --paths /tmp/ --recurse
{"level":"info","version":"v0.1.1-28-gf18e040","logging_level":"info","age_check_enabled":false,"size_min_check_enabled":true,"size_max_check_enabled":false,"caller":"github.com/atc0005/check-path/cmd/check_path/main.go:277","message":"1/1 specified paths pass min size validation checks (0 missing, 0 ignored by request)"}
OK: 1/1 specified paths pass min size validation checks (0 missing, 0 ignored by request)

**ERRORS**

* None

**THRESHOLDS**

* CRITICAL: [Min File size (bytes: 121097000, Human: 115.5 MiB)]
* WARNING: [Min File size (bytes: 150000000, Human: 143.1 MiB)]

**DETAILED INFO**

* Paths to check: [/tmp]
* Paths to ignore: []
* Recursive search: true
* Fail-Fast: false
* Plugin: check-path v0.1.1-28-gf18e040 (https://github.com/atc0005/check-path)
```

##### `CRITICAL`

Without recursion:

Here there is `594897382` bytes (directly within `/tmp`, not including
subdirectories), but we specified that only values greater than or equal to
`964523940` bytes were acceptable.

```ShellSession
$ sudo ./release_assets/check_path/check_path-v0.1.1-28-gf18e040-linux-amd64 --size-min-warning 964523940 --size-min-critical 944523940 --paths /tmp/
{"level":"error","version":"v0.1.1-28-gf18e040","logging_level":"info","error":"evaluated files in specified path too small (29 evaluated)","critical_size_min_bytes":944523940,"warning_size_min_bytes":964523940,"actual_size_bytes":594897382,"actual_size_hr":"567.3 MiB","size_min_check_enabled":true,"path":"/tmp","caller":"github.com/atc0005/check-path/cmd/check_path/size.go:125","message":"evaluated files in specified path too small (29 evaluated)"}
CRITICAL: minimum size threshold crossed; 567.3 MiB found in path "/tmp"

**ERRORS**

* evaluated files in specified path too small (29 evaluated)

**THRESHOLDS**

* CRITICAL: [Min File size (bytes: 944523940, Human: 900.8 MiB)]
* WARNING: [Min File size (bytes: 964523940, Human: 919.8 MiB)]

**DETAILED INFO**

* Paths to check: [/tmp]
* Paths to ignore: []
* Recursive search: false
* Fail-Fast: false
* Plugin: check-path v0.1.1-28-gf18e040 (https://github.com/atc0005/check-path)
* Size
** path: "/tmp"
** bytes: 594897382
** human-readable: 567.3 MiB
```

With recursion:

Here there is `854523940` bytes (within `/tmp`, including subdirectories), but
we specified that only values greater than or equal to `964523940` bytes were
acceptable.

```ShellSession
$ sudo ./release_assets/check_path/check_path-v0.1.1-28-gf18e040-linux-amd64 --size-min-warning 964523940 --size-min-critical 944523940 --paths /tmp/ --recurse
{"level":"error","version":"v0.1.1-28-gf18e040","logging_level":"info","error":"evaluated files in specified path too small (4675 evaluated)","critical_size_min_bytes":944523940,"warning_size_min_bytes":964523940,"actual_size_bytes":854523940,"actual_size_hr":"814.9 MiB","size_min_check_enabled":true,"path":"/tmp","caller":"github.com/atc0005/check-path/cmd/check_path/size.go:125","message":"evaluated files in specified path too small (4675 evaluated)"}
CRITICAL: minimum size threshold crossed; 814.9 MiB found in path "/tmp"

**ERRORS**

* evaluated files in specified path too small (4675 evaluated)

**THRESHOLDS**

* CRITICAL: [Min File size (bytes: 944523940, Human: 900.8 MiB)]
* WARNING: [Min File size (bytes: 964523940, Human: 919.8 MiB)]

**DETAILED INFO**

* Paths to check: [/tmp]
* Paths to ignore: []
* Recursive search: true
* Fail-Fast: false
* Plugin: check-path v0.1.1-28-gf18e040 (https://github.com/atc0005/check-path)
* Size
** path: "/tmp"
** bytes: 854523940
** human-readable: 814.9 MiB
```

##### `WARNING`

Without recursion:

```ShellSession
$ sudo ./release_assets/check_path/check_path-v0.1.1-28-gf18e040-linux-amd64 --size-min-warning 854523941 --size-min-critical 844523940 --paths /tmp/
{"level":"error","version":"v0.1.1-28-gf18e040","logging_level":"info","error":"evaluated files in specified path too small (29 evaluated)","critical_size_min_bytes":844523940,"warning_size_min_bytes":854523941,"actual_size_bytes":594897382,"actual_size_hr":"567.3 MiB","size_min_check_enabled":true,"path":"/tmp","caller":"github.com/atc0005/check-path/cmd/check_path/size.go:125","message":"evaluated files in specified path too small (29 evaluated)"}
CRITICAL: minimum size threshold crossed; 567.3 MiB found in path "/tmp"

**ERRORS**

* evaluated files in specified path too small (29 evaluated)

**THRESHOLDS**

* CRITICAL: [Min File size (bytes: 844523940, Human: 805.4 MiB)]
* WARNING: [Min File size (bytes: 854523941, Human: 814.9 MiB)]

**DETAILED INFO**

* Paths to check: [/tmp]
* Paths to ignore: []
* Recursive search: false
* Fail-Fast: false
* Plugin: check-path v0.1.1-28-gf18e040 (https://github.com/atc0005/check-path)
* Size
** path: "/tmp"
** bytes: 594897382
** human-readable: 567.3 MiB
```

With recursion:

```ShellSession
$ sudo ./release_assets/check_path/check_path-v0.1.1-28-gf18e040-linux-amd64 --size-min-warning 854523941 --size-min-critical 844523940 --paths /tmp/ --recurse
{"level":"error","version":"v0.1.1-28-gf18e040","logging_level":"info","error":"evaluated files in specified path too small (4675 evaluated)","critical_size_min_bytes":844523940,"warning_size_min_bytes":854523941,"actual_size_bytes":854523940,"actual_size_hr":"814.9 MiB","size_min_check_enabled":true,"path":"/tmp","caller":"github.com/atc0005/check-path/cmd/check_path/size.go:125","message":"evaluated files in specified path too small (4675 evaluated)"}
WARNING: minimum size threshold crossed; 814.9 MiB found in path "/tmp"

**ERRORS**

* evaluated files in specified path too small (4675 evaluated)

**THRESHOLDS**

* CRITICAL: [Min File size (bytes: 844523940, Human: 805.4 MiB)]
* WARNING: [Min File size (bytes: 854523941, Human: 814.9 MiB)]

**DETAILED INFO**

* Paths to check: [/tmp]
* Paths to ignore: []
* Recursive search: true
* Fail-Fast: false
* Plugin: check-path v0.1.1-28-gf18e040 (https://github.com/atc0005/check-path)
* Size
** path: "/tmp"
** bytes: 854523940
** human-readable: 814.9 MiB
```

##### `CRITICAL`, enable `fail-fast` behavior

Combining the `fail-fast` option with `size-min-warning` and
`size-min-critical` triggers a state change upon finding any **single** file
that does not meet the specified thresholds.

See the [Known issues](#known-issues) section for more details.

```ShellSession
$ sudo ./release_assets/check_path/check_path-v0.1.1-28-gf18e040-linux-amd64 --size-min-warning 900000000 --size-min-critical 700000000 --paths /tmp/ --recurse --fail-fast
{"level":"error","version":"v0.1.1-28-gf18e040","logging_level":"info","error":"evaluated files in specified path too small (1 evaluated)","critical_size_min_bytes":700000000,"warning_size_min_bytes":900000000,"actual_size_bytes":0,"actual_size_hr":"0 B","size_min_check_enabled":true,"path":"/tmp","caller":"github.com/atc0005/check-path/cmd/check_path/size.go:125","message":"evaluated files in specified path too small (1 evaluated)"}
CRITICAL: minimum size threshold crossed; 0 B found in path "/tmp"

**ERRORS**

* evaluated files in specified path too small (1 evaluated)

**THRESHOLDS**

* CRITICAL: [Min File size (bytes: 700000000, Human: 667.6 MiB)]
* WARNING: [Min File size (bytes: 900000000, Human: 858.3 MiB)]

**DETAILED INFO**

* Paths to check: [/tmp]
* Paths to ignore: []
* Recursive search: true
* Fail-Fast: true
* Plugin: check-path v0.1.1-28-gf18e040 (https://github.com/atc0005/check-path)
* Size
** path: "/tmp"
** bytes: 0
** human-readable: 0 B
```

### Username check

Our example file:

```ShellSession
$ ls -l /tmp/go1.15.3.linux-amd64.tar.gz
-rw-r--r-- 1 root ubuntu 121097663 Oct 15 05:23 /tmp/go1.15.3.linux-amd64.tar.gz
```

#### `CRITICAL`

```ShellSession
$ ./release_assets/check_path/check_path-v0.1.1-28-gf18e040-linux-amd64 --username-missing-critical ubuntu --paths /tmp/go1.15.3.linux-amd64.tar.gz
{"level":"error","version":"v0.1.1-28-gf18e040","logging_level":"info","error":"requested username not set on file/directory","username_check_enabled":true,"group_name_check_enabled":false,"path":"/tmp/go1.15.3.linux-amd64.tar.gz","caller":"github.com/atc0005/check-path/cmd/check_path/ids.go:71","message":"found username \"root\"; expected \"ubuntu\" [path: \"/tmp/go1.15.3.linux-amd64.tar.gz\"]"}
CRITICAL: found username "root"; expected "ubuntu" [path: "/tmp/go1.15.3.linux-amd64.tar.gz"]

**ERRORS**

* requested username not set on file/directory

**THRESHOLDS**

* Not specified

**DETAILED INFO**

* Paths to check: [/tmp/go1.15.3.linux-amd64.tar.gz]
* Paths to ignore: []
* Recursive search: false
* Fail-Fast: false
* Plugin: check-path v0.1.1-28-gf18e040 (https://github.com/atc0005/check-path)
```

#### `WARNING`

```ShellSession
$ ./release_assets/check_path/check_path-v0.1.1-28-gf18e040-linux-amd64 --username-missing-warning ubuntu --paths /tmp/go1.15.3.linux-amd64.tar.gz
{"level":"error","version":"v0.1.1-28-gf18e040","logging_level":"info","error":"requested username not set on file/directory","username_check_enabled":true,"group_name_check_enabled":false,"path":"/tmp/go1.15.3.linux-amd64.tar.gz","caller":"github.com/atc0005/check-path/cmd/check_path/ids.go:71","message":"found username \"root\"; expected \"ubuntu\" [path: \"/tmp/go1.15.3.linux-amd64.tar.gz\"]"}
WARNING: found username "root"; expected "ubuntu" [path: "/tmp/go1.15.3.linux-amd64.tar.gz"]

**ERRORS**

* requested username not set on file/directory

**THRESHOLDS**

* Not specified

**DETAILED INFO**

* Paths to check: [/tmp/go1.15.3.linux-amd64.tar.gz]
* Paths to ignore: []
* Recursive search: false
* Fail-Fast: false
* Plugin: check-path v0.1.1-28-gf18e040 (https://github.com/atc0005/check-path)
```

### Group Name check

Our example file:

```ShellSession
$ ls -l /tmp/go1.15.3.linux-amd64.tar.gz
-rw-r--r-- 1 root ubuntu 121097663 Oct 15 05:23 /tmp/go1.15.3.linux-amd64.tar.gz
```

#### `OK`

```ShellSession
$ ./release_assets/check_path/check_path-v0.1.1-28-gf18e040-linux-amd64 --group-name-missing-critical ubuntu --paths /tmp/go1.15.3.linux-amd64.tar.gz
{"level":"info","version":"v0.1.1-28-gf18e040","logging_level":"info","age_check_enabled":false,"size_min_check_enabled":false,"size_max_check_enabled":false,"caller":"github.com/atc0005/check-path/cmd/check_path/main.go:277","message":"1/1 specified paths pass group name validation checks (0 missing, 0 ignored by request)"}
OK: 1/1 specified paths pass group name validation checks (0 missing, 0 ignored by request)

**ERRORS**

* None

**THRESHOLDS**

* Not specified

**DETAILED INFO**

* Paths to check: [/tmp/go1.15.3.linux-amd64.tar.gz]
* Paths to ignore: []
* Recursive search: false
* Fail-Fast: false
* Plugin: check-path v0.1.1-28-gf18e040 (https://github.com/atc0005/check-path)
```

#### `CRITICAL`

This is a contrived example of checking for a file that we expect to find, but
doesn't exist.

```ShellSession
$ ./release_assets/check_path/check_path-v0.1.1-28-gf18e040-linux-amd64 --group-name-missing-warning ubuntu --paths /tmp/go1.15.1.linux-amd64.tar.gz
{"level":"error","version":"v0.1.1-28-gf18e040","logging_level":"info","error":"error examining path \"/tmp/go1.15.1.linux-amd64.tar.gz\": path does not exist: /tmp/go1.15.1.linux-amd64.tar.gz","recursive":false,"path":"/tmp/go1.15.1.linux-amd64.tar.gz","caller":"github.com/atc0005/check-path/cmd/check_path/main.go:148","message":"error processing path"}
CRITICAL: Error processing path: /tmp/go1.15.1.linux-amd64.tar.gz

**ERRORS**

* error examining path "/tmp/go1.15.1.linux-amd64.tar.gz": path does not exist: /tmp/go1.15.1.linux-amd64.tar.gz

**THRESHOLDS**

* Not specified

**DETAILED INFO**

* Paths to check: [/tmp/go1.15.1.linux-amd64.tar.gz]
* Paths to ignore: []
* Recursive search: false
* Fail-Fast: false
* Plugin: check-path v0.1.1-28-gf18e040 (https://github.com/atc0005/check-path)
```

#### `WARNING`

```ShellSession
$ ./release_assets/check_path/check_path-v0.1.1-28-gf18e040-linux-amd64 --group-name-missing-warning ubuntu --paths /tmp/go1.15.2.linux-amd64.tar.gz
{"level":"error","version":"v0.1.1-28-gf18e040","logging_level":"info","error":"requested group name not set on file/directory","username_check_enabled":false,"group_name_check_enabled":true,"path":"/tmp/go1.15.2.linux-amd64.tar.gz","caller":"github.com/atc0005/check-path/cmd/check_path/ids.go:112","message":"found group name \"root\"; expected \"ubuntu\" [path: \"/tmp/go1.15.2.linux-amd64.tar.gz\"]"}
WARNING: found group name "root"; expected "ubuntu" [path: "/tmp/go1.15.2.linux-amd64.tar.gz"]

**ERRORS**

* requested group name not set on file/directory

**THRESHOLDS**

* Not specified

**DETAILED INFO**

* Paths to check: [/tmp/go1.15.2.linux-amd64.tar.gz]
* Paths to ignore: []
* Recursive search: false
* Fail-Fast: false
* Plugin: check-path v0.1.1-28-gf18e040 (https://github.com/atc0005/check-path)
```

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
- <https://github.com/atc0005/check-path>
- <https://github.com/atc0005/check-illiad>

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

[go-supported-releases]: <https://go.dev/doc/devel/release#policy> "Go Release Policy"

<!-- []: PLACEHOLDER "DESCRIPTION_HERE" -->
