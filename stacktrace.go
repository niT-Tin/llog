package llog

import (
	"runtime"
	"sync"
)

type stacktrace struct {
	pcs     []uintptr
	frames  *runtime.Frames
	storage []uintptr
}

var _stacktracePool = &sync.Pool{
	New: func() interface{} {
		return &stacktrace{
			storage: make([]uintptr, 1),
		}
	},
}

func (s *stacktrace) Next() (runtime.Frame, bool) {
	return s.frames.Next()
}

func (s *stacktrace) Free() {
	s.pcs = nil
	_stacktracePool.Put(s)
}

func captureStacktrace(skip int) *stacktrace {
	// skip captureStacktrace and Log also for Next()
	// 3 = captureStacktrace + Log + Next
	const frameOffset = 3
	stack := _stacktracePool.Get().(*stacktrace)
	// TODO: 1 is enough?
	stack.pcs = stack.storage[:1]
	numFrames := runtime.Callers(skip+frameOffset, stack.pcs)
	stack.frames = runtime.CallersFrames(stack.pcs[:numFrames])
	return stack
}
