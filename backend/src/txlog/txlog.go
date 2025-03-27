package txlog

import (
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

const (
	LOG_LEVEL_DEBUG = iota
	LOG_LEVEL_INFO
	LOG_LEVEL_WARN
	LOG_LEVEL_ERROR
)

func logLevelValid(level int) bool {
    return LOG_LEVEL_DEBUG <= level && level <= LOG_LEVEL_ERROR
}

var txLogger *log.Logger
var logLevel int = LOG_LEVEL_INFO // Default

func InitLogging(level int, logFilepath string) error {

	file, err := os.OpenFile(logFilepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
        return errors.New(fmt.Sprintf("Failed to open log file: %s\n", logFilepath))
	}

	txLogger = log.New(file, "", 0)

    if (!logLevelValid(level)) {
        TxLogError("Invalid log level provided: defaulting to INFO: level=%d", level)
        logLevel = LOG_LEVEL_INFO
    } else {
        logLevel = level 
    }

    return nil
}

// Helper function to get the file and line number
func getCallerInfo() (string, int) {
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		return "unknown", 0
	}
	idx := strings.LastIndex(file, "/")
	file = file[idx+1:]
	return file, line
}

func Str2LogLevel(level string) int {
    level = strings.ToUpper(level)
    if level == "DEBUG" {
        return LOG_LEVEL_DEBUG
    }
    if level == "INFO" {
        return LOG_LEVEL_INFO
    }
    if level == "WARN" {
        return LOG_LEVEL_WARN
    }
    if level == "ERROR" {
        return LOG_LEVEL_ERROR
    }
    return -1
} 

func LogLevel2Str(level int) string {
	switch level {
	case LOG_LEVEL_INFO:
		return "INFO"
	case LOG_LEVEL_WARN:
		return "WARN"
	case LOG_LEVEL_ERROR:
		return "ERROR"
	case LOG_LEVEL_DEBUG:
		return "DEBUG"
	default:
		return ""
	}
}

func txLog(level int, format string, args ...any) {
	if level < logLevel || txLogger == nil {
		return
	}
	now := time.Now().Format("2006-01-02 15:04:05.000")
	file, line := getCallerInfo()
    txLogger.Printf("%s", fmt.Sprintf("%s %-5s [%15s:%4d] "+format, append([]any{now, LogLevel2Str(level), file, line}, args...)...))
}

func TxLogInfo(format string, args ...any) {
	txLog(LOG_LEVEL_INFO, format, args...)
}

func TxLogWarn(format string, args ...any) {
	txLog(LOG_LEVEL_WARN, format, args...)
}

func TxLogError(format string, args ...any) {
	txLog(LOG_LEVEL_ERROR, format, args...)
}

func TxLogDebug(format string, args ...any) {
	txLog(LOG_LEVEL_DEBUG, format, args...)
}
