// Copyright (c) 2020 Xelaj Software
//
// This file is a part of go-dry package.
// See https://github.com/xelaj/go-dry/blob/master/LICENSE for details

package dry

import (
	"bytes"
	"reflect"
	"sync"
	"testing"
)

func assertStringsEqual(t *testing.T, str1, str2 string) {
	if str1 != str2 {
		t.FailNow()
	}
}

func Test_BytesMD5(t *testing.T) {
	assertStringsEqual(
		t, "5d41402abc4b2a76b9719d911017c592",
		BytesMD5("hello"))
}

func Test_BytesEncodeBase64(t *testing.T) {
	assertStringsEqual(
		t, "aGVsbG8=",
		BytesEncodeBase64("hello"))
}

func Test_BytesDecodeBase64(t *testing.T) {
	assertStringsEqual(
		t, "hello",
		BytesDecodeBase64("aGVsbG8="))
}

func Test_BytesEncodeHex(t *testing.T) {
	assertStringsEqual(
		t, "68656c6c6f",
		BytesEncodeHex("hello"))
}

func Test_BytesDecodeHex(t *testing.T) {
	assertStringsEqual(
		t, "hello",
		BytesDecodeHex("68656C6C6F"))
}

//nolint:staticcheck // strongly protected by WaitGroup
func testCompressDecompress(t *testing.T,
	compressFunc func([]byte) []byte,
	decompressFunc func([]byte) []byte) {
	testFn := func(wg *sync.WaitGroup, testData []byte) {
		defer wg.Done()
		compressedData := compressFunc(testData)
		uncompressedData := decompressFunc(compressedData)
		if !bytes.Equal(testData, uncompressedData) {
			t.FailNow()
		}
	}
	wg := new(sync.WaitGroup)
	wg.Add(3)

	go testFn(wg, []byte("hello123"))
	go testFn(wg, []byte("gopher456"))
	go testFn(wg, []byte("dry789"))

	wg.Wait()
}

func Test_BytesDeflateInflate(t *testing.T) {
	testCompressDecompress(t, BytesDeflate, BytesInflate)
}

func Test_BytesGzipUnGzip(t *testing.T) {
	testCompressDecompress(t, BytesGzip, BytesUnGzip)
}

func bytesHeadTailTestHelper(
	t *testing.T,
	testMethod func([]byte, int) ([]string, []byte),
	lines []byte, n int,
	expectedLines []string, expectedRest []byte) {
	resultLines, resultRest := testMethod(lines, n)
	if !reflect.DeepEqual(resultLines, expectedLines) {
		t.FailNow()
	}
	if !bytes.Equal(resultRest, expectedRest) {
		t.FailNow()
	}
}

func Test_BytesHead(t *testing.T) {
	bytesHeadTailTestHelper(
		t, BytesHead,
		[]byte("line1\nline2\r\nline3\nline4\nline5"), 3,
		[]string{"line1", "line2", "line3"}, []byte("line4\nline5"))
	bytesHeadTailTestHelper(
		t, BytesHead,
		[]byte("line1\nline2\r\nline3\nline4\nline5"), 6,
		[]string{"line1", "line2", "line3", "line4", "line5"}, []byte(""))
}

func Test_BytesTail(t *testing.T) {
	bytesHeadTailTestHelper(
		t, BytesTail,
		[]byte("line1\nline2\nline3\nline4\r\nline5"), 2,
		[]string{"line5", "line4"}, []byte("line1\nline2\nline3"))
	bytesHeadTailTestHelper(
		t, BytesTail,
		[]byte("line1\nline2\r\nline3\nline4\nline5"), 6,
		[]string{"line5", "line4", "line3", "line2", "line1"}, []byte(""))
}

func Test_BytesMap(t *testing.T) {
	upper := func(b byte) byte {
		return b - ('a' - 'A')
	}
	result := BytesMap(upper, []byte("hello"))
	correct := []byte("HELLO")
	if len(result) != len(correct) {
		t.Fail()
	}
	for i := range result {
		if result[i] != correct[i] {
			t.Fail()
		}
	}
}

func Test_BytesFilter(t *testing.T) {
	azFunc := func(b byte) bool {
		return b >= 'A' && b <= 'Z'
	}
	result := BytesFilter(azFunc, []byte{1, 2, 3, 'A', 'f', 'R', 123})
	correct := []byte{'A', 'R'}
	if len(result) != len(correct) {
		t.Fail()
	}
	for i := range result {
		if result[i] != correct[i] {
			t.Fail()
		}
	}
}
