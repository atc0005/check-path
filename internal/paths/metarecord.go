// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/check-path
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package paths

import (
	"os"
	"sort"
	"time"

	"github.com/atc0005/check-path/internal/units"
	"github.com/phayes/permbits"
)

// ID is a collection of username and group name values, associated with a
// specific os.FileInfo value.
type ID struct {
	Username  string
	UID       int
	UIDStr    string
	GroupName string
	GID       int
	GIDStr    string
}

// MetaRecord represents a superset of statistics for a file. This includes
// os.FileInfo, the fully-qualified path to the file and potentially
// user/group information.
type MetaRecord struct {

	// Size, name and other common details for the associated path.
	os.FileInfo

	// Permissions for the associated path.
	Permissions permbits.PermissionBits

	// Username and group values.
	ID

	// FQPath is the original, fully-qualified path.
	FQPath string

	// ParentDir is the parent directory for the path.
	ParentDir string
}

// MetaRecords is a slice of MetaRecord objects intended for bulk processing.
type MetaRecords []MetaRecord

// TotalFileSize returns the cumulative size of all MetaRecord objects in the
// slice in bytes. Since this would also apply to directories, those objects
// are filtered out in an effort to provide a more accurate value.
func (mr MetaRecords) TotalFileSize() int64 {

	var totalSize int64

	for _, file := range mr {

		// Skip any directory (MetaRecord) entries that may have been added.
		if file.IsDir() {
			continue
		}

		totalSize += file.Size()
	}

	return totalSize

}

// TotalFileSizeHR returns a human-readable string of the cumulative size of
// all files in the slice of bytes.
func (mr MetaRecords) TotalFileSizeHR() string {
	return units.ByteCountIEC(mr.TotalFileSize())
}

// SizeHR returns a human-readable string of the size of a MetaRecord object.
// Unless filtered later, this also applies to directories.
func (mr MetaRecord) SizeHR() string {
	return units.ByteCountIEC(mr.Size())
}

// AgeExceeded indicates whether a path is older than the specified threshold
// in days. If the path age is younger or equal to the specified number of
// days then the threshold is considered uncrossed.
func AgeExceeded(file os.FileInfo, days int) bool {

	var oldFile bool

	now := time.Now()
	fileModTime := file.ModTime()

	// Flip user specified number of days negative so that we can wind
	// back that many days from the file modification time. This gives
	// us our threshold to compare file modification times against.
	daysBack := -(days)
	fileAgeThreshold := now.AddDate(0, 0, daysBack)

	switch {
	case fileModTime.Before(fileAgeThreshold):
		oldFile = true
	case fileModTime.Equal(fileAgeThreshold):
		oldFile = false
	case fileModTime.After(fileAgeThreshold):
		oldFile = false
	}

	return oldFile

}

// SortByModTimeAsc sorts slice of MetaRecord objects in ascending order with
// older values listed first.
func (mr MetaRecords) SortByModTimeAsc() {
	sort.Slice(mr, func(i, j int) bool {
		return mr[i].ModTime().Before(mr[j].ModTime())
	})
}

// SortByModTimeDesc sorts slice of MetaRecord objects in descending order with
// newer values listed first.
func (mr MetaRecords) SortByModTimeDesc() {
	sort.Slice(mr, func(i, j int) bool {
		return mr[i].ModTime().After(mr[j].ModTime())
	})
}

// SortBySizeAsc sorts slice of MetaRecord objects in ascending order with
// smaller values listed first.
func (mr MetaRecords) SortBySizeAsc() {
	sort.Slice(mr, func(i, j int) bool {
		return mr[i].FileInfo.Size() > mr[j].FileInfo.Size()
	})
}

// SortBySizeDesc sorts slice of MetaRecord objects in descending order with
// larger values listed first.
func (mr MetaRecords) SortBySizeDesc() {
	sort.Slice(mr, func(i, j int) bool {
		return mr[i].FileInfo.Size() < mr[j].FileInfo.Size()
	})
}
