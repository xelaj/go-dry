// Copyright (c) 2022 Xelaj Software
//
// This file is a part of go-dry package.
// See https://github.com/xelaj/go-dry/blob/master/LICENSE for details

package dry

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"runtime"
	"strings"
)

func StackTrace(skipFrames int) string {
	var b strings.Builder
	var lastFile string
	for i := 3; ; i++ {
		contin := fprintStackTraceLine(i, &lastFile, &b)
		if !contin {
			break
		}
	}
	return b.String()
}

func StackTraceLine(skipFrames int) string {
	var b strings.Builder
	var lastFile string
	fprintStackTraceLine(skipFrames, &lastFile, &b)
	return b.String()
}

func fprintStackTraceLine(i int, lastFile *string, b io.Writer) bool {
	var lines [][]byte

	pc, file, line, ok := runtime.Caller(i)
	if !ok {
		return false
	}

	// Print this much at least.  If we can't find the source, it won't show.
	fmt.Fprintf(b, "%s:%d (0x%x)\n", file, line, pc)
	if file != *lastFile {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			return true
		}
		lines = bytes.Split(data, []byte{'\n'})
		*lastFile = file
	}
	line-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	fmt.Fprintf(b, "\t%s: %s\n", function(pc), source(lines, line))
	return true
}

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	sep       = []byte("/")
)

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.Trim(lines[n], " \t")
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod

	if period := bytes.LastIndex(name, sep); period >= 0 {
		name = name[period+1:]
	}

	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.ReplaceAll(name, centerDot, dot)
	return name
}
