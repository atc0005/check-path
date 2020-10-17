// +build windows

// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/check-path
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package paths

import (
	"fmt"
	"os"
)

/*
	NOTE: These operations are currently not supported for Windows.

	While we use config validation in an effort to "fail fast" for unsupported
	operations, we also mock functions here to satisfy any calls not protected
	by build tags.
*/

// ResolveIDs accepts a MetaRecord pointer and if supported, resolves the uid
// and gid values from the underlying syscall.Stat_t and sets the values
// directly using the provided MetaRecord pointer. For unsupported operating
// systems, this function returns a hard-coded error.
func ResolveIDs(mr *MetaRecord) error {
	return fmt.Errorf("ResolveIDs unavailable; unsupported operating system")
}

// LookupIDs accepts a os.FileInfo and if supported, resolves the uid and gid
// values from the underlying syscall.Stat_t and returns those resolved values
// as an instance of the ID type for further processing. For unsupported
// operating systems, this function returns a hard-coded error.
func LookupIDs(fi os.FileInfo) (ID, error) {
	return ID{}, fmt.Errorf("LookupIDs unavailable; unsupported operating system")
}
