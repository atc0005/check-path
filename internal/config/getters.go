// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/check-path
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package config

import "path/filepath"

// PathsInclude returns the user-provided list of paths to check or an empty
// list if a user-specified list of paths was not provided. Each path in the
// list is processed by filepath.Clean. This allows for more reliable
// comparison against PathsExclude values to determine if a path should be
// ignored.
func (c Config) PathsInclude() []string {
	switch {
	case c.Search.PathsInclude != nil:
		cleanedPaths := make([]string, len(c.Search.PathsInclude))
		for i, path := range c.Search.PathsInclude {
			cleanedPaths[i] = filepath.Clean(path)
		}

		return cleanedPaths
	default:
		// this will probably not be reached due to validation checks ensuring
		// that a value was provided
		return []string{}
	}
}

// PathsExclude returns the user-provided list of paths to ignore or an empty
// list if a user-specified list of paths was not provided. Each path in the
// list is processed by filepath.Clean. This allows for more reliable
// comparison against PathsInclude values to determine if a path should be
// ignored.
func (c Config) PathsExclude() []string {
	switch {
	case c.Search.PathsExclude != nil:
		cleanedPaths := make([]string, len(c.Search.PathsExclude))
		for i, path := range c.Search.PathsExclude {
			cleanedPaths[i] = filepath.Clean(path)
		}

		return cleanedPaths

	default:
		return []string{}
	}
}

// LogLevel returns the user-provided logging level or the default value if
// not provided.
func (c Config) LogLevel() string {

	switch {
	case c.Logging.Level != nil:
		return *c.Logging.Level
	default:
		return defaultLogLevel
	}
}

// Recursive returns the user-provided choice of whether paths are checked
// recursively or the default value if not provided.
func (c Config) Recursive() bool {
	switch {
	case c.Search.Recursive != nil:
		return *c.Search.Recursive
	default:
		return defaultSearchRecursive
	}
}

// MissingOK returns the user-provided choice of whether missing paths are
// considered OK or the default value if not provided.
func (c Config) MissingOK() bool {
	switch {
	case c.Search.MissingOK != nil:
		return *c.Search.MissingOK
	default:
		return defaultSearchMissingOK
	}
}

// FailFast returns the user-provided choice of whether paths are processed in
// a way that prioritizes a first-fail result over a strict order of CRITICAL
// results before WARNING results. The default value is returned if not
// provided.
func (c Config) FailFast() bool {
	switch {
	case c.Search.FailFast != nil:
		return *c.Search.FailFast
	default:
		return defaultSearchFailFast
	}
}

// EmitBranding returns the user-provided choice of whether branded output is
// emitted with check results or the default value if not provided.
func (c Config) EmitBranding() bool {
	switch {
	case c.Logging.EmitBranding != nil:
		return *c.Logging.EmitBranding
	default:
		return defaultEmitBranding
	}
}

// Age returns the user-provided CRITICAL and WARNING thresholds in days for
// the specified paths.
func (c Config) Age() FileAgeThresholds {
	switch {
	case c.Search.AgeCritical != nil && c.Search.AgeWarning != nil:
		return FileAgeThresholds{
			Critical: *c.Search.AgeCritical,
			Warning:  *c.Search.AgeWarning,
			Set:      true,
		}
	default:
		return FileAgeThresholds{
			Set: false,
		}
	}
}

// SizeMin returns the user-provided CRITICAL and WARNING thresholds for
// minimum size in bytes for the specified paths.
func (c Config) SizeMin() FileSizeThresholds {
	switch {
	case c.Search.SizeMinCritical != nil && c.Search.SizeMinWarning != nil:
		return FileSizeThresholds{
			Description: sizeMinDescription,
			Critical:    *c.Search.SizeMinCritical,
			Warning:     *c.Search.SizeMinWarning,
			Set:         true,
		}
	default:
		return FileSizeThresholds{
			Description: sizeMinDescription,
			Set:         false,
		}
	}
}

// SizeMax returns the user-provided CRITICAL and WARNING thresholds for
// maximum size in bytes for the specified paths.
func (c Config) SizeMax() FileSizeThresholds {
	switch {
	case c.Search.SizeMaxCritical != nil && c.Search.SizeMaxWarning != nil:
		return FileSizeThresholds{
			Description: sizeMaxDescription,
			Critical:    *c.Search.SizeMaxCritical,
			Warning:     *c.Search.SizeMaxWarning,
			Set:         true,
		}
	default:
		return FileSizeThresholds{
			Description: sizeMaxDescription,
			Set:         false,
		}
	}
}

// PathExistsCritical indicates whether the existence of specified paths is
// considered a CRITICAL state.
func (c Config) PathExistsCritical() bool {
	return c.Search.ExistsCritical != nil && *c.Search.ExistsCritical
}

// PathExistsWarning indicates whether the existence of specified paths is
// considered a WARNING state.
func (c Config) PathExistsWarning() bool {
	return c.Search.ExistsWarning != nil && *c.Search.ExistsWarning
}

// Username returns the user-provided username set via the
// username-missing-critical or username-missing-warning flags or the
// default value if not provided.
func (c Config) Username() string {
	switch {
	case c.Search.UsernameMissingCritical != nil:
		return *c.Search.UsernameMissingCritical
	case c.Search.UsernameMissingWarning != nil:
		return *c.Search.UsernameMissingWarning
	default:
		return defaultUsername
	}
}

// UsernameCritical indicates whether user opted to check for username
// mismatches. Failing results indicate a CRITICAL state.
func (c Config) UsernameCritical() bool {
	return c.Search.UsernameMissingCritical != nil
}

// UsernameWarning indicates whether user opted to check for group name
// mismatches. Failing results indicate a WARNING state.
func (c Config) UsernameWarning() bool {
	return c.Search.UsernameMissingWarning != nil
}

// GroupName returns the user-provided group name set via the
// group-name-missing-critical or group-name-missing-warning flags or the
// default value if not provided.
func (c Config) GroupName() string {
	switch {
	case c.Search.GroupNameMissingCritical != nil:
		return *c.Search.GroupNameMissingCritical
	case c.Search.GroupNameMissingWarning != nil:
		return *c.Search.GroupNameMissingWarning
	default:
		return defaultGroupName
	}
}

// GroupNameCritical indicates whether user opted to check for group name
// mismatches. Failing results indicate a CRITICAL state.
func (c Config) GroupNameCritical() bool {
	return c.Search.GroupNameMissingCritical != nil
}

// GroupNameWarning indicates whether user opted to check for group name
// mismatches. Failing results indicate a WARNING state.
func (c Config) GroupNameWarning() bool {
	return c.Search.GroupNameMissingWarning != nil
}

// ResolveIDs returns a ResolveIDs type which indicates whether user opted to
// resolve user and group id values to name values and if so, at which exit
// state values.
func (c Config) ResolveIDs() ResolveIDs {
	return ResolveIDs{
		IDs: IDs{
			Username:  c.Username(),
			GroupName: c.GroupName(),
		},
		UsernameCheck:     c.UsernameCritical() || c.UsernameWarning(),
		UsernameCritical:  c.UsernameCritical(),
		UsernameWarning:   c.UsernameWarning(),
		GroupNameCheck:    c.GroupNameCritical() || c.GroupNameWarning(),
		GroupNameCritical: c.GroupNameCritical(),
		GroupNameWarning:  c.GroupNameWarning(),
	}
}
