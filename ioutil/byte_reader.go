// Copyright (c) 2022 Xelaj Software
//
// This file is a part of go-dry package.
// See https://github.com/xelaj/go-dry/blob/master/LICENSE for details

package ioutil

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

func BytesReader(data interface{}) io.Reader {
	switch s := data.(type) {
	case io.Reader:
		return s
	case []byte:
		return bytes.NewReader(s)
	case string:
		return strings.NewReader(s)
	case fmt.Stringer:
		return strings.NewReader(s.String())
	case error:
		return strings.NewReader(s.Error())
	}
	return nil
}
