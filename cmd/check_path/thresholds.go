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
func setThresholdDescriptions(cfg *config.Config, nes *nagios.Plugin) {
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

	if sizeMin := cfg.SizeMin(); sizeMin.Set {

		sizeMinCriticalThreshold := fmt.Sprintf(
			"[Min File size (bytes: %d, Human: %s)]",
			sizeMin.Critical,
			units.ByteCountIEC(sizeMin.Critical),
		)

		switch {
		case nes.CriticalThreshold != "":
			nes.CriticalThreshold = strings.Join(
				[]string{
					nes.CriticalThreshold,
					sizeMinCriticalThreshold,
				},
				", ",
			)
		default:
			nes.CriticalThreshold = sizeMinCriticalThreshold
		}

		sizeMinWarningThreshold := fmt.Sprintf(
			"[Min File size (bytes: %d, Human: %s)]",
			sizeMin.Warning,
			units.ByteCountIEC(sizeMin.Warning),
		)

		switch {
		case nes.WarningThreshold != "":
			nes.WarningThreshold = strings.Join(
				[]string{
					nes.WarningThreshold,
					sizeMinWarningThreshold,
				},
				", ",
			)
		default:
			nes.WarningThreshold = sizeMinWarningThreshold
		}

	}

	if sizeMax := cfg.SizeMax(); sizeMax.Set {

		sizeMaxCriticalThreshold := fmt.Sprintf(
			"[Max File size (bytes: %d, Human: %s)]",
			sizeMax.Critical,
			units.ByteCountIEC(sizeMax.Critical),
		)

		switch {
		case nes.CriticalThreshold != "":
			nes.CriticalThreshold = strings.Join(
				[]string{
					nes.CriticalThreshold,
					sizeMaxCriticalThreshold,
				},
				", ",
			)
		default:
			nes.CriticalThreshold = sizeMaxCriticalThreshold
		}

		sizeMaxWarningThreshold := fmt.Sprintf(
			"[Max File size (bytes: %d, Human: %s)]",
			sizeMax.Warning,
			units.ByteCountIEC(sizeMax.Warning),
		)

		switch {
		case nes.WarningThreshold != "":
			nes.WarningThreshold = strings.Join(
				[]string{
					nes.WarningThreshold,
					sizeMaxWarningThreshold,
				},
				", ",
			)
		default:
			nes.WarningThreshold = sizeMaxWarningThreshold
		}

	}

}
