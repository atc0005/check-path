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

// checkIDs is a helper variadic function that accepts one or many MetaRecord
// values for username, group name evaluation. If the specified values are not
// present, the provided *nagios.ExitState is updated and an error is returned
// to signal that this specific check has found missing username or group name
// values.
func checkIDs(path string, resolveIDs config.ResolveIDs, zlog *zerolog.Logger, nes *nagios.ExitState, mrs ...paths.MetaRecord) error {

	// G601: Implicit memory aliasing in for loop. (gosec)
	// for _, record := range mrs {

	for i := range mrs {

		// Avoid taking the address of a loop variable and instead use
		// indexing to get at the original memory address for use with
		// paths.ResolveIDs
		// G601: Implicit memory aliasing in for loop. (gosec)
		// https://stackoverflow.com/questions/62446118/implicit-memory-aliasing-in-for-loop
		record := mrs[i]

		zlog.Debug().Msg("Username, Group name resolution enabled")

		resolveErr := paths.ResolveIDs(&record)
		if resolveErr != nil {
			zlog.Error().Err(resolveErr).
				Str("path", path).
				Msg(resolveErr.Error())

			nes.AddError(resolveErr)
			nes.ServiceOutput = fmt.Sprintf(
				"%s: failed to resolve IDs: %v [path: %q]",
				nagios.StateCRITICALLabel,
				resolveErr.Error(),
				path,
			)

			return resolveErr
		}

		if resolveIDs.UsernameCheck &&
			record.Username != resolveIDs.Username {

			statusMsg := fmt.Sprintf(
				"found username %q; expected %q [path: %q]",
				record.Username,
				resolveIDs.Username,
				record.FQPath,
			)

			zlog.Error().Err(paths.ErrPathMissingUsername).
				Bool("username_check_enabled", resolveIDs.UsernameCheck).
				Bool("group_name_check_enabled", resolveIDs.GroupNameCheck).
				Str("path", path).
				Msg(statusMsg)

			nes.AddError(paths.ErrPathMissingUsername)

			switch {
			case resolveIDs.UsernameCritical:
				nes.ServiceOutput = fmt.Sprintf(
					"%s: %s",
					nagios.StateCRITICALLabel,
					statusMsg,
				)
				nes.ExitStatusCode = nagios.StateCRITICALExitCode

				return paths.ErrPathMissingUsername

			case resolveIDs.UsernameWarning:
				nes.ServiceOutput = fmt.Sprintf(
					"%s: %s",
					nagios.StateWARNINGLabel,
					statusMsg,
				)
				nes.ExitStatusCode = nagios.StateWARNINGExitCode

				return paths.ErrPathMissingUsername
			}
		}

		if resolveIDs.GroupNameCheck &&
			record.GroupName != resolveIDs.GroupName {

			statusMsg := fmt.Sprintf(
				"found group name %q; expected %q [path: %q]",
				record.GroupName,
				resolveIDs.GroupName,
				record.FQPath,
			)

			zlog.Error().Err(paths.ErrPathMissingGroupName).
				Bool("username_check_enabled", resolveIDs.UsernameCheck).
				Bool("group_name_check_enabled", resolveIDs.GroupNameCheck).
				Str("path", path).
				Msg(statusMsg)

			nes.AddError(paths.ErrPathMissingGroupName)

			switch {
			case resolveIDs.GroupNameCritical:
				nes.ServiceOutput = fmt.Sprintf(
					"%s: %s",
					nagios.StateCRITICALLabel,
					statusMsg,
				)
				nes.ExitStatusCode = nagios.StateCRITICALExitCode

				return paths.ErrPathMissingGroupName

			case resolveIDs.GroupNameWarning:
				nes.ServiceOutput = fmt.Sprintf(
					"%s: %s",
					nagios.StateWARNINGLabel,
					statusMsg,
				)
				nes.ExitStatusCode = nagios.StateWARNINGExitCode

				return paths.ErrPathMissingGroupName
			}
		}
	}

	return nil
}
