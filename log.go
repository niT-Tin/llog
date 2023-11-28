package llog

var DefaultLogger = NewStdLogger()

type Logger interface {
	Log(l Level, keyvals ...any) error
	AddCallerSkip(skip int) Logger
	GetCallerSkip() int
	Clone() Logger
}

// TODO: maybe later
// type logger struct {
// 	logger Logger
// 	ctx    context.Context
// }

// func (l *logger) Log(lvl Level, keyvals ...any) error {
// 	return l.logger.Log(lvl, keyvals...)
// }
