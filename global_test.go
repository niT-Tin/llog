package llog

import (
	"testing"
)

func TestGLogger(t *testing.T) {
	Glogger.Log(Warn, "password", "global val")
	// stdlog := NewStdLogger(WithStdWriter(os.Stdout))
	stdlog := NewStdLogger()
	SetGLogger(NewFilter(
		stdlog,
		WithFilterKeys("password"),
	))
	Glogger.Log(Warn, "password", "global val")
	// stdlog.Log(Warn, "password", "global val")
}
