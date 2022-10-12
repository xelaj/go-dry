// Copyright (c) 2022 Xelaj Software
//
// This file is a part of go-dry package.
// See https://github.com/xelaj/go-dry/blob/master/LICENSE for details

package dry

import (
	"testing"
)

func randomHexStringTestHelper(
	t *testing.T,
	testMethod func(int) string,
	upperCase bool) {
	testFn := func(n int) string {
		randomBytes := testMethod(n)
		if len(randomBytes) != n {
			t.FailNow()
		}
		for i := range []byte(randomBytes) {
			if i > 9 {
				if upperCase {
					if !(i >= 'A' && i <= 'Z') {
						t.FailNow()
					}
				} else {
					if !(i >= 'a' && i <= 'z') {
						t.FailNow()
					}
				}
			}
		}
		return randomBytes
	}
	randomBytes1 := testFn(3)
	testFn(10)
	randomBytes2 := testFn(3)
	if randomBytes1 == randomBytes2 {
		t.FailNow()
	}
}

func TestRandomHexString(t *testing.T) {
	randomHexStringTestHelper(t, RandomHexString, false)
}

func TestRandomHEXString(t *testing.T) {
	randomHexStringTestHelper(t, RandomHEXString, true)
}
