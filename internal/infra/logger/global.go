package logger

import "sync"

var (
	_globalLog BaseLogger = newNop()
	_globalMu  sync.RWMutex
)

func GetLog() BaseLogger {
	_globalMu.RLock()
	defer _globalMu.RUnlock()

	log := _globalLog

	return log
}

func ReplaceGlobal(log BaseLogger) {
	_globalMu.Lock()
	defer _globalMu.Unlock()

	_globalLog = log
}
