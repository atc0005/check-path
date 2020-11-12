// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/check-path
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package config

const (

	// MyAppName is the public name of this application.
	myAppName string = "check-path"

	// MyAppURL is the location of the repo for this application.
	myAppURL string = "https://github.com/atc0005/check-path"

	// MyAppDescription is the description for this application shown in
	// HelpText output.
	myAppDescription string = "Go-based tooling to check/verify filesystem paths as part of a Nagios service check"
)

// Default (flag, config file, etc) settings if not overridden by user input.
const (
	defaultLogLevel        string = "info"
	defaultSearchRecursive bool   = false
	defaultSearchMissingOK bool   = false
	defaultSearchFailFast  bool   = false
	defaultEmitBranding    bool   = false

	// these values have to be supplied via flag by the sysadmin to be useful
	defaultUsername  string = ""
	defaultGroupName string = ""
)

// used by SizeMin and SizeMax getter methods for threshold descriptions
const (
	sizeMinDescription string = "minimum"
	sizeMaxDescription string = "maximum"
)
