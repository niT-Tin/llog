package llog

import (
	"os"
	"testing"
)

func TestHelper(t *testing.T) {

	helper := NewHelper(
		NewFilter(
			DefaultLogger,
			WithLevel(Warn),
			WithFilterKeys("password", "mobile"),
		),
	)
	helper.Debug("debug info should not be seen")
	helper.Warn("warn info")
	helper.Error("error info")
	helper.Errorf("error %s, %d", "error info", 12)
	helper.Errorw("error", "value")
	// helper.Fatal("fatal value")
}

func TestHelperWrite(t *testing.T) {
	file, err := os.OpenFile("file.log", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	logger := NewStdLogger(WithStdWriter(file))
	helper := NewHelper(
		NewFilter(
			logger,
			// WithLevel(Warn),
			WithFilterKeys("password", "mobile"),
		),
	)
	// password should be masked
	helper.Debugw("password", "122145jasdfxd", "name", "bruce wayne")

	// password should NOT be masked
	helper.Debug("password", "122145jasdfxd", "name", "bruce wayne")
	helper.Infof("format: %s %s %s %s", "password", "122145jasdfxd", "name", "bruce wayne")
}
