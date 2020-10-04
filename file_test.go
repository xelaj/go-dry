package dry

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_FileGetString(t *testing.T) {
	const hello = "Hello World!"

	t.Run("wrong file", func(t *testing.T) {
		_, err := FileGetString("invalid_file")
		require.Error(t, err)
	})

	t.Run("new file", func(t *testing.T) {
		tempFile := filepath.Join(t.TempDir(), "testfile.txt")
		err := ioutil.WriteFile(tempFile, []byte(hello), 0600)
		require.NoError(t, err)

		str, err := FileGetString(tempFile)
		require.NoError(t, err)
		require.Equal(t, hello, str)
	})

	t.Run("http", func(t *testing.T) {
		addr := "https://raw.githubusercontent.com/xelaj/go-dry/master/file_test.go"

		req, err := http.NewRequest("GET", addr, nil)
		require.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		raw, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err)

		str, err := FileGetString("https://raw.githubusercontent.com/xelaj/go-dry/master/file_test.go")
		require.NoError(t, err)

		require.Equal(t, string(raw), str)
	})
}

func Test_FileIsDir(t *testing.T) {
	if FileIsDir("testfile.txt") {
		t.Fail()
	}
	if !FileIsDir(".") {
		t.Fail()
	}
}
