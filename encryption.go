// Copyright (c) 2022 Xelaj Software
//
// This file is a part of go-dry package.
// See https://github.com/xelaj/go-dry/blob/master/LICENSE for details

package dry

import (
	"crypto/sha1"
)

func Sha1(input string) []byte {
	r := sha1.Sum([]byte(input))
	return r[:]
}

func Sha1Byte(input []byte) []byte {
	r := sha1.Sum(input)
	return r[:]
}
