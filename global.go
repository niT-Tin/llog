package llog

import "sync"

var (
	lock    sync.RWMutex
	Glogger Logger
	Ghelper *Helper
)

func init() {
	Glogger = DefaultLogger
	Ghelper = NewHelper(Glogger.Clone())
}

func SetGLogger(l Logger) {
	lock.Lock()
	defer lock.Unlock()
	Glogger = l
}

func SetGHelper(h *Helper) {
	lock.Lock()
	defer lock.Unlock()
	Ghelper = h
}
