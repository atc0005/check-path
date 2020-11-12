// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/check-path
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/alexflint/go-arg"
	"github.com/atc0005/check-path/internal/caller"
)

// see `constants.go`, `logging.go` for other related values

// version reflects the application version. This is overridden via Makefile
// for release builds.
var version string = "dev build"

// Version emits application name, version and repo location.
func Version() string {
	return fmt.Sprintf("%s %s (%s)", myAppName, version, myAppURL)
}

// Branding accepts a message and returns a function that concatenates that
// message with version information. This returned function implements
// nagios.ExitCallBackFunc in order to optionally supply branding details with
// the service check output.
func Branding(msg string) func() string {
	return func() string {
		return strings.Join([]string{msg, Version()}, "")
	}
}

// String implements the Stringer interface in order to display all
// initialized (user-provided or default) values.
func (c Config) String() string {
	return fmt.Sprintf(
		"{ PathsInclude: %v, "+
			"PathsExclude: %v, "+
			"LogLevel: %v, "+
			"Recursive: %v, "+
			"MissingOK: %v, "+
			"EmitBranding: %v, "+
			"Age: [Critical: %v, Warning: %v, Set: %v], "+
			"SizeMin: [Critical: %v, Warning: %v, Set: %v], "+
			"SizeMax: [Critical: %v, Warning: %v, Set: %v], "+
			"PathExists: [Critical: %v, Warning: %v], "+
			"User: [Name: %q, Critical: %v, Warning: %v], "+
			"Group: [Name: %q, Critical: %v, Warning: %v] }",
		c.PathsInclude(),
		c.PathsExclude(),
		c.LogLevel(),
		c.Recursive(),
		c.MissingOK(),
		c.EmitBranding(),
		c.Age().Critical,
		c.Age().Warning,
		c.Age().Set,
		c.SizeMin().Critical,
		c.SizeMin().Warning,
		c.SizeMin().Set,
		c.SizeMax().Critical,
		c.SizeMax().Warning,
		c.SizeMax().Set,
		c.PathExistsCritical(),
		c.PathExistsWarning(),
		c.Username(),
		c.UsernameCritical(),
		c.UsernameWarning(),
		c.GroupName(),
		c.GroupNameCritical(),
		c.GroupNameWarning(),
	)
}

// Version reuses the package-level Version function to emit version
// information and associated branding details whenever the user specifies the
// `--version` flag. The application exits after displaying this information.
func (c Config) Version() string {
	return Version() + "\n"
}

// Description emits branding information whenever the user specifies the `-h`
// flag. The application uses this as a header prior to displaying available
// CLI flag options.
func (c Config) Description() string {
	return fmt.Sprintf("\n%s", myAppDescription)
}

// New is a factory function that produces a new Config object based on user
// provided flag and where applicable, default values.
func New() (*Config, error) {

	myFuncName := caller.GetFuncName()

	var config Config

	// Bundle the returned `*.arg.Parser` for potential later use.
	config.flagParser = arg.MustParse(&config)

	if err := config.validate(); err != nil {
		// As of Nagios 3.x, stderr is not processed, so this is visible to
		// the user running the plugin from CLI only.
		// config.flagParser.WriteHelp(os.Stderr)
		// config.flagParser.Fail(err.Error())
		config.flagParser.WriteUsage(os.Stderr)

		return nil, fmt.Errorf(
			"%s: failed to validate configuration: %w",
			myFuncName,
			err,
		)
	}

	config.configureLogging()

	return &config, nil

}
