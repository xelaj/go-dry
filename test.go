package dry

import (
	"os"
	"strings"
)

func IsTestMode() bool {
	if len(os.Args) < 1 {
		return false
	}

	return strings.HasSuffix(os.Args[0], ".test")
}
