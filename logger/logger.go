package logger

import (
	"go.uber.org/zap"
	"log"
	"os"
)

type Logger struct {
	zapLogger *zap.Logger
	baseLog   *log.Logger
}

// New creates a new Logger Agent
func NewLogAgent(zapLogger *zap.Logger) *Logger {
	return &Logger{zapLogger: zapLogger, baseLog: log.New(os.Stderr, "", log.LstdFlags)}
}

func (log *Logger) Debug(msg string) {
	if nil != log.zapLogger {
		log.zapLogger.Debug(msg)
	} else {
		log.baseLog.Println(msg)
	}
}

func (log *Logger) Info(msg string) {
	if nil != log.zapLogger {
		log.zapLogger.Info(msg)
	} else {
		log.baseLog.Println(msg)
	}
}

func (log *Logger) Warn(msg string) {
	if nil != log.zapLogger {
		log.zapLogger.Warn(msg)
	} else {
		log.baseLog.Println(msg)
	}
}

func (log *Logger) Error(msg string) {
	if nil != log.zapLogger {
		log.zapLogger.Error(msg)
	} else {
		log.baseLog.Println(msg)
	}
}

func (log *Logger) Panic(msg string) {
	if nil != log.zapLogger {
		log.zapLogger.Panic(msg)
	} else {
		log.baseLog.Panic(msg)
	}
}

func (log *Logger) DPanic(msg string) {
	if nil != log.zapLogger {
		log.zapLogger.DPanic(msg)
	} else {
		log.baseLog.Println(msg)
	}
}

func (log *Logger) Fatal(msg string) {
	if nil != log.zapLogger {
		log.zapLogger.Fatal(msg)
	} else {
		log.baseLog.Fatal(msg)
	}
}
