// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/check-path
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package main

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/atc0005/check-path/internal/config"
	"github.com/atc0005/check-path/internal/paths"
	"github.com/atc0005/go-nagios"
)

func main() {

	// Set initial "state" as valid, adjust as we go.
	var nagiosExitState = nagios.ExitState{
		LastError:      nil,
		ExitStatusCode: nagios.StateOKExitCode,
	}

	// defer this from the start so it is the last deferred function to run
	defer nagiosExitState.ReturnCheckResults()

	cfg, configErr := config.New()
	if configErr != nil {
		nagiosExitState.LastError = configErr
		nagiosExitState.ExitStatusCode = nagios.StateCRITICALExitCode
		log.Err(configErr).Msg("Error validating configuration")

		nagiosExitState.ServiceOutput = fmt.Sprintf(
			"%s: Failed to load configuration: %v",
			nagios.StateCRITICALLabel,
			configErr,
		)

		// no need to go any further, we *want* to exit right away; we don't
		// have a working configuration and there isn't anything further to do
		return
	}

	// If enabled, show application details at end of notification
	if cfg.EmitBranding() {
		nagiosExitState.BrandingCallback = config.Branding("Notification generated by ")
	}

	// contexts & cancellation are awesome
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// If sysadmin opted to ignore missing paths (e.g., Age or Size checks),
	// we track them in this slice for later reference in the summary output.
	missingOKPaths := make([]string, 0, 5)

	// If sysadmin opted to ignore specific paths, we track them in this slice
	// for later reference in the summary output.
	ignoredPaths := make([]string, 0, 5)

	// Resolve uid and gid values if sysadmin specified a username or
	// group name to compare against files in specified path
	resolveIDs := cfg.ResolveIDs()

	// Flesh out nagiosExitState with some additional common details now that
	// configuration flags have been parsed.
	nagiosExitState.LongServiceOutput = fmt.Sprintf(
		"* Paths to check: %v%s"+
			"* Paths to ignore: %v%s"+
			"* Recursive search: %v%s"+
			"* Fail-Fast: %v%s"+
			"* Plugin: %v%s",
		cfg.PathsInclude(),
		nagios.CheckOutputEOL,
		cfg.PathsExclude(),
		nagios.CheckOutputEOL,
		cfg.Recursive(),
		nagios.CheckOutputEOL,
		cfg.FailFast(),
		nagios.CheckOutputEOL,
		config.Version(),
		nagios.CheckOutputEOL,
	)

	setThresholdDescriptions(cfg, &nagiosExitState)

	// Check for existence of paths. NOTE: This check is not compatible with
	// other checks (e.g., Age, Size), so we exit ASAP after finishing.
	if cfg.PathExistsCritical() || cfg.PathExistsWarning() {
		checkExists(
			cfg.PathsInclude(),
			cfg.PathExistsCritical(),
			cfg.PathExistsWarning(),
			&cfg.Log,
			&nagiosExitState,
		)

		return
	}

	for _, path := range cfg.PathsInclude() {

		cfg.Log.Debug().Msgf("Processing path %s ...", path)

		// This placement is intentional. The channel is closed once
		// paths.Process returns, so we recreate the channel on the next
		// iteration of the paths list. If this is moved, a panic will likely
		// occur unless the current logic is reworked.
		results := make(chan paths.ProcessResult)

		// Process continues walking the path until complete, one of the
		// returned paths.MetaRecord values fails evaluation, or an error
		// occurs, whichever comes first.
		go paths.Process(ctx, path, cfg.PathsExclude(), cfg.Recursive(), results)

		// Collection of "records processed thus far" for the current path out
		// of the specified list that we're evaluating.
		var metaRecords paths.MetaRecords

		for result := range results {

			// fail early on errors from goroutine
			if result.Error != nil {

				// unless the error is an expected type
				switch {
				case errors.Is(result.Error, paths.ErrPathDoesNotExist):
					if cfg.MissingOK() {
						missingOKPaths = append(missingOKPaths, result.MetaRecord.FQPath)
						continue
					}

				case errors.Is(result.Error, paths.ErrPathIgnored):
					ignoredPaths = append(ignoredPaths, result.MetaRecord.FQPath)
					continue
				}

				cfg.Log.Error().Err(result.Error).
					Bool("recursive", cfg.Recursive()).
					Str("path", path).
					Msg("error processing path")

				nagiosExitState.LastError = result.Error
				nagiosExitState.ServiceOutput = fmt.Sprintf(
					"%s: Error processing path: %s",
					nagios.StateCRITICALLabel,
					path,
				)
				nagiosExitState.ExitStatusCode = nagios.StateCRITICALExitCode

				return
			}

			// no error thus far
			metaRecords = append(metaRecords, result.MetaRecord)

			if cfg.FailFast() {

				ageCheck := cfg.Age()
				if ageCheck.Set {
					ageCheckErr := checkAge(path, ageCheck, &cfg.Log, &nagiosExitState, result.MetaRecord)
					if ageCheckErr != nil {
						return
					}
				}

				sizeMaxCheck := cfg.SizeMax()
				sizeMinCheck := cfg.SizeMin()
				if sizeMaxCheck.Set || sizeMinCheck.Set {
					thsMinMax := config.FileSizeThresholdsMinMax{
						SizeMin: sizeMinCheck,
						SizeMax: sizeMaxCheck,
					}

					// evaluate the entire set of MetaRecord values each time
					// (instead of one at a time) in order to fail-fast when
					// the accumulated content size first crosses specified
					// size thresholds.
					sizeCheckErr := checkSize(path, thsMinMax, &cfg.Log, &nagiosExitState, metaRecords...)
					if sizeCheckErr != nil {
						return
					}

				}

				// if this is set, then sysadmin requested that we assert that
				// provided username or group name is present on all items
				// (including directories) in the specified paths.
				if resolveIDs.GroupNameCheck || resolveIDs.UsernameCheck {
					idsErr := checkIDs(path, resolveIDs, &cfg.Log, &nagiosExitState, result.MetaRecord)
					if idsErr != nil {
						return
					}
				}

			}

		}

		if !cfg.FailFast() {
			ageCheck := cfg.Age()
			if ageCheck.Set {
				ageCheckErr := checkAge(path, ageCheck, &cfg.Log, &nagiosExitState, metaRecords...)
				if ageCheckErr != nil {
					return
				}
			}

			sizeMaxCheck := cfg.SizeMax()
			sizeMinCheck := cfg.SizeMin()
			if sizeMaxCheck.Set || sizeMinCheck.Set {
				thsMinMax := config.FileSizeThresholdsMinMax{
					SizeMin: sizeMinCheck,
					SizeMax: sizeMaxCheck,
				}

				sizeCheckErr := checkSize(path, thsMinMax, &cfg.Log, &nagiosExitState, metaRecords...)
				if sizeCheckErr != nil {
					return
				}

			}

			if resolveIDs.GroupNameCheck || resolveIDs.UsernameCheck {
				idsErr := checkIDs(path, resolveIDs, &cfg.Log, &nagiosExitState, metaRecords...)
				if idsErr != nil {
					return
				}
			}
		}

	}

	// if we made it here, everything checked out
	otherChecksApplied := make([]string, 0, 2)

	if cfg.SizeMin().Set {
		otherChecksApplied = append(otherChecksApplied, "min size")
	}
	if cfg.SizeMax().Set {
		otherChecksApplied = append(otherChecksApplied, "max size")
	}
	if cfg.Age().Set {
		otherChecksApplied = append(otherChecksApplied, "age")
	}
	if resolveIDs.UsernameCheck {
		otherChecksApplied = append(otherChecksApplied, "username")
	}
	if resolveIDs.GroupNameCheck {
		otherChecksApplied = append(otherChecksApplied, "group name")
	}

	skippedEval := len(missingOKPaths)
	ignoredEval := len(ignoredPaths)
	okEval := len(cfg.PathsInclude()) - (skippedEval + ignoredEval)

	statusMsg := fmt.Sprintf(
		"%d/%d specified paths pass %v validation checks (%d missing, %d ignored by request)",
		okEval,
		len(cfg.PathsInclude()),
		strings.Join(otherChecksApplied, ", "),
		skippedEval,
		ignoredEval,
	)

	cfg.Log.Info().
		Bool("age_check_enabled", cfg.Age().Set).
		Bool("size_min_check_enabled", cfg.SizeMin().Set).
		Bool("size_max_check_enabled", cfg.SizeMax().Set).
		Msg(statusMsg)

	nagiosExitState.LastError = nil
	nagiosExitState.ServiceOutput = fmt.Sprintf(
		"%s: %s",
		nagios.StateOKLabel,
		statusMsg,
	)

	nagiosExitState.ExitStatusCode = nagios.StateOKExitCode

	// implied return, allow nagiosExitState.ReturnCheckResults() to run
	// return

}
