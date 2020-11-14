// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/check-path
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package main

import (
	"fmt"

	"github.com/atc0005/check-path/internal/config"
	"github.com/atc0005/check-path/internal/paths"
	"github.com/atc0005/go-nagios"
	"github.com/rs/zerolog"
)

// checkSize is a helper variadic function that accepts one or many MetaRecord
// values for size evaluation. If the specified size threshold values are
// crossed, the provided *nagios.ExitState is updated and an error is returned
// to signal that this specific check has found files which exceed specified
// size values (either too small or too large).
func checkSize(path string, ths config.FileSizeThresholdsMinMax, zlog *zerolog.Logger, nes *nagios.ExitState, mrs ...paths.MetaRecord) error {

	// type conversion to expose desired methods
	metaRecords := paths.MetaRecords(mrs)

	// sort if more than one entry
	if len(mrs) > 1 {
		metaRecords.SortBySizeAsc()
	}

	// NOTE: Directories (themselves) are not included in the size values,
	// just the contents of said directories.
	actualSizeHR := metaRecords.TotalFileSizeHR()
	actualSizeBytes := metaRecords.TotalFileSize()

	// warning threshold required, so we can use that to reduce
	// conditional check logic complexity
	if (ths.SizeMin.Set && actualSizeBytes < ths.SizeMin.Warning) ||
		(ths.SizeMax.Set && actualSizeBytes > ths.SizeMax.Warning) {

		sizeOfFilesTooLargeErr := fmt.Errorf(
			"%w (%d evaluated)",
			paths.ErrSizeOfFilesTooLarge,
			len(metaRecords),
		)

		sizeOfFilesTooSmallErr := fmt.Errorf(
			"%w (%d evaluated)",
			paths.ErrSizeOfFilesTooSmall,
			len(metaRecords),
		)

		serviceOutputTmpl := fmt.Sprintf(
			"size threshold crossed; %v found in path %q",
			actualSizeHR,
			path,
		)

		nes.LongServiceOutput += fmt.Sprintf(
			"* Size %s** path: %q%s** bytes: %v%s** human-readable: %v%s",
			nagios.CheckOutputEOL,
			path,
			nagios.CheckOutputEOL,
			actualSizeBytes,
			nagios.CheckOutputEOL,
			actualSizeHR,
			nagios.CheckOutputEOL,
		)

		var stateLabel string
		var exitCode int

		// configure exit state details based on how the
		// thresholds were crossed. return after all exit state
		// details are recorded
		switch {

		case ths.SizeMax.Set &&
			(actualSizeBytes > ths.SizeMax.Critical || actualSizeBytes > ths.SizeMax.Warning):
			zlog.Error().Err(sizeOfFilesTooLargeErr).
				Int64("critical_size_max_bytes", ths.SizeMax.Critical).
				Int64("warning_size_max_bytes", ths.SizeMax.Warning).
				Int64("actual_size_bytes", actualSizeBytes).
				Str("actual_size_hr", actualSizeHR).
				Bool("size_max_check_enabled", ths.SizeMax.Set).
				Str("path", path).
				Msg(sizeOfFilesTooLargeErr.Error())

			nes.LastError = sizeOfFilesTooLargeErr

			if actualSizeBytes > ths.SizeMax.Critical {
				stateLabel = nagios.StateCRITICALLabel
				exitCode = nagios.StateCRITICALExitCode
			}

			if actualSizeBytes > ths.SizeMax.Warning {
				stateLabel = nagios.StateWARNINGLabel
				exitCode = nagios.StateWARNINGExitCode
			}

			nes.ServiceOutput = fmt.Sprintf(
				"%s: %s %s",
				stateLabel,
				ths.SizeMax.Description,
				serviceOutputTmpl,
			)

			nes.ExitStatusCode = exitCode

			return sizeOfFilesTooLargeErr

		case ths.SizeMin.Set &&
			(actualSizeBytes < ths.SizeMin.Critical || actualSizeBytes < ths.SizeMin.Warning):
			zlog.Error().Err(sizeOfFilesTooSmallErr).
				Int64("critical_size_min_bytes", ths.SizeMin.Critical).
				Int64("warning_size_min_bytes", ths.SizeMin.Warning).
				Int64("actual_size_bytes", actualSizeBytes).
				Str("actual_size_hr", actualSizeHR).
				Bool("size_min_check_enabled", ths.SizeMin.Set).
				Str("path", path).
				Msg(sizeOfFilesTooSmallErr.Error())

			nes.LastError = sizeOfFilesTooSmallErr

			if actualSizeBytes < ths.SizeMin.Critical {
				stateLabel = nagios.StateCRITICALLabel
				exitCode = nagios.StateCRITICALExitCode
			}

			if actualSizeBytes < ths.SizeMin.Warning {
				stateLabel = nagios.StateWARNINGLabel
				exitCode = nagios.StateWARNINGExitCode
			}

			nes.ServiceOutput = fmt.Sprintf(
				"%s: %s %s",
				stateLabel,
				ths.SizeMin.Description,
				serviceOutputTmpl,
			)

			nes.ExitStatusCode = exitCode

			return sizeOfFilesTooSmallErr

		}

	}

	return nil

}
