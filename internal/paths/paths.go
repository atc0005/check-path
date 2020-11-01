// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/check-path
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package paths

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Application-specific errors for common path checks.
var (
	ErrPathExists        = errors.New("path exists")
	ErrPathDoesNotExist  = errors.New("path does not exist")
	ErrPathEmptyString   = errors.New("specified path is empty string")
	ErrPathCheckFailed   = errors.New("failed to check path")
	ErrPathCheckCanceled = errors.New("path check canceled")
	ErrPathOldFilesFound = errors.New("old files found in path")
)

// ProcessResult is a superset of a MetaRecord and any associated error
// encountered while processing a path.
type ProcessResult struct {
	MetaRecord
	Error error
}

// AssertNotExists accepts a list of paths to process and returns the specific
// error if unable to check the path, a set MetaRecord value and a package
// specific associated error for a path that is found or an empty MetaRecord and
// nil if all paths do not exist.
func AssertNotExists(list []string) (MetaRecord, error) {
	for _, path := range list {

		pathInfo, err := Info(path)

		switch {

		// desired state
		case errors.Is(err, ErrPathDoesNotExist):
			continue

		// some other error occurred
		case err != nil:
			return MetaRecord{}, fmt.Errorf(
				"unable to assert non-existence of %s: %w",
				path,
				err,
			)

		// path found
		case err == nil:

			var pathType string
			pathType = "file"
			if pathInfo.IsDir() {
				pathType = "directory"
			}

			return pathInfo, fmt.Errorf(
				"%s %q: %w",
				pathType, path, ErrPathExists,
			)

		}

	}

	// no errors, paths not found
	return MetaRecord{}, nil
}

// Info is a helper function used to quickly gather details on a specified
// path. A MetaRecord value and nil is returned for successful path evaluation,
// otherwise an empty MetaRecord value and appropriate error is returned.
func Info(path string) (MetaRecord, error) {

	// Make sure path isn't empty
	if strings.TrimSpace(path) == "" {
		return MetaRecord{}, ErrPathEmptyString
	}

	pathInfo, statErr := os.Stat(path)
	if statErr != nil {
		if !os.IsNotExist(statErr) {
			// ERROR: another error occurred aside from file not found
			return MetaRecord{}, fmt.Errorf(
				"error checking path %s: %w",
				path,
				statErr,
			)
		}
		// path not found
		return MetaRecord{}, fmt.Errorf("%w: %s", ErrPathDoesNotExist, path)
	}

	// path found
	return MetaRecord{
		FileInfo:  pathInfo,
		FQPath:    path,
		ParentDir: filepath.Dir(path),
	}, nil
}

// Exists is a helper function used to quickly determine whether a specified
// path exists.
func Exists(path string) (bool, error) {

	// Make sure path isn't empty
	if strings.TrimSpace(path) == "" {
		return false, fmt.Errorf("specified path is empty string")
	}

	_, statErr := os.Stat(path)
	if statErr != nil {
		if !os.IsNotExist(statErr) {
			// ERROR: another error occurred aside from file not found
			return false, fmt.Errorf(
				"error checking path %s: %w",
				path,
				statErr,
			)
		}
		// file not found
		return false, nil
	}

	// file found
	return true, nil
}

// Process evalutes the specified path, either at a flat level or if
// specified, recursively. ProcessResult values are sent back by way of a
// results channel.
func Process(ctx context.Context, path string, recurse bool, results chan<- ProcessResult) {

	// NOTE: This is safe to close *ONLY* because we recreate the channel on
	// each iteration of the specified paths (e.g., one path at a time) before
	// calling this function.
	defer close(results)

	// Qualify the path for later use, otherwise bail.
	fqPath, absErr := filepath.Abs(path)
	if absErr != nil {
		results <- ProcessResult{
			Error: absErr,
		}
		return
	}

	walkErr := filepath.Walk(fqPath, func(path string, info os.FileInfo, err error) error {

		// If we return a non-nil error, this will stop the filepath.Walk()
		// function from continuing to walk the path.
		switch {

		// error: context canceled
		case ctx.Err() != nil:
			results <- ProcessResult{
				Error: fmt.Errorf("%s: %w", ctx.Err().Error(), ErrPathCheckCanceled),
			}
			return ctx.Err()

		// error: path does not exist
		case os.IsNotExist(err):
			return fmt.Errorf("%w: %s", ErrPathDoesNotExist, path)

		// error: other
		case err != nil:
			return err

		// is a directory & not fully-qualified, specified path; skip if
		// recurse is not enabled
		case info.IsDir() && path != fqPath:
			if !recurse {
				return filepath.SkipDir
			}
		}

		// send back metadata for further processing
		results <- ProcessResult{
			MetaRecord: MetaRecord{
				FileInfo:  info,
				FQPath:    path,
				ParentDir: filepath.Dir(path),
			},
		}

		// indicate no error to filepath.Walk() so that it will continue to
		// the next item in the path (if applicable)
		return nil
	})

	if walkErr != nil {
		results <- ProcessResult{
			Error: fmt.Errorf("error examining path %q: %w", path, walkErr),
		}
	}

}
