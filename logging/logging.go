package logging

import "fmt"

type LOG_LEVEL int

const (
	LOG_LEVEL_TRACE LOG_LEVEL = 100
	LOG_LEVEL_DEBUG LOG_LEVEL = 200
	LOG_LEVEL_INFO  LOG_LEVEL = 300
	LOG_LEVEL_WARN  LOG_LEVEL = 400
	LOG_LEVEL_ERROR LOG_LEVEL = 500
	LOG_LEVEL_FATAL LOG_LEVEL = 600
)

var GLOBAL_LOG_LEVEL = LOG_LEVEL_TRACE

func logf(level LOG_LEVEL, msg_format string, args ...interface{}) {
	if level <= GLOBAL_LOG_LEVEL {
		return
	}
	if len(args) == 0 {
		fmt.Println(msg_format)
	} else {
		fmt.Printf(msg_format, args...)
		fmt.Println()
	}
}

var DEBUG = func(msg_format string, args ...interface{}) {
	logf(LOG_LEVEL_DEBUG, msg_format, args...)
}

var TRACE = func(msg_format string, args ...interface{}) {
	logf(LOG_LEVEL_TRACE, msg_format, args...)
}
var INFO = func(msg_format string, args ...interface{}) {
	logf(LOG_LEVEL_INFO, msg_format, args...)
}
var WARN = func(msg_format string, args ...interface{}) {
	logf(LOG_LEVEL_WARN, msg_format, args...)
}

var ERROR = func(msg_format string, args ...interface{}) {
	logf(LOG_LEVEL_ERROR, msg_format, args...)
}

var FATAL = func(msg_format string, args ...interface{}) {
	logf(LOG_LEVEL_FATAL, msg_format, args...)
}

func SetLogLevel(level int) {
	GLOBAL_LOG_LEVEL = LOG_LEVEL(level)
}
