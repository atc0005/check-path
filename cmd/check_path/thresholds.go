// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/check-path
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package main

import (
	"fmt"
	"strings"

	"github.com/atc0005/check-path/internal/config"
	"github.com/atc0005/check-path/internal/units"
	"github.com/atc0005/go-nagios"
)

// setThresholdDescriptions is a helper function for conditionally setting
// CRITICAL and WARNING threshold descriptions based on user-specified flags
// and values.
func setThresholdDescriptions(cfg *config.Config, nes *nagios.ExitState) {
	if cfg.PathExistsCritical() {
		nes.CriticalThreshold = "[Paths exist]"
		nes.WarningThreshold = "N/A"
	}

	if cfg.PathExistsWarning() {
		nes.CriticalThreshold = "N/A"
		nes.WarningThreshold = "[Paths exist]"
	}

	if age := cfg.Age(); age.Set {
		ageCriticalThreshold := fmt.Sprintf(
			"[File age in days: %d]",
			age.Critical,
		)

		switch {
		case nes.CriticalThreshold != "":
			nes.CriticalThreshold = strings.Join(
				[]string{
					nes.CriticalThreshold,
					ageCriticalThreshold,
				},
				", ",
			)
		default:
			nes.CriticalThreshold = ageCriticalThreshold
		}

		ageWarningThreshold := fmt.Sprintf(
			"[File age in days: %d]",
			age.Warning,
		)

		switch {
		case nes.WarningThreshold != "":
			nes.WarningThreshold = strings.Join(
				[]string{
					nes.WarningThreshold,
					ageWarningThreshold,
				},
				", ",
			)
		default:
			nes.WarningThreshold = ageWarningThreshold
		}

	}

	if size := cfg.Size(); size.Set {

		sizeCriticalThreshold := fmt.Sprintf(
			"[File size (bytes: %d, Human: %s)]",
			size.Critical,
			units.ByteCountIEC(size.Critical),
		)

		switch {
		case nes.CriticalThreshold != "":
			nes.CriticalThreshold = strings.Join(
				[]string{
					nes.CriticalThreshold,
					sizeCriticalThreshold,
				},
				", ",
			)
		default:
			nes.CriticalThreshold = sizeCriticalThreshold
		}

		sizeWarningThreshold := fmt.Sprintf(
			"[File size (bytes: %d, Human: %s)]",
			size.Warning,
			units.ByteCountIEC(size.Warning),
		)

		switch {
		case nes.WarningThreshold != "":
			nes.WarningThreshold = strings.Join(
				[]string{
					nes.WarningThreshold,
					sizeWarningThreshold,
				},
				", ",
			)
		default:
			nes.WarningThreshold = sizeWarningThreshold
		}

	}

}
