// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/check-path
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package main

import (
	"fmt"
	"time"

	"github.com/atc0005/check-path/internal/config"
	"github.com/atc0005/check-path/internal/paths"
	"github.com/atc0005/go-nagios"
	"github.com/rs/zerolog"
)

// checkAge is a helper variadic function that accepts one or many MetaRecord
// values for age evaluation. If the specified age threshold values are
// crossed, the provided *nagios.Plugin is updated and an error is returned
// to signal that this specific check has found old files.
func checkAge(path string, ths config.FileAgeThresholds, zlog *zerolog.Logger, nes *nagios.Plugin, mrs ...paths.MetaRecord) error {

	// type conversion to expose desired methods
	metaRecords := paths.MetaRecords(mrs)

	// sort if more than one entry
	if len(mrs) > 1 {
		metaRecords.SortByModTimeAsc()
	}

	for _, record := range metaRecords {

		// skip age check for directories
		if record.IsDir() {
			continue
		}

		criticalAgeFile := paths.AgeExceeded(
			record.FileInfo, ths.Critical)

		warningAgeFile := paths.AgeExceeded(
			record.FileInfo, ths.Warning)

		if criticalAgeFile || warningAgeFile {
			zlog.Error().Err(paths.ErrPathOldFilesFound).
				Int("critical_age_days", ths.Critical).
				Int("warning_age_days", ths.Warning).
				Bool("age_check_enabled", ths.Set).
				Str("path", path).
				Msg("old files found")

			nes.AddError(fmt.Errorf(
				"%d files & directories evaluated: %w",
				len(metaRecords),
				paths.ErrPathOldFilesFound,
			))

			fileAge := time.Since(record.ModTime()).Hours() / 24

			nes.LongServiceOutput += fmt.Sprintf(
				"* File %s** parent dir: %q%s** name: %q%s** age: %v%s",
				nagios.CheckOutputEOL,
				record.ParentDir,
				nagios.CheckOutputEOL,
				record.Name(),
				nagios.CheckOutputEOL,
				fileAge,
				nagios.CheckOutputEOL,
			)

			switch {
			case criticalAgeFile:
				nes.ServiceOutput = fmt.Sprintf(
					"%s: file older than %d days (%.2f) found [path: %q]",
					nagios.StateCRITICALLabel,
					ths.Critical,
					fileAge,
					path,
				)

				nes.ExitStatusCode = nagios.StateCRITICALExitCode

				return paths.ErrPathOldFilesFound

			case warningAgeFile:
				nes.ServiceOutput = fmt.Sprintf(
					"%s: file older than %d days (%.2f) found [path: %q]",
					nagios.StateWARNINGLabel,
					ths.Warning,
					fileAge,
					path,
				)

				nes.ExitStatusCode = nagios.StateWARNINGExitCode

				return paths.ErrPathOldFilesFound
			}

		}
	}

	return nil

}
