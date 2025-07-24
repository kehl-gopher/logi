package utils

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	lg "gopkg.in/natefinch/lumberjack.v2"
)

type LogLevel int

const (
	InfoLevel LogLevel = iota + 1
	DebugLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

type fLogger interface {
	fatalLog(message string, args ...interface{})
}

type Log struct {
	lmutex *sync.Mutex
	log    *slog.Logger
	fLog   fLogger
}

type fLog struct {
	log *slog.Logger
}

func newFlogger(l *lg.Logger) fLogger {
	faLog := slog.New(slog.NewJSONHandler(l, &slog.HandlerOptions{
		Level: slog.LevelError,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey && a.Value.String() == slog.LevelError.String() {
				return slog.String(a.Key, "FATAL")
			}
			return a
		},
	}))
	return &fLog{log: faLog}
}
func NewLogger() *Log {
	dir, _ := BasePath()
	LOG_DIR := filepath.Join(dir, "logs", "app.log")
	lumLog := &lg.Logger{
		Filename:   LOG_DIR,
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     30,
	}
	handler := slog.New(slog.NewJSONHandler(lumLog, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	faLog := newFlogger(lumLog)
	return &Log{log: handler, fLog: faLog, lmutex: new(sync.Mutex)}
}

func (l *Log) infoLog(message string, args ...interface{}) {
	l.lmutex.Lock()
	defer l.lmutex.Unlock()
	if args == nil {
		l.log.Info(message)
		return
	}
	l.log.Info(message, args...)
}

func (l *Log) debugLog(message string, args ...interface{}) {
	l.lmutex.Lock()
	defer l.lmutex.Unlock()
	if args == nil {
		l.log.Debug(message)
		return
	}
	l.log.Info(message)
}

func (l *Log) warnLog(message string, args ...interface{}) {
	l.lmutex.Lock()
	defer l.lmutex.Unlock()
	if args == nil {
		l.log.Warn(message)
		return
	}
	l.log.Warn(message, args...)
}

func (l *Log) errorLog(message string, args ...interface{}) {
	l.lmutex.Lock()
	defer l.lmutex.Unlock()
	if args == nil {
		l.log.Error(message)
		return
	}
	l.log.Error(message, args...)
}

func (l *fLog) fatalLog(message string, args ...interface{}) {
	if args == nil {
		l.log.Error(message)
	}
	l.log.Error(message, args...)
}

func PrintLog(log *Log, message string, level LogLevel, args ...interface{}) {
	switch level {
	case 1:
		log.infoLog(message, args...)
	case 2:
		log.debugLog(message, args...)
	case 3:
		log.warnLog(message, args...)
	case 4:
		log.errorLog(message, args...)
	case 5:
		log.lmutex.Lock()
		log.fLog.fatalLog(message, args)
		log.lmutex.Unlock()
	}

	buf := bufio.NewWriter(os.Stdout)
	log.lmutex.Lock()
	_, _ = buf.Write([]byte(message))
	fmt.Fprintln(buf, args...)
	buf.Flush()
	log.lmutex.Unlock()
}
