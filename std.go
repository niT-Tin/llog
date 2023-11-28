package llog

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

var _ Logger = (*stdLogger)(nil)

type StdOption interface {
	apply(*StdConfig) error
}

var colorsLen = 6

const (
	b = 1 << (10 * iota) // byte
	Kb
	Mb
	Gb
	Tb
	perm = 0700
	// Pb
)

type StdConfig struct {
	// 日志输出基础文件名
	LogFile string
	// 是否在达到指定大小时切换文件
	IsBackup bool
	// 日志文件大小
	MaxLogSize int64
	// 是否启用标准输出
	IsStdOut bool
	// 是否启用标准错误
	IsStdErr bool
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
	pool       *sync.Pool
	callerSkip int
	stdcfg     *StdConfig
	logfile    *os.File
	// io         *IO
}

func WithStdOut(c bool) StdOption {
	return stdopfunc(func(sc *StdConfig) error {
		sc.IsStdOut = c
		return nil
	})
}

func WithStdErr(c bool) StdOption {
	return stdopfunc(func(sc *StdConfig) error {
		sc.IsStdErr = c
		return nil
	})
}
func WithMaxLogSize(sz int64) StdOption {
	return stdopfunc(func(sc *StdConfig) error {
		sc.MaxLogSize = sz
		return nil
	})
}

func WithBackUp(c bool) StdOption {
	return stdopfunc(func(sc *StdConfig) error {
		sc.IsBackup = c
		return nil
	})
}

func WithLogFile(f string) StdOption {
	return stdopfunc(func(sc *StdConfig) error {
		sc.LogFile = f
		return nil
	})
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

func WithStdColors(cs map[Level][]Color) StdOption {
	return stdopfunc(func(sc *StdConfig) error {
		// cs's length should be le then 6
		if len(cs) > colorsLen {
			return errors.New("colors should be no more than 6")
		}
		for level, colour := range cs {
			sc.Colors[level] = colour
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
		Colors: map[Level][]Color{
			Info: {
				// FgWhite,
				FgGreen,
			},
			Debug: {
				FgCyan,
				// BgCyan,
			},
			Warn: {
				// FgWhite,
				FgYellow,
			},
			Error: {
				// FgWhite,
				FgRed,
			},
			Fatal: {
				// BgRed,
				FgRed,
			},
			// Fatal
		},
		TimeFormat: "2006-01-02 15:04:05.000",
		TimeZone:   time.Now().UTC().Location(),
		IsBackup:   false,    // 默认在文件到达MaxLogSize大小时要切换文件
		MaxLogSize: 100 * Mb, // 默认文件大小为100日志mb
		IsStdOut:   true,     // 默认开启标准输出
		IsStdErr:   false,    // 默认关闭标准错误
		LogFile:    "llog.log",
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
	if !cfg.IsBackup {
		l.stdcfg = cfg
		return l
	}
	f := OpenFile(cfg.LogFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, perm)
	if f == nil {
		l.stdcfg = cfg
		return l
	}
	l.logfile = f
	l.stdcfg = cfg
	return l
}

// 底层logger的caller skip
func (l *stdLogger) AddCallerSkip(skip int) Logger {
	l.callerSkip += skip
	return l
}

func (l *stdLogger) GetCallerSkip() int {
	return l.callerSkip
}

func (l *stdLogger) Clone() Logger {
	cloned := &stdLogger{}
	*cloned = *l
	return cloned
}

func buildstring(l *stdLogger, level Level, keyvals ...any) (builded string, raw string) {
	if len(keyvals) == 0 {
		return "", ""
	}
	if (len(keyvals) & 1) == 1 {
		keyvals = append(keyvals, "KEYVALS UNPAIRED")
	}
	buf := l.pool.Get().(*bytes.Buffer)
	stack := captureStacktrace(l.callerSkip + 1)
	defer stack.Free()
	frame, _ := stack.Next()
	file, line := frame.File, frame.Line
	var ws string
	if !l.stdcfg.Colored {
		ws = level.String()
		ws = fmt.Sprintf("%v %9s", time.Now().Format(l.stdcfg.TimeFormat), ws)
		// 构建没有颜色的字符串前部份
		raw += ws
		goto blank
	}
	// 构建有颜色的字符串前部份
	builded = fmt.Sprintf("%v %20s", time.Now().Format(l.stdcfg.TimeFormat), level.String())
	raw = fmt.Sprintf("%v %20s", time.Now().Format(l.stdcfg.TimeFormat), level.String())
	switch level {
	case Debug:
		ws = WithColor(builded, l.stdcfg.Colors[Debug]...)
	case Info:
		ws = WithColor(builded, l.stdcfg.Colors[Info]...)
	case Warn:
		ws = WithColor(builded, l.stdcfg.Colors[Warn]...)
	case Error:
		ws = WithColor(builded, l.stdcfg.Colors[Error]...)
	case Fatal:
		ws = WithColor(builded, l.stdcfg.Colors[Fatal]...)
	default:
		ws = WithColor(builded, l.stdcfg.Colors[Info]...)
	}
blank:
	buf.WriteString(ws)
	var path string
	idx := strings.LastIndexByte(file, '/')
	if idx == -1 {
		path = file + ":" + fmt.Sprintf("%d", line)
	} else {
		idx = strings.LastIndexByte(file[:idx], '/')
		path = file[idx+1:] + ":" + fmt.Sprintf("%d", line)
	}
	fmt.Fprintf(buf, " [%s]", path)
	// TODO: maybe this should be colored？
	for i := 0; i < len(keyvals); i += 2 {
		raw += fmt.Sprintf(" %s: %v", keyvals[i], keyvals[i+1])
		_, _ = fmt.Fprintf(buf, " %s: %v", keyvals[i], keyvals[i+1])
	}
	res := buf.String()
	buf.Reset()
	l.pool.Put(buf)
	return res, raw
}

func (l *stdLogger) backup() bool {
	if !l.stdcfg.IsBackup {
		return true
	}
	fi, err := l.logfile.Stat()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return false
	}
	// 如果当前文件大小已经达到设置的备份大小，则进行备份
	if fi.Size() >= l.stdcfg.MaxLogSize {
		// 重新创建并打开一个新文件
		f := OpenFile(fmt.Sprintf("%s-%d.log", l.stdcfg.LogFile, time.Now().UnixNano()), os.O_CREATE|os.O_APPEND|os.O_RDWR, perm)
		if f == nil {
			fmt.Fprintln(os.Stderr, "new log file created failed")
			return false
		}
		// 关闭当前日志文件
		l.logfile.Close()
		l.logfile = f
		return true
	}
	return true
}

func (l *stdLogger) Log(level Level, keyvals ...any) error {
	s, raw := buildstring(l, level, keyvals...)
	if l.stdcfg.IsStdOut {
		// TODO: 暂时stdout直接使用fmt
		fmt.Fprintln(os.Stdout, s)
	}
	if l.stdcfg.IsStdErr {
		fmt.Fprintln(os.Stderr, s)
	}
	if l.stdcfg.IsBackup {
		if !l.backup() {
			return errors.New("back to file error")
		}
		_, err := l.logfile.Write([]byte(raw + "\n"))
		if err != nil {
			fmt.Fprintln(os.Stderr, "log to file error")
			return errors.New("log to file error")
		}
	}
	// _ = l.log.Output(4, buf.String())
	if level == Fatal {
		os.Exit(1)
	}
	return nil
}
