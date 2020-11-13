// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/check-path
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/atc0005/go-nagios"
)

// ErrUserNameHasSpaces is returned by validation checks if a username
// contains spaces.
var ErrUserNameHasSpaces = errors.New("username contains spaces")

// ErrGroupNameHasSpaces is returned by validation checks if a group name
// contains spaces.
var ErrGroupNameHasSpaces = errors.New("group name contains spaces")

// ErrUsernameIsEmpty is returned by validation checks if a username is an
// empty string.
var ErrUsernameIsEmpty = errors.New("username is empty string")

// ErrGroupNameIsEmpty is returned by validation checks if a group name is an
// empty string.
var ErrGroupNameIsEmpty = errors.New("group name is empty string")

// usernameValidation is intended to help concentrate validation checks
// specific to usernames in one place.
func usernameValidation(username string) error {

	// TODO: Extend validation to cover most common username, group name
	// checks
	// https://unix.stackexchange.com/questions/157426/what-is-the-regex-to-validate-linux-users
	// https://github.com/systemd/systemd/issues/6237

	if username == "" {
		return ErrUsernameIsEmpty
	}

	if strings.Contains(username, " ") {
		return ErrUserNameHasSpaces
	}

	// TODO: Extend with further checks

	return nil
}

// groupNameValidation is intended to help concentrate validation checks
// specific to group names in one place.
func groupNameValidation(groupName string) error {

	// TODO: Extend validation to cover most common username, group name
	// checks
	// https://unix.stackexchange.com/questions/157426/what-is-the-regex-to-validate-linux-users
	// https://github.com/systemd/systemd/issues/6237

	if groupName == "" {
		return ErrGroupNameIsEmpty
	}

	if strings.Contains(groupName, " ") {
		return ErrGroupNameHasSpaces
	}

	// TODO: Extend with further checks

	return nil
}

// pathSizeValidation is used as a helper validation function for size checks
// to reduce code duplication.
func pathSizeValidation(ths FileSizeThresholds, sizeCritical *int64, sizeWarning *int64) error {

	const (
		tmplNotSetErrMsg                     string = "%s size in bytes not specified for %s threshold; both values required if checking %s file size"
		tmplTooSmallErrMsg                   string = "provided %s size in bytes (%d) not valid for %s threshold"
		tmplWarningGreaterThanCriticalErrMsg string = "provided %s %s size in bytes (%d) greater than %s %s size in bytes (%d)"
		tmplWarningLessThanCriticalErrMsg    string = "provided %s %s size in bytes (%d) less than %s %s size in bytes (%d)"
		tmplWarningEqualToCriticalErrMsg     string = "provided %s %s size in bytes (%d) equal to %s %s size in bytes (%d)"
	)

	if sizeCritical == nil {
		return fmt.Errorf(
			tmplNotSetErrMsg,
			ths.Description,
			nagios.StateCRITICALLabel,
			ths.Description,
		)
	}

	if sizeWarning == nil {
		return fmt.Errorf(
			tmplNotSetErrMsg,
			ths.Description,
			nagios.StateWARNINGLabel,
			ths.Description,
		)
	}

	if *sizeCritical <= 0 {
		return fmt.Errorf(
			tmplTooSmallErrMsg,
			ths.Description,
			*sizeCritical,
			nagios.StateCRITICALLabel,
		)
	}

	if *sizeWarning <= 0 {
		return fmt.Errorf(
			tmplTooSmallErrMsg,
			ths.Description,
			*sizeWarning,
			nagios.StateWARNINGLabel,
		)
	}

	switch {
	case ths.Description == sizeMaxDescription:
		if *sizeWarning > *sizeCritical {
			return fmt.Errorf(
				tmplWarningGreaterThanCriticalErrMsg,
				ths.Description,
				nagios.StateWARNINGLabel,
				*sizeWarning,
				ths.Description,
				nagios.StateCRITICALLabel,
				*sizeCritical,
			)
		}
	case ths.Description == sizeMinDescription:
		if *sizeWarning < *sizeCritical {
			return fmt.Errorf(
				tmplWarningLessThanCriticalErrMsg,
				ths.Description,
				nagios.StateWARNINGLabel,
				*sizeWarning,
				ths.Description,
				nagios.StateCRITICALLabel,
				*sizeCritical,
			)
		}
	}

	if *sizeWarning == *sizeCritical {
		return fmt.Errorf(
			tmplWarningEqualToCriticalErrMsg,
			ths.Description,
			nagios.StateWARNINGLabel,
			*sizeWarning,
			ths.Description,
			nagios.StateCRITICALLabel,
			*sizeCritical,
		)
	}

	return nil

}

// validate verifies that user-provided and/or default values are acceptable.
//
// Where possible/reliable, getter methods are checked instead of directly
// referencing the config struct. This approach is used because the getter
// methods pass user-provided values through without modification; if a user
// did not specify a value, the default value is passed through for
// validation.
func (c Config) validate() error {

	if c.Search.PathsInclude == nil {
		return fmt.Errorf("one or more paths not provided")
	}

	// TODO: Search.PathsExclude - how to handle this one? The file or
	// directory not existing should not be treated as a problem.

	switch c.LogLevel() {
	case LogLevelDisabled:
	case LogLevelPanic:
	case LogLevelFatal:
	case LogLevelError:
	case LogLevelWarn:
	case LogLevelInfo:
	case LogLevelDebug:
	case LogLevelTrace:
	default:
		return fmt.Errorf("invalid log level provided: %v", c.LogLevel())
	}

	// Search.Recursive is optional and boolean
	// Search.MissingOK is optional and boolean
	// Logging.EmitBranding is optional and boolean

	// One pair of AgeCritical, AgeWarning or SizeCritical and SizeWarning has
	// to be specified. If one of AgeCritical or AgeWarning is specified, the
	// other has to be specified also.

	existsCriticalSet := c.Search.ExistsCritical != nil
	existsWarningSet := c.Search.ExistsWarning != nil

	ageCriticalSet := c.Search.AgeCritical != nil
	ageWarningSet := c.Search.AgeWarning != nil

	sizeMaxCriticalSet := c.Search.SizeMaxCritical != nil
	sizeMaxWarningSet := c.Search.SizeMaxWarning != nil
	sizeMaxSet := c.Search.SizeMaxCritical != nil && c.Search.SizeMaxWarning != nil

	sizeMinCriticalSet := c.Search.SizeMinCritical != nil
	sizeMinWarningSet := c.Search.SizeMinWarning != nil
	sizeMinSet := c.Search.SizeMinCritical != nil && c.Search.SizeMinWarning != nil

	usernameMissingCriticalSet := c.Search.UsernameMissingCritical != nil
	usernameMissingWarningSet := c.Search.UsernameMissingWarning != nil

	groupNameMissingCriticalSet := c.Search.GroupNameMissingCritical != nil
	groupNameMissingWarningSet := c.Search.GroupNameMissingWarning != nil

	// Needs to be maintained to list all potential conflicts.
	// TODO: What is a better way to handle this?
	if (existsCriticalSet || existsWarningSet) &&
		(sizeMaxCriticalSet ||
			sizeMaxWarningSet ||
			sizeMinCriticalSet ||
			sizeMinWarningSet ||
			ageCriticalSet ||
			ageWarningSet ||
			usernameMissingCriticalSet ||
			usernameMissingWarningSet ||
			groupNameMissingCriticalSet ||
			groupNameMissingWarningSet) {

		if existsCriticalSet {
			return fmt.Errorf(
				"'exists-critical' incompatible with other options",
			)
		}
		if existsWarningSet {
			return fmt.Errorf(
				"'exists-warning' incompatible with other options",
			)
		}
	}

	if existsCriticalSet && existsWarningSet {
		return fmt.Errorf(
			"'exists-critical' and 'exists-warning' specified; " +
				"only one is permitted",
		)
	}

	if ageCriticalSet || ageWarningSet {

		notSetErrMsg :=
			"minimum age in days not specified for %s threshold; " +
				"both values required if checking file age"

		tooSmallErrMsg :=
			"provided age in days (%d) not valid for %s threshold"

		warningGreaterThanCriticalMsg :=
			"provided %s age in days (%d) greater than %s age in days (%d)"

		warningEqualToCriticalMsg :=
			"provided %s age in days (%d) equal to %s age in days (%d)"

		if !ageCriticalSet {
			return fmt.Errorf(notSetErrMsg, nagios.StateCRITICALLabel)
		}

		if !ageWarningSet {
			return fmt.Errorf(notSetErrMsg, nagios.StateWARNINGLabel)
		}

		ageThresholds := c.Age()

		if ageThresholds.Critical <= 0 {
			return fmt.Errorf(
				tooSmallErrMsg,
				ageThresholds.Critical,
				nagios.StateCRITICALLabel,
			)
		}

		if ageThresholds.Warning <= 0 {
			return fmt.Errorf(
				tooSmallErrMsg,
				ageThresholds.Warning,
				nagios.StateWARNINGLabel,
			)
		}

		if ageThresholds.Warning > ageThresholds.Critical {
			return fmt.Errorf(
				warningGreaterThanCriticalMsg,
				nagios.StateWARNINGLabel,
				ageThresholds.Warning,
				nagios.StateCRITICALLabel,
				ageThresholds.Critical,
			)
		}

		if ageThresholds.Warning == ageThresholds.Critical {
			return fmt.Errorf(
				warningEqualToCriticalMsg,
				nagios.StateWARNINGLabel,
				ageThresholds.Warning,
				nagios.StateCRITICALLabel,
				ageThresholds.Critical,
			)
		}
	}

	if sizeMaxCriticalSet || sizeMaxWarningSet {
		sizeErr := pathSizeValidation(
			c.SizeMax(),
			c.Search.SizeMaxCritical,
			c.Search.SizeMaxWarning,
		)
		if sizeErr != nil {
			return sizeErr
		}
	}

	if sizeMinCriticalSet || sizeMinWarningSet {
		sizeErr := pathSizeValidation(
			c.SizeMin(),
			c.Search.SizeMinCritical,
			c.Search.SizeMinWarning,
		)
		if sizeErr != nil {
			return sizeErr
		}
	}

	if usernameMissingCriticalSet && usernameMissingWarningSet {
		return fmt.Errorf(
			"username-missing-critical' and " +
				"'username-missing-warning' specified; " +
				"only one is permitted",
		)
	}

	if usernameMissingCriticalSet || usernameMissingWarningSet {

		if usernameMissingCriticalSet {
			if osWindows {
				return fmt.Errorf(
					"username-missing-critical' specified; " +
						"not currently supported for Windows",
				)
			}
			if err := usernameValidation(*c.Search.UsernameMissingCritical); err != nil {
				return fmt.Errorf(
					"invalid value %q specified for username-missing-critical: %w",
					*c.Search.UsernameMissingCritical,
					err,
				)
			}
		}

		if usernameMissingWarningSet {
			if osWindows {
				return fmt.Errorf(
					"username-missing-warning' specified; " +
						"not currently supported for Windows",
				)
			}
			if err := usernameValidation(*c.Search.UsernameMissingWarning); err != nil {
				return fmt.Errorf(
					"invalid value %q specified for username-missing-warning: %w",
					*c.Search.UsernameMissingWarning,
					err,
				)
			}
		}
	}

	if groupNameMissingCriticalSet && groupNameMissingWarningSet {
		return fmt.Errorf(
			"group-name-missing-critical' and " +
				"'group-name-missing-warning' specified; " +
				"only one is permitted",
		)
	}

	if groupNameMissingCriticalSet || groupNameMissingWarningSet {

		if groupNameMissingCriticalSet {
			if osWindows {
				return fmt.Errorf(
					"group-name-missing-critical' specified; " +
						"not currently supported for Windows",
				)
			}
			if err := groupNameValidation(*c.Search.GroupNameMissingCritical); err != nil {
				return fmt.Errorf(
					"invalid value %q specified for group-name-missing-critical: %w",
					*c.Search.GroupNameMissingCritical,
					err,
				)
			}
		}

		if groupNameMissingWarningSet {
			if osWindows {
				return fmt.Errorf(
					"group-name-missing-warning' specified; " +
						"not currently supported for Windows",
				)
			}
			if err := groupNameValidation(*c.Search.GroupNameMissingWarning); err != nil {
				return fmt.Errorf(
					"invalid value %q specified for group-name-missing-warning: %w",
					*c.Search.GroupNameMissingWarning,
					err,
				)
			}
		}

	}

	// if neither size (both), age (both), existence (only one), username
	// (only one) or group name (only one) is provided, then configuration is
	// incomplete
	if !(sizeMinSet || sizeMaxSet) &&
		!(ageCriticalSet && ageWarningSet) &&
		!(existsCriticalSet || existsWarningSet) &&
		!(usernameMissingCriticalSet || usernameMissingWarningSet) &&
		!(groupNameMissingCriticalSet || groupNameMissingWarningSet) {
		return fmt.Errorf(
			"no values specified for age, minimum size, maximum size, username, group name or existence",
		)
	}

	// declare that all options are valid if we make it this far
	return nil

}
