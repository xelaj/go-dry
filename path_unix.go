// Copyright (c) 2022 Xelaj Software
//
// This file is a part of go-dry package.
// See https://github.com/xelaj/go-dry/blob/master/LICENSE for details

// +build !windows

package dry

import (
	"path/filepath"

	"golang.org/x/sys/unix"
)

func pathIsWirtable(path string) bool {
	inspectPath, _ := filepath.Abs(path)

	nearestPath := PathNearestExisting(inspectPath)
	if nearestPath == "" {
		return false
	}

	return unix.Access(nearestPath, unix.W_OK) == nil
}
