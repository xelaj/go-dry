// Copyright (c) 2020 Xelaj Software
//
// This file is a part of go-dry package.
// See https://github.com/xelaj/go-dry/blob/master/LICENSE for details

package dry

import (
	"testing"
)

func Test_EncryptionAES(t *testing.T) {
	key := []byte("0123456789ABCDEF")
	data := "Hello World 1234"

	ciphertext := EncryptAES(key, []byte(data))
	if string(DecryptAES(key, ciphertext)) != data {
		t.Fail()
	}
}
