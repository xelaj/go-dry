// Copyright (c) 2022 Xelaj Software
//
// This file is a part of go-dry package.
// See https://github.com/xelaj/go-dry/blob/master/LICENSE for details

package sys

import "testing"

func TestEndianSafeSplitUint16(t *testing.T) {
	least, most := EndianSafeSplitUint16(1)
	if !(least == 1 && most == 0) {
		t.Fail()
	}

	least, most = EndianSafeSplitUint16(256)
	if !(least == 0 && most == 1) {
		t.Fail()
	}
}
