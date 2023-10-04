package llog

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

var _ Logger = (*stdLogger)(nil)

type StdOption interface {
	apply(*StdConfig) error
}

var colorsLen = 6

type StdConfig struct {
	Writer io.Writer
	// Info Debug Warn Error Fatal
	Colors     map[Level][]Color
	Colored    bool
	TimeFormat string
	TimeZone   *time.Location
}

type stdopfunc func(*StdConfig) error

func (s stdopfunc) apply(sc *StdConfig) error {
	return s(sc)
}

type stdLogger struct {
	log    *log.Logger
	pool   *sync.Pool
	stdcfg *StdConfig
}

func WithStdTimeFormat(format string) StdOption {
	return stdopfunc(func(sc *StdConfig) error {
		if format == "" {
			return errors.New("time format empty")
		}
		sc.TimeFormat = format
		return nil
	})
}

func WithStdTimeZone(tz *time.Location) StdOption {
	return stdopfunc(func(sc *StdConfig) error {
		if tz == nil {
			return errors.New("time zone nil")
		}
		sc.TimeZone = tz
		return nil
	})
}

func WithStdWriter(w io.Writer) StdOption {
	return stdopfunc(func(sc *StdConfig) error {
		sc.Writer = w
		return nil
	})
}

func WithStdColors(cs map[Level][]Color) StdOption {
	return stdopfunc(func(sc *StdConfig) error {
		// cs's length should be le then 6
		if len(cs) > colorsLen {
			return errors.New("colors should be no more than 6")
		}
		for level, colour := range cs {
			if len(colour) != len(sc.Colors[level]) {
				return errors.New(fmt.Sprintf("two much color attributes for: %d", level))
			}
			copy(sc.Colors[level], colour)
		}
		return nil
	})
}

func WithStdColored(c bool) StdOption {
	return stdopfunc(func(sc *StdConfig) error {
		sc.Colored = c
		return nil
	})
}

// TODO: lock
func NewStdLogger(opts ...StdOption) Logger {
	cfg := &StdConfig{
		Colored: true,
		Writer:  os.Stdin,
		Colors: map[Level][]Color{
			Info: {
				FgWhite,
				BgWhite,
			},
			Debug: {
				FgWhite,
				BgCyan,
			},
			Warn: {
				FgWhite,
				BgYellow,
			},
			Error: {
				FgWhite,
				BgRed,
			},
			Fatal: {
				BgRed,
				FgBlack,
			},
			// Fatal
		},
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   time.Now().UTC().Location(),
	}
	l := &stdLogger{
		pool: &sync.Pool{
			New: func() any {
				return new(bytes.Buffer)
			},
		},
	}
	for _, o := range opts {
		o.apply(cfg)
	}
	l.log = log.New(cfg.Writer, "", 0)
	l.stdcfg = cfg
	return l
}

func (l *stdLogger) Log(level Level, keyvals ...any) error {
	if len(keyvals) == 0 {
		return nil
	}
	if (len(keyvals) & 1) == 1 {
		keyvals = append(keyvals, "KEYVALS UNPAIRED")
	}
	buf := l.pool.Get().(*bytes.Buffer)
	var ws string
	if !l.stdcfg.Colored {
		ws = level.String()
		goto blank
	}
	switch level {
	case Debug:
		ws = WithColor(level.String(), l.stdcfg.Colors[Debug]...)
	case Info:
		ws = WithColor(level.String(), l.stdcfg.Colors[Info]...)
	case Warn:
		ws = WithColor(level.String(), l.stdcfg.Colors[Warn]...)
	case Error:
		ws = WithColor(level.String(), l.stdcfg.Colors[Error]...)
	case Fatal:
		ws = WithColor(level.String(), l.stdcfg.Colors[Fatal]...)
	default:
		ws = WithColor(level.String(), l.stdcfg.Colors[Info]...)
	}
blank:
	buf.WriteString(ws)
	fmt.Fprintf(buf, " time: %v", time.Now().Format(l.stdcfg.TimeFormat))
	// TODO: maybe this should be colored？
	for i := 0; i < len(keyvals); i += 2 {
		_, _ = fmt.Fprintf(buf, " %s: %v", keyvals[i], keyvals[i+1])
	}
	_ = l.log.Output(4, buf.String())
	buf.Reset()
	l.pool.Put(buf)
	return nil
}
