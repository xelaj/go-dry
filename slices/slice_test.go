// Copyright (c) 2022 Xelaj Software
//
// This file is a part of go-dry package.
// See https://github.com/xelaj/go-dry/blob/master/LICENSE for details

package slices

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteIndex(t *testing.T) {
	_ = DeleteIndex([]string{"1", "2", "3"}, 2).([]string)
}

func TestSliceRemoveFunc(t *testing.T) {
	tests := []struct {
		slice any
		f     func(i any) bool
		want  any
	}{
		{
			[]string{"a", "b", "", "c"},
			func(i any) bool { return i.(string) != "" },
			[]string{"a", "b", "c"},
		},
		{
			[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, -1},
			func(i any) bool { return i.(int) == 0 },
			[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			assert.Equal(t, tt.want, RemoveFunc(tt.slice, tt.f))
		})
	}
}

func TestSlicePopFunc(t *testing.T) {
	tests := []struct {
		slice      any
		f          func(i any) bool
		wantRes    any
		wantPopped any
	}{
		{
			[]string{"a", "b", "", "c"},
			func(i any) bool { return i.(string) == "" },
			[]string{"a", "b", "c"},
			[]string{""},
		},
		{
			[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, -1},
			func(i any) bool { return i.(int) != 0 },
			[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			[]int{-1},
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			res, popped := PopFunc(tt.slice, tt.f)
			assert.Equal(t, tt.wantRes, res)
			assert.Equal(t, tt.wantPopped, popped)
		})
	}
}
