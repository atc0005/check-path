// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/check-path
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package config

import (
	"github.com/alexflint/go-arg"
	"github.com/rs/zerolog"
)

// FileAgeThresholds represents the user-specified file age thresholds for
// specified paths.
type FileAgeThresholds struct {
	Critical int
	Warning  int
	Set      bool
}

// FileSizeThresholds represents the user-specified file size thresholds for
// specified paths.
type FileSizeThresholds struct {
	Description string
	Critical    int64
	Warning     int64
	Set         bool
}

// FileSizeThresholdsMinMax represents the combined minimum and maximum
// user-specified file size thresholds for specified paths.
type FileSizeThresholdsMinMax struct {
	SizeMin FileSizeThresholds
	SizeMax FileSizeThresholds
}

// ResolveIDs is a helper struct to record whether user opted to resolve user
// and group id values to name values and if so, at which exit state values.
type ResolveIDs struct {
	UsernameCheck     bool
	UsernameCritical  bool
	UsernameWarning   bool
	GroupNameCheck    bool
	GroupNameCritical bool
	GroupNameWarning  bool
	IDs
}

// IDs represents username and group name values.
type IDs struct {
	Username  string
	GroupName string
}

// Search represents options specific to controlling how this application
// performs searches in the filesystem.
type Search struct {
	PathsInclude             []string `arg:"--paths,env:CHECK_PATH_PATHS_INCLUDE" help:"List of comma or space-separated paths to check."`
	PathsExclude             []string `arg:"--ignore,env:CHECK_PATH_PATHS_IGNORE" help:"List of comma or space-separated paths to ignore. Does not apply to existence checks."`
	Recursive                *bool    `arg:"--recurse,env:CHECK_PATH_RECURSE" help:"Perform recursive search into subdirectories per provided path."`
	MissingOK                *bool    `arg:"--missing-ok,env:CHECK_PATH_MISSING_OK" help:"Whether a missing path is considered OK. Incompatible with exists-critical or exists-warning options."`
	FailFast                 *bool    `arg:"--fail-fast,env:CHECK_PATH_FAIL_FAST" help:"Whether this plugin prioritizes speed of check results over always returning a CRITICAL state result before a WARNING state. This can be useful for processing large collections of content."`
	AgeCritical              *int     `arg:"--age-critical,env:CHECK_PATH_AGE_CRITICAL" help:"Assert that age for specified paths is less than or equal to the specified age in days, otherwise consider state to be CRITICAL."`
	AgeWarning               *int     `arg:"--age-warning,env:CHECK_PATH_AGE_WARNING" help:"Assert that age for specified paths is less than or equal to the specified age in days, otherwise consider state to be WARNING."`
	SizeMinCritical          *int64   `arg:"--size-min-critical,env:CHECK_PATH_SIZE_MIN_CRITICAL" help:"Assert that size for specified paths is the specified size in bytes or greater, otherwise consider state to be CRITICAL."`
	SizeMinWarning           *int64   `arg:"--size-min-warning,env:CHECK_PATH_SIZE_MIN_WARNING" help:"Assert that size for specified paths is the specified size in bytes or greater, otherwise consider state to be WARNING."`
	SizeMaxCritical          *int64   `arg:"--size-max-critical,env:CHECK_PATH_SIZE_MAX_CRITICAL" help:"Assert that size for specified paths is the specified size in bytes or less, otherwise consider state to be CRITICAL."`
	SizeMaxWarning           *int64   `arg:"--size-max-warning,env:CHECK_PATH_SIZE_MAX_WARNING" help:"Assert that size for specified paths is the specified size in bytes or less , otherwise consider state to be WARNING."`
	ExistsCritical           *bool    `arg:"--exists-critical,env:CHECK_PATH_EXISTS_CRITICAL" help:"Assert that specified paths are missing, otherwise consider state to be CRITICAL."`
	ExistsWarning            *bool    `arg:"--exists-warning,env:CHECK_PATH_EXISTS_WARNING" help:"Assert that specified paths are missing, otherwise consider state to be WARNING."`
	UsernameMissingCritical  *string  `arg:"--username-missing-critical,env:CHECK_PATH_USERNAME_MISSING_CRITICAL" help:"Assert that specified owner/username is present on all content in specified paths, otherwise consider state to be CRITICAL."`
	UsernameMissingWarning   *string  `arg:"--username-missing-warning,env:CHECK_PATH_USERNAME_MISSING_WARNING" help:"Assert that specified owner/username is present on all content in specified paths, otherwise consider state to be WARNING."`
	GroupNameMissingCritical *string  `arg:"--group-name-missing-critical,env:CHECK_PATH_GROUP_NAME_MISSING_CRITICAL" help:"Assert that specified group name is present on all content in specified paths, otherwise consider state to be CRITICAL."`
	GroupNameMissingWarning  *string  `arg:"--group-name-missing-warning,env:CHECK_PATH_GROUP_NAME_MISSING_WARNING" help:"Assert that specified group name is present on all content in specified paths, otherwise consider state to be WARNING."`
}

// Logging represents options specific to how this application handles
// logging.
type Logging struct {
	Level        *string `arg:"--log-level,env:CHECK_PATH_LOG_LEVEL" help:"Maximum log level at which messages will be logged. Log messages below this threshold will be discarded."`
	EmitBranding *bool   `arg:"--emit-branding,env:CHECK_PATH_EMIT_BRANDING" help:"Whether 'generated by' text is included at the bottom of application output. This output is included in the Nagios dashboard and notifications. This output may not mix well with branding output from other tools such as atc0005/send2teams which also insert their own branding output."`
}

// Config is a unified set of configuration values for this application. This
// struct is configured via command-line flags or (maybe in the future) TOML
// configuration file provided by the user. Values held by this object are
// intended to be retrieved via "getter" methods.
type Config struct {
	Logging
	Search

	// Log is an embedded zerolog Logger initialized via config.New().
	Log zerolog.Logger `arg:"-"`

	flagParser *arg.Parser `arg:"-"`
}
