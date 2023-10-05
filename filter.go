package llog

import (
	"context"
)

type FilterOption interface {
	apply(*Filter)
}

type filterOptionfunc func(*Filter)

func (ff filterOptionfunc) apply(f *Filter) {
	ff(f)
}

const fuzzyStr = "***"

var _ Logger = (*Filter)(nil)

// some from kratos
type Filter struct {
	ctx    context.Context
	logger Logger
	level  Level
	// 对应key对应的value需要屏蔽
	key map[interface{}]struct{}
	// 对应的value需要屏蔽
	value  map[interface{}]struct{}
	filter func(level Level, keyvals ...interface{}) bool
}

func NewFilter(logger Logger, opts ...FilterOption) *Filter {
	f := &Filter{
		// 添加一层caller skip，Filter内调用
		logger: logger.Clone().AddCallerSkip(1),
		key:    make(map[interface{}]struct{}),
		value:  make(map[interface{}]struct{}),
	}
	for _, o := range opts {
		o.apply(f)
	}
	return f
}

func WithFilterKeys(keys ...string) FilterOption {
	return filterOptionfunc(func(f *Filter) {
		for _, k := range keys {
			f.key[k] = struct{}{}
		}
	})
}

func WithFilterVlaues(values ...string) FilterOption {
	return filterOptionfunc(func(f *Filter) {
		for _, v := range values {
			f.value[v] = struct{}{}
		}
	})
}

func WithFilterFunc(f func(level Level, keyvals ...interface{}) bool) FilterOption {
	return filterOptionfunc(func(o *Filter) {
		o.filter = f
	})
}

func WithLevel(level Level) FilterOption {
	return filterOptionfunc(func(o *Filter) {
		o.level = level
	},
	)
}

// Filter 中间层调用底部logger的caller skip
func (f *Filter) AddCallerSkip(skip int) Logger {
	f.logger.AddCallerSkip(skip)
	return f
}

func (f *Filter) GetCallerSkip() int {
	return f.logger.GetCallerSkip()
}

func (f *Filter) Clone() Logger {
	// cloned := &Filter{}
	cloend_logger := f.logger.Clone()
	cloned := *f
	cloned.logger = cloend_logger
	return &cloned
}

func (f *Filter) Log(l Level, keyvals ...any) error {
	if l < f.level {
		return nil
	}
	if f.filter != nil && (f.filter(l, keyvals...)) {
		return nil
	}
	if len(f.key) > 0 || len(f.value) > 0 {
		for i := 0; i < len(keyvals); i += 2 {
			v := i + 1
			if v >= len(keyvals) {
				continue
			}
			if _, ok := f.key[keyvals[i]]; ok {
				keyvals[v] = fuzzyStr
			}
			if _, ok := f.value[keyvals[v]]; ok {
				keyvals[v] = fuzzyStr
			}
		}
	}
	return f.logger.Log(l, keyvals...)
}
