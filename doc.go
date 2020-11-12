// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/check-path
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

/*

This repo provides a Nagios plugin used to verify the ownership, group,
permissions, age or size of specific files or directories.

PROJECT HOME

See our GitHub repo (https://github.com/atc0005/check-path) for the latest
code, to file an issue or submit improvements for review and potential
inclusion into the project.

PURPOSE

Verify the ownership, group, age, permissions, size or existence of specific
files or directories.

FEATURES

• Age checks: CRITICAL and WARNING thresholds

• Size checks: minimum and maximum CRITICAL and WARNING thresholds

• Existence checks: CRITICAL or WARNING (as specified) if present

• Username checks: CRITICAL or WARNING (as specified) if missing

• Group Name checks: CRITICAL or WARNING (as specified) if missing

• Optional directory recursion

• Optional "missing OK" toggle, compatible with most checks

• Optional exclusion of specific paths from evaluation

• Optional "fail fast" behavior in an effort to avoid I/O churn over deep paths. See "Known issues" section of README for potential issues with this option.

USAGE

See our main README for supported settings and examples.

*/
package main
