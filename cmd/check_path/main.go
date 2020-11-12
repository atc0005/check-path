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
	"time"

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
	resolveUsername := cfg.UsernameCritical() || cfg.UsernameWarning()
	resolveGroupName := cfg.GroupNameCritical() || cfg.GroupNameWarning()
	resolveIDs := resolveUsername || resolveGroupName

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
		checkPaths(
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
		go paths.Process(ctx, path, cfg.PathsExclude(), cfg.Recursive(), cfg.FailFast(), results)

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

				// shutdown goroutine
				cancel()

				return
			}

			// no error thus far
			metaRecords = append(metaRecords, result.MetaRecord)

			ageCheck := cfg.Age()
			if ageCheck.Set && !result.MetaRecord.IsDir() {

				criticalAgeFile := paths.AgeExceeded(
					result.MetaRecord.FileInfo, ageCheck.Critical)

				warningAgeFile := paths.AgeExceeded(
					result.MetaRecord.FileInfo, ageCheck.Warning)

				if criticalAgeFile || warningAgeFile {
					cfg.Log.Error().Err(paths.ErrPathOldFilesFound).
						Int("critical_age_days", ageCheck.Critical).
						Int("warning_age_days", ageCheck.Warning).
						Bool("age_check_enabled", ageCheck.Set).
						Str("path", path).
						Msg("old files found")

					nagiosExitState.LastError = paths.ErrPathOldFilesFound
					if cfg.FailFast() {
						nagiosExitState.LastError = fmt.Errorf(
							"%d files & directories evaluated thus far: %w",
							len(metaRecords),
							paths.ErrPathOldFilesFound,
						)
					}

					fileAge := time.Since(result.MetaRecord.ModTime()).Hours() / 24

					nagiosExitState.LongServiceOutput += fmt.Sprintf(
						"* File %s** parent dir: %q%s** name: %q%s** age: %v%s",
						nagios.CheckOutputEOL,
						result.MetaRecord.ParentDir,
						nagios.CheckOutputEOL,
						result.MetaRecord.Name(),
						nagios.CheckOutputEOL,
						fileAge,
						nagios.CheckOutputEOL,
					)

					switch {
					case criticalAgeFile:
						nagiosExitState.ServiceOutput = fmt.Sprintf(
							"%s: file older than %d days (%.2f) found [path: %q]",
							nagios.StateCRITICALLabel,
							ageCheck.Critical,
							fileAge,
							path,
						)

						nagiosExitState.ExitStatusCode = nagios.StateCRITICALExitCode

						return

					case warningAgeFile:
						nagiosExitState.ServiceOutput = fmt.Sprintf(
							"%s: file older than %d days (%.2f) found [path: %q]",
							nagios.StateWARNINGLabel,
							ageCheck.Warning,
							fileAge,
							path,
						)

						nagiosExitState.ExitStatusCode = nagios.StateWARNINGExitCode

						return
					}

				}
			}

			sizeMaxCheck := cfg.SizeMax()
			sizeMinCheck := cfg.SizeMin()
			if (sizeMaxCheck.Set || sizeMinCheck.Set) && !result.MetaRecord.IsDir() {
				actualSizeHR := metaRecords.TotalFileSizeHR()
				actualSizeBytes := metaRecords.TotalFileSize()
				sizeOfFilesTooLargeErr := errors.New("evaluated files in specified path too large")
				sizeOfFilesTooSmallErr := errors.New("evaluated files in specified path too small")

				// warning threshold required, so we can use that to reduce
				// conditional check logic complexity
				if actualSizeBytes < sizeMinCheck.Warning || actualSizeBytes > sizeMaxCheck.Warning {

					if cfg.FailFast() {

						sizeOfFilesTooLargeErr = fmt.Errorf(
							"%s (%d thus far)",
							sizeOfFilesTooLargeErr.Error(),
							len(metaRecords),
						)

						sizeOfFilesTooSmallErr = fmt.Errorf(
							"%s (%d thus far)",
							sizeOfFilesTooSmallErr.Error(),
							len(metaRecords),
						)

					}

					serviceOutputTmpl := fmt.Sprintf(
						"size threshold crossed; %v found in path %q",
						actualSizeHR,
						path,
					)

					nagiosExitState.LongServiceOutput += fmt.Sprintf(
						"* Size %s** path: %q%s** bytes: %v%s** human-readable: %v%s",
						nagios.CheckOutputEOL,
						path,
						nagios.CheckOutputEOL,
						actualSizeBytes,
						nagios.CheckOutputEOL,
						actualSizeHR,
						nagios.CheckOutputEOL,
					)

					// configure exit state details based on how the
					// thresholds were crossed. return after all exit state
					// details are recorded
					switch {

					case actualSizeBytes > sizeMaxCheck.Critical || actualSizeBytes > sizeMaxCheck.Warning:
						cfg.Log.Error().Err(sizeOfFilesTooLargeErr).
							Int64("critical_size_max_bytes", sizeMaxCheck.Critical).
							Int64("warning_size_max_bytes", sizeMaxCheck.Warning).
							Int64("actual_size_bytes", actualSizeBytes).
							Str("actual_size_hr", actualSizeHR).
							Bool("size_max_check_enabled", sizeMaxCheck.Set).
							Str("path", path).
							Msg(sizeOfFilesTooLargeErr.Error())

						nagiosExitState.LastError = sizeOfFilesTooLargeErr

						var stateLabel string
						var exitCode int

						if actualSizeBytes > sizeMaxCheck.Critical {
							stateLabel = nagios.StateCRITICALLabel
							exitCode = nagios.StateCRITICALExitCode
						}

						if actualSizeBytes > sizeMaxCheck.Warning {
							stateLabel = nagios.StateWARNINGLabel
							exitCode = nagios.StateWARNINGExitCode
						}

						nagiosExitState.ServiceOutput = fmt.Sprintf(
							"%s: %s %s",
							stateLabel,
							sizeMaxCheck.Description,
							serviceOutputTmpl,
						)

						nagiosExitState.ExitStatusCode = exitCode

					case actualSizeBytes < sizeMinCheck.Critical || actualSizeBytes < sizeMinCheck.Warning:
						cfg.Log.Error().Err(sizeOfFilesTooSmallErr).
							Int64("critical_size_min_bytes", sizeMinCheck.Critical).
							Int64("warning_size_min_bytes", sizeMinCheck.Warning).
							Int64("actual_size_bytes", actualSizeBytes).
							Str("actual_size_hr", actualSizeHR).
							Bool("size_min_check_enabled", sizeMinCheck.Set).
							Str("path", path).
							Msg(sizeOfFilesTooSmallErr.Error())

						nagiosExitState.LastError = sizeOfFilesTooSmallErr

						var stateLabel string
						var exitCode int

						if actualSizeBytes < sizeMinCheck.Critical {
							stateLabel = nagios.StateCRITICALLabel
							exitCode = nagios.StateCRITICALExitCode
						}

						if actualSizeBytes < sizeMinCheck.Warning {
							stateLabel = nagios.StateWARNINGLabel
							exitCode = nagios.StateWARNINGExitCode
						}

						nagiosExitState.ServiceOutput = fmt.Sprintf(
							"%s: %s %s",
							stateLabel,
							sizeMinCheck.Description,
							serviceOutputTmpl,
						)

						nagiosExitState.ExitStatusCode = exitCode

					}

					// Size check tripped, we're done here
					return

				}

			}

			// if this is set, then sysadmin requested that we assert that
			// provided username or group name is present on all items
			// (including directories) in the specified paths.
			if resolveIDs {

				cfg.Log.Debug().Msg("Username, Group name resolution enabled")

				resolveErr := paths.ResolveIDs(&result.MetaRecord)
				if resolveErr != nil {
					cfg.Log.Error().Err(resolveErr).
						Str("path", path).
						Msg(resolveErr.Error())

					nagiosExitState.LastError = resolveErr
					nagiosExitState.ServiceOutput = fmt.Sprintf(
						"%s: failed to resolve IDs: %v [path: %q]",
						nagios.StateCRITICALLabel,
						resolveErr.Error(),
						path,
					)

					return
				}

				if resolveUsername &&
					result.MetaRecord.Username != cfg.Username() {

					err := fmt.Errorf("requested username not set on file/directory")
					errMsg := fmt.Errorf(
						"found username %q; expected %q [path: %q]",
						result.MetaRecord.Username,
						cfg.Username(),
						result.MetaRecord.FQPath,
					)

					cfg.Log.Error().Err(err).
						Bool("username_check_enabled", resolveUsername).
						Bool("group_name_check_enabled", resolveGroupName).
						Str("path", path).
						Msg(errMsg.Error())

					nagiosExitState.LastError = err

					switch {
					case cfg.UsernameCritical():
						nagiosExitState.ServiceOutput = fmt.Sprintf(
							"%s: %s",
							nagios.StateCRITICALLabel,
							errMsg.Error(),
						)
						nagiosExitState.ExitStatusCode = nagios.StateCRITICALExitCode

						return

					case cfg.UsernameWarning():
						nagiosExitState.ServiceOutput = fmt.Sprintf(
							"%s: %s",
							nagios.StateWARNINGLabel,
							errMsg.Error(),
						)
						nagiosExitState.ExitStatusCode = nagios.StateWARNINGExitCode

						return
					}
				}

				if resolveGroupName &&
					result.MetaRecord.GroupName != cfg.GroupName() {

					err := fmt.Errorf("requested group name not set on file/directory")
					errMsg := fmt.Errorf(
						"found group name %q; expected %q [path: %q]",
						result.MetaRecord.GroupName,
						cfg.GroupName(),
						result.MetaRecord.FQPath,
					)

					cfg.Log.Error().Err(err).
						Bool("username_check_enabled", resolveUsername).
						Bool("group_name_check_enabled", resolveGroupName).
						Str("path", path).
						Msg(errMsg.Error())

					nagiosExitState.LastError = err

					switch {
					case cfg.GroupNameCritical():
						nagiosExitState.ServiceOutput = fmt.Sprintf(
							"%s: %s",
							nagios.StateCRITICALLabel,
							errMsg.Error(),
						)
						nagiosExitState.ExitStatusCode = nagios.StateCRITICALExitCode

						return

					case cfg.GroupNameWarning():
						nagiosExitState.ServiceOutput = fmt.Sprintf(
							"%s: %s",
							nagios.StateWARNINGLabel,
							errMsg.Error(),
						)
						nagiosExitState.ExitStatusCode = nagios.StateWARNINGExitCode

						return
					}
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
	if resolveUsername {
		otherChecksApplied = append(otherChecksApplied, "username")
	}
	if resolveGroupName {
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
