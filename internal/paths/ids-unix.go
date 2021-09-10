//go:build !windows
// +build !windows

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
	"os/user"
	"strconv"
	"syscall"
)

// ResolveIDs accepts a MetaRecord pointer and if supported, resolves the uid
// and gid values from the underlying syscall.Stat_t and sets the values
// directly using the provided MetaRecord pointer. For unsupported operating
// systems, this function is effectively a NOOP.
func ResolveIDs(mr *MetaRecord) error {

	if mr == nil {
		return fmt.Errorf("received nil MetaRecord pointer")
	}

	// FIXME: Refactor so that we can call LookupIDs instead of duplicating the
	// logic below?
	stat, statOK := mr.Sys().(*syscall.Stat_t)
	if !statOK {
		return fmt.Errorf("failed to access syscall.Stat_t")
	}

	mr.UID = int(stat.Uid)
	mr.GID = int(stat.Gid)

	mr.GIDStr = strconv.Itoa(mr.GID)
	var groupLookupErr error
	groupResult, groupLookupErr := user.LookupGroupId(mr.GIDStr)
	if groupLookupErr != nil {
		return fmt.Errorf("failed to resolve gid to group name: %w", groupLookupErr)
	}
	mr.GroupName = groupResult.Name

	mr.UIDStr = strconv.Itoa(mr.UID)
	var userLookupErr error
	userResult, userLookupErr := user.LookupId(mr.UIDStr)
	if userLookupErr != nil {
		return fmt.Errorf("failed to resolve uid to username: %w", userLookupErr)
	}
	mr.Username = userResult.Username

	return nil

}

// LookupIDs accepts a os.FileInfo and if supported, resolves the uid and gid
// values from the underlying syscall.Stat_t and returns those resolved values
// as an instance of the ID type for further processing. For unsupported
// operating systems, this function is effectively a NOOP.
func LookupIDs(fi os.FileInfo) (ID, error) {
	stat, statOK := fi.Sys().(*syscall.Stat_t)
	if !statOK {
		return ID{}, fmt.Errorf("failed to access syscall.Stat_t")
	}

	id := ID{}

	id.UID = int(stat.Uid)
	id.GID = int(stat.Gid)

	id.GIDStr = strconv.Itoa(id.GID)
	groupResult, err := user.LookupGroupId(id.GIDStr)
	if err != nil {
		return ID{}, fmt.Errorf("failed to resolve gid to group name: %w", err)
	}
	id.GroupName = groupResult.Name

	id.UIDStr = strconv.Itoa(id.UID)
	userResult, err := user.LookupId(id.UIDStr)
	if err != nil {
		return ID{}, fmt.Errorf("failed to resolve uid to username: %w", err)
	}
	id.Username = userResult.Username

	return id, nil

}
