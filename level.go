package llog

import "strings"

type Level int8

const (
	Debug Level = iota + 1
	Info
	Warn
	Error
	Fatal
)

func (l Level) String() string {
	switch l {
	case Debug:
		return "DEBU"
	case Info:
		return "INFO"
	case Warn:
		return "WARN"
	case Error:
		return "ERRO"
	case Fatal:
		return "FATA"
	default:
		return ""
	}
}

// can I generate this?

func ParseLevel(s string) Level {
	switch strings.ToUpper(s) {
	case "DEBU":
		return Debug
	case "INFO":
		return Info
	case "WARN":
		return Warn
	case "ERRO":
		return Error
	case "FATA":
		return Fatal
	}
	return Info
}
