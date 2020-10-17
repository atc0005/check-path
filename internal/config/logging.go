// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/check-path
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package config

import (
	"os"

	"github.com/rs/zerolog"
)

const (

	// LogLevelDisabled maps to zerolog.Disabled logging level
	LogLevelDisabled string = "disabled"

	// LogLevelPanic maps to zerolog.PanicLevel logging level
	LogLevelPanic string = "panic"

	// LogLevelFatal maps to zerolog.FatalLevel logging level
	LogLevelFatal string = "fatal"

	// LogLevelError maps to zerolog.ErrorLevel logging level
	LogLevelError string = "error"

	// LogLevelWarn maps to zerolog.WarnLevel logging level
	LogLevelWarn string = "warn"

	// LogLevelInfo maps to zerolog.InfoLevel logging level
	LogLevelInfo string = "info"

	// LogLevelDebug maps to zerolog.DebugLevel logging level
	LogLevelDebug string = "debug"

	// LogLevelTrace maps to zerolog.TraceLevel logging level
	LogLevelTrace string = "trace"
)

// configureLogging is a wrapper function to enable setting requested logging
// settings. This is called *after* configuration validation has been
// performed in order to reject any invalid user-provided settings.
func (c *Config) configureLogging() {

	switch c.LogLevel() {
	case LogLevelDisabled:
		zerolog.SetGlobalLevel(zerolog.Disabled)
	case LogLevelPanic:
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case LogLevelFatal:
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case LogLevelError:
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case LogLevelWarn:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case LogLevelInfo:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case LogLevelDebug:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case LogLevelTrace:
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}

	// c.Log = zerolog.New(os.Stderr).With().Caller().
	// 	Str("version", version).
	// 	Str("logging_level", config.LoggingLevel).
	// 	Str("server", config.Server).
	// 	Int("port", config.Port).
	// 	Int("age_warning", config.AgeWarning).
	// 	Int("age_critical", config.AgeCritical).
	// 	Str("expected_sans_entries", config.SANsEntries.String()).Logger()

	c.Log = zerolog.New(os.Stderr).With().Caller().
		Str("version", version).
		Str("logging_level", c.LogLevel()).Logger()

}
