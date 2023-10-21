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
		WithStdTimeFormat("2006-01-02 15:04"),
	)
	logger.Log(Info, "STD: Infokey", "Infoval")
	logger.Log(Debug, "STD: Debugkey", "Debugval")
	logger.Log(Warn, "STD: Warnkey", "Warnval")
	logger.Log(Error, "STD: Errorkey", "Errorval")
	// logger.Log(Fatal, "STD: Fatalkey", "Fatalval")
	// assert.Equal(t, "hello", "hello")

	colorMap := make(map[Level][]Color, 1)
	colorMap[Info] = []Color{FgYellow, BgGreen}
	newLogger := NewStdLogger(
		WithStdColors(colorMap),
	)
	nextNewLogger := NewStdLogger(
	// WithStdTimeFormat("2006-01-02 15:04"),
	)
	newLogger.Log(Info, "STD: Infokey", "Infoval")
	nextNewLogger.Log(Info, "STD: Infokey", "Infoval")
}
