package dry

// see ./test.go doc how it works
var _ = func() bool {
	testMode = true
	return true
}()
