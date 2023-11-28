package llog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelper(t *testing.T) {

	helper := NewHelper(
		NewFilter(
			// NewStdLogger(WithStdWriter(os.Stdout)),
			NewStdLogger(),
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

func TestClone(t *testing.T) {
	// stdlog := NewStdLogger(WithStdWriter(os.Stdout))
	stdlog := NewStdLogger()
	filter_1 := NewFilter(stdlog)
	filter_2 := filter_1.Clone().(*Filter)
	filter_1.AddCallerSkip(1)
	assert.NotEqual(t, filter_1.logger.GetCallerSkip(), filter_2.logger.GetCallerSkip())
	// filter_1.key["password"] = struct{}{}
	// assert.NotEqual(t, len(filter_1.key), len(filter_2.key))
}

func TestHelperClone(t *testing.T) {
	// stdlog := NewStdLogger(WithStdWriter(os.Stdout))
	stdlog := NewStdLogger()
	filter_1 := NewFilter(stdlog)
	helper_1 := NewHelper(filter_1)
	filter_2 := filter_1.Clone().(*Filter)
	helper_2 := NewHelper(filter_2)

	filter_1.Log(Error, "filter_1 Warn", "world")
	helper_1.Debug("helper_1 Debug", "world")
	helper_2.Info("helper_2 Info", "world")
	filter_2.Log(Warn, "filter_2 Warn", "world")
	stdlog.Log(Error, "stdlog Warn", "world")
}

func TestHelperWrite(t *testing.T) {
	logger := NewStdLogger(
		WithBackUp(true),
		WithLogFile("file.log"),
		WithMaxLogSize(1*Kb),
		WithStdOut(false),
	)
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
