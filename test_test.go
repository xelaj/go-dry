// Copyright (c) 2022 Xelaj Software
//
// This file is a part of go-dry package.
// See https://github.com/xelaj/go-dry/blob/master/LICENSE for details

package dry

// see ./test.go doc how it works
var _ = func() bool {
	testMode = true
	return true
}()
