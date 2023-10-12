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
		WithFilterValues("world"),
		WithFilterFunc(func(level Level, keyvals ...interface{}) bool {
			if level == Warn {
				for i := 0; i < len(keyvals); i += 2 {
					if keyvals[i] == "password" {
						keyvals[i+1] = "******"
						return true
					}
				}
			}
			return false
		}),
	)
	// nothing
	filter.Log(Info, "asdfasdfasdfasdf")
	filter.Log(Warn, "hello", "world")
	filter.Log(Warn, "password", "world")
	// filter.Log(Fatal, "mobile", "world")
}
