package llog

import (
	"fmt"
	"os"
)

type HelperOption interface {
	apply(*Helper) error
}

type helperOptionFunc func(*Helper) error

func (f helperOptionFunc) apply(h *Helper) error {
	return f(h)
}

const DefaultMessageKey = "msg"

// from kratos
type Helper struct {
	logger  Logger
	msgKey  string
	sprint  func(...interface{}) string
	sprintf func(format string, a ...any) string
}

func WithMessageKey(msgkey string) HelperOption {
	return helperOptionFunc(func(h *Helper) error {
		h.msgKey = msgkey
		return nil
	})
}

func WithSprint(sp func(...interface{}) string) HelperOption {
	return helperOptionFunc(func(h *Helper) error {
		h.sprint = sp
		return nil
	})
}

func WithSprintf(spf func(format string, a ...any) string) HelperOption {
	return helperOptionFunc(func(h *Helper) error {
		h.sprintf = spf
		return nil
	})
}

func NewHelper(logger Logger, opts ...HelperOption) *Helper {
	options := &Helper{
		msgKey:  DefaultMessageKey,
		logger:  logger.Clone().AddCallerSkip(1),
		sprint:  fmt.Sprint,
		sprintf: fmt.Sprintf,
	}
	for _, o := range opts {
		if err := o.apply(options); err != nil {
			// TODO: panic ???
			panic(err)
		}
	}
	return options
}

func (h *Helper) Log(level Level, keyvals ...any) {
	_ = h.logger.Log(level, keyvals...)
}

func (h *Helper) Debug(a ...any) {
	_ = h.logger.Log(Debug, h.msgKey, h.sprint(a...))
}

func (h *Helper) Debugf(format string, a ...any) {
	_ = h.logger.Log(Debug, h.msgKey, h.sprintf(format, a...))
}

func (h *Helper) Debugw(kvs ...any) {
	_ = h.logger.Log(Debug, kvs...)
}

func (h *Helper) Info(a ...any) {
	_ = h.logger.Log(Info, h.msgKey, h.sprint(a...))
}

func (h *Helper) Infof(format string, a ...any) {
	_ = h.logger.Log(Info, h.msgKey, h.sprintf(format, a...))
}

func (h *Helper) Infow(kvs ...any) {
	_ = h.logger.Log(Info, kvs...)
}

func (h *Helper) Warn(a ...any) {
	_ = h.logger.Log(Warn, h.msgKey, h.sprint(a...))
}

func (h *Helper) Warnf(format string, a ...any) {
	_ = h.logger.Log(Warn, h.msgKey, h.sprintf(format, a...))
}

func (h *Helper) Warnw(kvs ...any) {
	_ = h.logger.Log(Warn, kvs...)
}

func (h *Helper) Error(a ...any) {
	_ = h.logger.Log(Error, h.msgKey, h.sprint(a...))
}

func (h *Helper) Errorf(format string, a ...any) {
	_ = h.logger.Log(Error, h.msgKey, h.sprintf(format, a...))
}

func (h *Helper) Errorw(kvs ...any) {
	_ = h.logger.Log(Error, kvs...)
}

func (h *Helper) Fatal(a ...any) {
	_ = h.logger.Log(Fatal, h.msgKey, h.sprint(a...))
	os.Exit(1)
}

func (h *Helper) Fatalf(format string, a ...any) {
	_ = h.logger.Log(Fatal, h.msgKey, h.sprintf(format, a...))
	os.Exit(1)
}

func (h *Helper) Fatalw(kvs ...any) {
	_ = h.logger.Log(Error, kvs...)
	os.Exit(1)
}
