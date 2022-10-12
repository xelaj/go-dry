// Copyright (c) 2022 Xelaj Software
//
// This file is a part of go-dry package.
// See https://github.com/xelaj/go-dry/blob/master/LICENSE for details

package dry

import (
	"path/filepath"
	"runtime"
)

var testMode bool

func IsTestMode() bool {
	return testMode
}

func TestGetCurrentPackagePath() string {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		panic("No caller information")
	}
	return filepath.Dir(filename)
}
