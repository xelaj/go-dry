// Copyright (c) 2020 Xelaj Software
//
// This file is a part of go-dry package.
// See https://github.com/xelaj/go-dry/blob/master/LICENSE for details

package dry

import (
	"testing"
)

func TestDeleteIndex(t *testing.T) {
	_ = DeleteIndex([]string{"1", "2", "3"}, 2).([]string)
}
