package llog

import (
	"log"
	"testing"
)

func TestStdLog(t *testing.T) {
	logger := NewStdLogger(
		WithStdColored(true),
		// WithStdColored(false),
		WithStdWriter(log.Writer()),
	)
	logger.Log(Info, "Infokey", "Infoval")
	logger.Log(Debug, "Debugkey", "Debugval")
	logger.Log(Warn, "Warnkey", "Warnval")
	logger.Log(Error, "Errorkey", "Errorval")
	logger.Log(Fatal, "Fatalkey", "Fatalval")
	// assert.Equal(t, "hello", "hello")
}
