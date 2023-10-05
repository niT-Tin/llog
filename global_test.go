package llog

import (
	"os"
	"testing"
)

func TestGLogger(t *testing.T) {
	Glogger.Log(Warn, "password", "global val")
	stdlog := NewStdLogger(WithStdWriter(os.Stdout))
	SetGLogger(NewFilter(
		stdlog,
		WithFilterKeys("password"),
	))
	Glogger.Log(Warn, "password", "global val")
	// stdlog.Log(Warn, "password", "global val")
}
