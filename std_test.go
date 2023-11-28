package llog

import (
	"testing"
)

func TestStdLog(t *testing.T) {
	// f, err := os.OpenFile("test.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0700)
	// defer func() {
	// 	f.Sync()
	// 	f.Close()
	// }()

	// if err != nil {
	// 	panic(err)
	// }
	logger := NewStdLogger(
		WithStdColored(true),
		WithBackUp(true),
		WithMaxLogSize(1*Kb),
		// WithBackUp(false),
		WithLogFile("/tmp/log/llog.log"),
		// WithStdColored(false),
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
