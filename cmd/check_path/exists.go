// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/check-path
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package main

import (
	"errors"
	"fmt"

	"github.com/atc0005/check-path/internal/paths"
	"github.com/atc0005/go-nagios"
	"github.com/rs/zerolog"
)

func checkPaths(list []string, critical bool, warning bool, zlog *zerolog.Logger, nes *nagios.ExitState) {

	pathInfo, err := paths.AssertNotExists(list)

	switch {

	// error: path found
	case errors.Is(err, paths.ErrPathExists):

		nes.LastError = err
		nes.LongServiceOutput += fmt.Sprintf(
			"* Path %q%s** Last Modified: %v%s",
			pathInfo.FQPath,
			nagios.CheckOutputEOL,
			pathInfo.ModTime(),
			nagios.CheckOutputEOL,
		)

		if critical {
			nes.ServiceOutput = fmt.Sprintf(
				"%s: %v",
				nagios.StateCRITICALLabel,
				err,
			)
			nes.ExitStatusCode = nagios.StateCRITICALExitCode
		}

		if warning {
			nes.ServiceOutput = fmt.Sprintf(
				"%s: %v",
				nagios.StateWARNINGLabel,
				err,
			)
			nes.ExitStatusCode = nagios.StateWARNINGExitCode
		}

		return

	// error: other
	case err != nil:
		zlog.Err(err).Msg("Error checking path")
		nes.LastError = err
		nes.ServiceOutput = fmt.Sprintf(
			"%s: Error checking path: %v",
			nagios.StateCRITICALLabel,
			err,
		)
		nes.ExitStatusCode = nagios.StateCRITICALExitCode

		return

	// OK: desired state; paths not found
	case err == nil:

		nes.LastError = nil
		nes.ServiceOutput = fmt.Sprintf(
			"%s: Specified (unwanted) paths do not exist",
			nagios.StateOKLabel,
		)

		zlog.Info().Msg(
			"OK: All specified paths were not found; expected result per " +
				"exists-critical or exists-warning flag",
		)

		nes.ExitStatusCode = nagios.StateOKExitCode

		return

	}

}
