package llog

import (
	"log"
	"testing"
)

func TestFilter(t *testing.T) {
	logger := NewStdLogger(
		WithStdWriter(log.Writer()),
	)
	filter := NewFilter(
		logger,
		WithFilterKeys("password", "mobile"),
		WithLevel(Warn),
	)
	// nothing
	filter.Log(Info, "asdfasdfasdfasdf")
	filter.Log(Warn, "hello", "world")
	filter.Log(Warn, "password", "world")
	filter.Log(Fatal, "mobile", "world")
}
