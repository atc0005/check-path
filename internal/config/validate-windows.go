//go:build windows
// +build windows

// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/check-path
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package config

// osWindows is a small alternative to importing the runtime package in order
// to determine if the current OS is Windows (or not). Build tags handle
// setting this value appropriately.
const osWindows bool = true
