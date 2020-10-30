// Copyright (c) 2020 Xelaj Software
//
// This file is a part of go-dry package.
// See https://github.com/xelaj/go-dry/blob/master/LICENSE for details

package dry

import "testing"

func TestEnvironMap(t *testing.T) {
	ret := environToMap([]string{
		"a=b",
		"b=c=d",
	})

	if len(ret) != 2 {
		t.Fail()
	}

	if ret["a"] != "b" {
		t.Fail()
	}

	if ret["b"] != "c=d" {
		t.Fail()
	}
}

func TestGetenvDefault(t *testing.T) {
	if GetenvDefault("GO_DRY_BOGUS_ENVIRONMENT_VARIABLE", "default") != "default" {
		t.Fail()
	}
}
