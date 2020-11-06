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
- [Known issues](#known-issues)
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
  - [Help output](#help-output)
  - [Existence check](#existence-check)
    - [`CRITICAL`](#critical)
    - [`WARNING`](#warning)
  - [Age check](#age-check)
    - [`error examining path`](#error-examining-path)
    - [`CRITICAL`](#critical-1)
    - [`WARNING`](#warning-1)
  - [Size check](#size-check)
    - [`CRITICAL`](#critical-2)
    - [`WARNING`, enable `fail-fast` option](#warning-enable-fail-fast-option)
  - [Username check](#username-check)
    - [`CRITICAL`](#critical-3)
    - [`WARNING`](#warning-2)
  - [Group Name check](#group-name-check)
    - [`OK`](#ok)
    - [`CRITICAL`](#critical-4)
    - [`WARNING`](#warning-3)
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
  - **NOTE**: this check is not supported on Windows
- Group Name checks
  - `CRITICAL` or `WARNING` (as specified) if missing
  - **NOTE**: this check is not supported on Windows
- Optional recursive evaluation toggle
- Optional "missing OK" toggle for all checks aside from the "existence"
  checks
- Optional "fail fast" behavior in an effort to avoid I/O churn over deep
  paths
  - see [Known issues](#known-issues) for potential issues with this option

## Known issues

If using the "early exit" behavior provided by the `fail-fast` flag, this
plugin will exit ASAP once a non-OK state is determined, regardless of whether
the first non-OK state is `CRITICAL` or `WARNING`. Receiving a `WARNING` state
for a path with files/directories which warrant a `CRITICAL` state result
could be confusing to troubleshoot, which is why `v0.1.1` changed the
default behavior to more closely match other Nagios plugins.

If you enable this option, just be aware of the tradeoff.

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

| Option                        | Required | Default      | Repeat | Possible                                                                | Description                                                                                                                                                                                      |
| ----------------------------- | -------- | ------------ | ------ | ----------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `h`, `help`                   | No       | `false`      | No     | `h`, `help`                                                             | Show Help text along with the list of supported flags.                                                                                                                                           |
| `emit-branding`               | No       | `false`      | No     | `true`, `false`                                                         | Toggles emission of branding details with plugin status details. This output is disabled by default.                                                                                             |
| `log-level`                   | No       | `info`       | No     | `disabled`, `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace` | Log message priority filter. Log messages with a lower level are ignored.                                                                                                                        |
| `paths`                       | Yes      | *empty list* | No     | *one or more valid files and directories*                               | List of comma or space-separated paths to process.                                                                                                                                               |
| `recursive`                   | No       | `false`      | No     | `true`, `false`                                                         | Perform recursive search into subdirectories.                                                                                                                                                    |
| `missing-ok`                  | No       | `false`      | No     | `true`, `false`                                                         | Whether a missing path is considered `OK`. Incompatible with `exists-critical` or `exists-warning` options.                                                                                      |
| `fail-fast`                   | No       | `false`      | No     | `true`, `false`                                                         | Whether this plugin prioritizes speed of check results over always returning a `CRITICAL` state result before a `WARNING` state. This can be useful for processing large collections of content. |
| `age-critical`                | No       | `0`          | No     | `2+` (*minimum 1 greater than warning*)                                 | Assert that age for specified paths is less than specified age in days, otherwise consider state to be `CRITICAL`.                                                                               |
| `age-warning`                 | No       | `0`          | No     | `1+` (*minimum of 1*)                                                   | Assert that age for specified paths is less than specified age in days, otherwise consider state to be `WARNING`.                                                                                |
| `size-critical`               | No       | `0`          | No     | `2+` (*minimum 1 greater than warning*)                                 | Assert that size for specified paths is less than specified size in bytes, otherwise consider state to be `CRITICAL`.                                                                            |
| `size-warning`                | No       | `0`          | No     | `1+` (*minimum of 1*)                                                   | Assert that size for specified paths is less than specified size in bytes, otherwise consider state to be `WARNING`.                                                                             |
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

| Flag Name                     | Environment Variable Name                | Notes | Example (mostly using default values)                     |
| ----------------------------- | ---------------------------------------- | ----- | --------------------------------------------------------- |
| `emit-branding`               | `CHECK_PATH_EMIT_BRANDING`               |       | `CHECK_PATH_EMIT_BRANDING="false"`                        |
| `log-level`                   | `CHECK_PATH_LOG_LEVEL`                   |       | `CHECK_PATH_LOG_LEVEL="info"`                             |
| `paths`                       | `CHECK_PATH_PATHS_LIST`                  |       | `CHECK_PATH_PATHS_LIST="/var/log/apache2 /var/log/samba"` |
| `recursive`                   | `CHECK_PATH_RECURSE`                     |       | `CHECK_PATH_RECURSE="false"`                              |
| `missing-ok`                  | `CHECK_PATH_MISSING_OK`                  |       | `CHECK_PATH_MISSING_OK="false"`                           |
| `fail-fast`                   | `CHECK_PATH_FAIL_FAST`                   |       | `CHECK_PATH_FAIL_FAST="false"`                            |
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

### Help output

```ShellSession
$ ./release_assets/check_path/check_path-v0.1.0-1-gdc622b8-linux-amd64 --help

Go-based tooling to check/verify filesystem paths as part of a Nagios service check
check-path v0.1.0-1-gdc622b8 (https://github.com/atc0005/check-path)

Usage: check_path-v0.1.0-1-gdc622b8-linux-amd64 [--log-level LOG-LEVEL] [--emit-branding] [--paths PATHS] [--recurse] [--missing-ok] [--fail-fast] [--age-critical AGE-CRITICAL] [--age-warning AGE-WARNING] [--size-critical SIZE-CRITICAL] [--size-warning SIZE-WARNING] [--exists-critical] [--exists-warning] [--username-missing-critical USERNAME-MISSING-CRITICAL] [--username-missing-warning USERNAME-MISSING-WARNING] [--group-name-missing-critical GROUP-NAME-MISSING-CRITICAL] [--group-name-missing-warning GROUP-NAME-MISSING-WARNING]

Options:
  --log-level LOG-LEVEL
                         Maximum log level at which messages will be logged. Log messages below this threshold will be discarded.
  --emit-branding        Whether 'generated by' text is included at the bottom of application output. This output is included in the Nagios dashboard and notifications. This output may not mix well with branding output from other tools such as atc0005/send2teams which also insert their own branding output.
  --paths PATHS          List of comma or space-separated paths to process.
  --recurse              Perform recursive search into subdirectories per provided path.
  --missing-ok           Whether a missing path is considered OK. Incompatible with exists-critical or exists-warning options.
  --fail-fast            Whether this plugin prioritizes speed of check results over always returning a CRITICAL state result before a WARNING state. This can be useful for processing large collections of content.
  --age-critical AGE-CRITICAL
                         Assert that age for specified paths is less than specified age in days, otherwise consider state to be CRITICAL.
  --age-warning AGE-WARNING
                         Assert that age for specified paths is less than specified age in days, otherwise consider state to be WARNING.
  --size-critical SIZE-CRITICAL
                         Assert that size for specified paths is less than specified size in bytes, otherwise consider state to be CRITICAL.
  --size-warning SIZE-WARNING
                         Assert that size for specified paths is less than specified size in bytes, otherwise consider state to be WARNING.
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
$ ./release_assets/check_path/check_path-v0.1.0-1-gdc622b8-linux-amd64 --exists-critical --paths /tmp/go1.15.3.linux-amd64.tar.gz
CRITICAL: file "/tmp/go1.15.3.linux-amd64.tar.gz": path exists

**ERRORS**

* file "/tmp/go1.15.3.linux-amd64.tar.gz": path exists

**THRESHOLDS**

* CRITICAL: [Paths exist]
* WARNING: N/A

**DETAILED INFO**

* Paths specified: [/tmp/go1.15.3.linux-amd64.tar.gz]
* Recursive search: false
* Fail-Fast: false
* Plugin: check-path v0.1.0-1-gdc622b8 (https://github.com/atc0005/check-path)
* Path "/tmp/go1.15.3.linux-amd64.tar.gz"
** Last Modified: 2020-10-15 05:23:50.0738968 -0500 CDT
```

#### `WARNING`

```ShellSession
$ ./release_assets/check_path/check_path-v0.1.0-1-gdc622b8-linux-amd64 --exists-warning --paths /tmp/go1.15.3.linux-amd64.tar.gz
WARNING: file "/tmp/go1.15.3.linux-amd64.tar.gz": path exists

**ERRORS**

* file "/tmp/go1.15.3.linux-amd64.tar.gz": path exists

**THRESHOLDS**

* CRITICAL: N/A
* WARNING: [Paths exist]

**DETAILED INFO**

* Paths specified: [/tmp/go1.15.3.linux-amd64.tar.gz]
* Recursive search: false
* Fail-Fast: false
* Plugin: check-path v0.1.0-1-gdc622b8 (https://github.com/atc0005/check-path)
* Path "/tmp/go1.15.3.linux-amd64.tar.gz"
** Last Modified: 2020-10-15 05:23:50.0738968 -0500 CDT
```

### Age check

#### `error examining path`

An example where `sudo` is needed to handle permission errors.

```ShellSession
$ ./release_assets/check_path/check_path-v0.1.0-1-gdc622b8-linux-amd64 --age-warning 15 --age-critical 30 --paths /tmp/
{"level":"error","version":"v0.1.0-1-gdc622b8","logging_level":"info","error":"error examining path \"/tmp/\": open /tmp/tmp0dyy3wu9: permission denied","recursive":false,"path":"/tmp/","caller":"github.com/atc0005/check-path/cmd/check_path/main.go:147","message":"error processing path"}
CRITICAL: Error processing path: /tmp/

**ERRORS**

* error examining path "/tmp/": open /tmp/tmp0dyy3wu9: permission denied

**THRESHOLDS**

* CRITICAL: [File age in days: 30]
* WARNING: [File age in days: 15]

**DETAILED INFO**

* Paths specified: [/tmp/]
* Recursive search: false
* Fail-Fast: false
* Plugin: check-path v0.1.0-1-gdc622b8 (https://github.com/atc0005/check-path)
```

#### `CRITICAL`

```ShellSession
$ sudo ./release_assets/check_path/check_path-v0.1.0-1-gdc622b8-linux-amd64 --age-warning 15 --age-critical 30 --paths /tmp/
{"level":"error","version":"v0.1.0-1-gdc622b8","logging_level":"info","error":"old files found in path","critical_age_days":30,"warning_age_days":15,"age_check_enabled":true,"path":"/tmp/","caller":"github.com/atc0005/check-path/cmd/check_path/main.go:181","message":"old files found"}
CRITICAL: file older than 30 days (56.06) found [path: "/tmp/"]

**ERRORS**

* old files found in path

**THRESHOLDS**

* CRITICAL: [File age in days: 30]
* WARNING: [File age in days: 15]

**DETAILED INFO**

* Paths specified: [/tmp/]
* Recursive search: false
* Fail-Fast: false
* Plugin: check-path v0.1.0-1-gdc622b8 (https://github.com/atc0005/check-path)
* File
** parent dir: "/tmp"
** name: "go1.15.2.linux-amd64.tar.gz"
** age: 56.05608141632987
```

#### `WARNING`

```ShellSession
$ sudo ./release_assets/check_path/check_path-v0.1.0-1-gdc622b8-linux-amd64 --age-warning 30 --age-critical 60 --paths /tmp/
{"level":"error","version":"v0.1.0-1-gdc622b8","logging_level":"info","error":"old files found in path","critical_age_days":60,"warning_age_days":30,"age_check_enabled":true,"path":"/tmp/","caller":"github.com/atc0005/check-path/cmd/check_path/main.go:181","message":"old files found"}
WARNING: file older than 30 days (56.06) found [path: "/tmp/"]

**ERRORS**

* old files found in path

**THRESHOLDS**

* CRITICAL: [File age in days: 60]
* WARNING: [File age in days: 30]

**DETAILED INFO**

* Paths specified: [/tmp/]
* Recursive search: false
* Fail-Fast: false
* Plugin: check-path v0.1.0-1-gdc622b8 (https://github.com/atc0005/check-path)
* File
** parent dir: "/tmp"
** name: "go1.15.2.linux-amd64.tar.gz"
** age: 56.05668234489005
```

### Size check

Our example file:

```ShellSession
$ ls -l /tmp/go1.15.3.linux-amd64.tar.gz
-rw-r--r-- 1 root ubuntu 121097663 Oct 15 05:23 /tmp/go1.15.3.linux-amd64.tar.gz
```

#### `CRITICAL`

```ShellSession
$ sudo ./release_assets/check_path/check_path-v0.1.0-1-gdc622b8-linux-amd64 --size-warning 100000000 --size-critical 121097000 --paths /tmp/
{"level":"error","version":"v0.1.0-1-gdc622b8","logging_level":"info","error":"evaluated files in specified path too large","critical_size_bytes":121097000,"warning_size_bytes":100000000,"actual_size_bytes":121149509,"actual_size_hr":"115.5 MiB","size_check_enabled":true,"path":"/tmp/","caller":"github.com/atc0005/check-path/cmd/check_path/main.go:259","message":"evaluated files too large"}
CRITICAL: size threshold crossed; 115.5 MiB found in path "/tmp/"

**ERRORS**

* evaluated files in specified path too large

**THRESHOLDS**

* CRITICAL: [File size (bytes: 121097000, Human: 115.5 MiB)]
* WARNING: [File size (bytes: 100000000, Human: 95.4 MiB)]

**DETAILED INFO**

* Paths specified: [/tmp/]
* Recursive search: false
* Fail-Fast: false
* Plugin: check-path v0.1.0-1-gdc622b8 (https://github.com/atc0005/check-path)
* Size
** path: "/tmp/"
** bytes: 121149509
** human-readable: 115.5 MiB
```

#### `WARNING`, enable `fail-fast` option

Mostly a repeat of the `CRITICAL` state `Size` example, but here we reduce the
`WARNING` threshold even further to just 1K bytes and enable the `fail-fast`
logic to illustrate the indeterminate nature of the setting.

```ShellSession
$ sudo ./release_assets/check_path/check_path-v0.1.0-1-gdc622b8-linux-amd64 --size-warning 1000 --size-critical 121097000 --paths /tmp/ --fail-fast
{"level":"error","version":"v0.1.0-1-gdc622b8","logging_level":"info","error":"evaluated files in specified path too large (4 thus far)","critical_size_bytes":121097000,"warning_size_bytes":1000,"actual_size_bytes":3400457,"actual_size_hr":"3.2 MiB","size_check_enabled":true,"path":"/tmp/","caller":"github.com/atc0005/check-path/cmd/check_path/main.go:259","message":"evaluated files too large"}
WARNING: size threshold crossed; 3.2 MiB found in path "/tmp/"

**ERRORS**

* evaluated files in specified path too large (4 thus far)

**THRESHOLDS**

* CRITICAL: [File size (bytes: 121097000, Human: 115.5 MiB)]
* WARNING: [File size (bytes: 1000, Human: 1000 B)]

**DETAILED INFO**

* Paths specified: [/tmp/]
* Recursive search: false
* Fail-Fast: true
* Plugin: check-path v0.1.0-1-gdc622b8 (https://github.com/atc0005/check-path)
* Size
** path: "/tmp/"
** bytes: 3400457
** human-readable: 3.2 MiB
```

### Username check

Our example file:

```ShellSession
$ ls -l /tmp/go1.15.3.linux-amd64.tar.gz
-rw-r--r-- 1 root ubuntu 121097663 Oct 15 05:23 /tmp/go1.15.3.linux-amd64.tar.gz
```

#### `CRITICAL`

```ShellSession
$ ./release_assets/check_path/check_path-v0.1.0-1-gdc622b8-linux-amd64 --username-missing-critical ubuntu --paths /tmp/go1.15.3.linux-amd64.tar.gz
{"level":"error","version":"v0.1.0-1-gdc622b8","logging_level":"info","error":"requested username not set on file/directory","username_check_enabled":true,"group_name_check_enabled":false,"path":"/tmp/go1.15.3.linux-amd64.tar.gz","caller":"github.com/atc0005/check-path/cmd/check_path/main.go:346","message":"found username \"root\"; expected \"ubuntu\" [path: \"/tmp/go1.15.3.linux-amd64.tar.gz\"]"}
CRITICAL: found username "root"; expected "ubuntu" [path: "/tmp/go1.15.3.linux-amd64.tar.gz"]

**ERRORS**

* requested username not set on file/directory

**THRESHOLDS**

* Not specified

**DETAILED INFO**

* Paths specified: [/tmp/go1.15.3.linux-amd64.tar.gz]
* Recursive search: false
* Fail-Fast: false
* Plugin: check-path v0.1.0-1-gdc622b8 (https://github.com/atc0005/check-path)
```

#### `WARNING`

```ShellSession
$ ./release_assets/check_path/check_path-v0.1.0-1-gdc622b8-linux-amd64 --username-missing-warning ubuntu --paths /tmp/go1.15.3.linux-amd64.tar.gz
{"level":"error","version":"v0.1.0-1-gdc622b8","logging_level":"info","error":"requested username not set on file/directory","username_check_enabled":true,"group_name_check_enabled":false,"path":"/tmp/go1.15.3.linux-amd64.tar.gz","caller":"github.com/atc0005/check-path/cmd/check_path/main.go:346","message":"found username \"root\"; expected \"ubuntu\" [path: \"/tmp/go1.15.3.linux-amd64.tar.gz\"]"}
WARNING: found username "root"; expected "ubuntu" [path: "/tmp/go1.15.3.linux-amd64.tar.gz"]

**ERRORS**

* requested username not set on file/directory

**THRESHOLDS**

* Not specified

**DETAILED INFO**

* Paths specified: [/tmp/go1.15.3.linux-amd64.tar.gz]
* Recursive search: false
* Fail-Fast: false
* Plugin: check-path v0.1.0-1-gdc622b8 (https://github.com/atc0005/check-path)
```

### Group Name check

Our example file:

```ShellSession
$ ls -l /tmp/go1.15.3.linux-amd64.tar.gz
-rw-r--r-- 1 root ubuntu 121097663 Oct 15 05:23 /tmp/go1.15.3.linux-amd64.tar.gz
```

#### `OK`

```ShellSession
$ ./release_assets/check_path/check_path-v0.1.0-1-gdc622b8-linux-amd64 --group-name-missing-critical ubuntu --paths /tmp/go1.15.3.linux-amd64.tar.gz
{"level":"info","version":"v0.1.0-1-gdc622b8","logging_level":"info","age_check_enabled":false,"size_check_enabled":false,"caller":"github.com/atc0005/check-path/cmd/check_path/main.go:451","message":"1/1 specified paths pass group name validation checks (0 missing & ignored by request)"}
OK: 1/1 specified paths pass group name validation checks (0 missing & ignored by request)

**ERRORS**

* None

**THRESHOLDS**

* Not specified

**DETAILED INFO**

* Paths specified: [/tmp/go1.15.3.linux-amd64.tar.gz]
* Recursive search: false
* Fail-Fast: false
* Plugin: check-path v0.1.0-1-gdc622b8 (https://github.com/atc0005/check-path)
```

#### `CRITICAL`

This is a contrived example of checking for a file that we expect to find, but
doesn't exist.

```ShellSession
$ ./release_assets/check_path/check_path-v0.1.0-1-gdc622b8-linux-amd64 --group-name-missing-warning ubuntu --paths /tmp/go1.15.1.linux-amd64.tar.gz
{"level":"error","version":"v0.1.0-1-gdc622b8","logging_level":"info","error":"error examining path \"/tmp/go1.15.1.linux-amd64.tar.gz\": path does not exist: /tmp/go1.15.1.linux-amd64.tar.gz","recursive":false,"path":"/tmp/go1.15.1.linux-amd64.tar.gz","caller":"github.com/atc0005/check-path/cmd/check_path/main.go:147","message":"error processing path"}
CRITICAL: Error processing path: /tmp/go1.15.1.linux-amd64.tar.gz

**ERRORS**

* error examining path "/tmp/go1.15.1.linux-amd64.tar.gz": path does not exist: /tmp/go1.15.1.linux-amd64.tar.gz

**THRESHOLDS**

* Not specified

**DETAILED INFO**

* Paths specified: [/tmp/go1.15.1.linux-amd64.tar.gz]
* Recursive search: false
* Fail-Fast: false
* Plugin: check-path v0.1.0-1-gdc622b8 (https://github.com/atc0005/check-path
```

#### `WARNING`

```ShellSession
$ ./release_assets/check_path/check_path-v0.1.0-1-gdc622b8-linux-amd64 --group-name-missing-warning ubuntu --paths /tmp/go1.15.2.linux-amd64.tar.gz
{"level":"error","version":"v0.1.0-1-gdc622b8","logging_level":"info","error":"requested group name not set on file/directory","username_check_enabled":false,"group_name_check_enabled":true,"path":"/tmp/go1.15.2.linux-amd64.tar.gz","caller":"github.com/atc0005/check-path/cmd/check_path/main.go:388","message":"found group name \"root\"; expected \"ubuntu\" [path: \"/tmp/go1.15.2.linux-amd64.tar.gz\"]"}
WARNING: found group name "root"; expected "ubuntu" [path: "/tmp/go1.15.2.linux-amd64.tar.gz"]

**ERRORS**

* requested group name not set on file/directory

**THRESHOLDS**

* Not specified

**DETAILED INFO**

* Paths specified: [/tmp/go1.15.2.linux-amd64.tar.gz]
* Recursive search: false
* Fail-Fast: false
* Plugin: check-path v0.1.0-1-gdc622b8 (https://github.com/atc0005/check-path)
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
- <https://github.com/atc0005/check-cert>
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

<!-- []: PLACEHOLDER "DESCRIPTION_HERE" -->
