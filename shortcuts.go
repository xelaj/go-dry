// Copyright (c) 2020 Xelaj Software
//
// This file is a part of go-dry package.
// See https://github.com/xelaj/go-dry/blob/master/LICENSE for details

package dry

// FirstArg returns the first passed argument,
// can be used to extract first result value
// from a function call to pass it on to functions like fmt.Printf
func FirstArg(args ...any) any {
	return args[0]
}
