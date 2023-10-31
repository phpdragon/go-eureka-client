package logger

import (
	"go.uber.org/zap"
	"log"
	"os"
)

type Logger struct {
	zapLogger *zap.SugaredLogger
	baseLog   *log.Logger
}

// New creates a new Logger Agent
func NewLogAgent(zapLogger *zap.SugaredLogger) *Logger {
	return &Logger{zapLogger: zapLogger, baseLog: log.New(os.Stderr, "", log.LstdFlags)}
}

func (log *Logger) Debug(args ...interface{}) {
	if nil != log.zapLogger {
		log.zapLogger.Debug(args)
	} else {
		log.baseLog.Println(args)
	}
}

func (log *Logger) Info(args ...interface{}) {
	if nil != log.zapLogger {
		log.zapLogger.Info(args)
	} else {
		log.baseLog.Println(args)
	}
}

func (log *Logger) Warn(args ...interface{}) {
	if nil != log.zapLogger {
		log.zapLogger.Warn(args)
	} else {
		log.baseLog.Println(args)
	}
}

func (log *Logger) Error(args ...interface{}) {
	if nil != log.zapLogger {
		log.zapLogger.Error(args)
	} else {
		log.baseLog.Println(args)
	}
}

func (log *Logger) Panic(args ...interface{}) {
	if nil != log.zapLogger {
		log.zapLogger.Panic(args)
	} else {
		log.baseLog.Panic(args)
	}
}

func (log *Logger) DPanic(args ...interface{}) {
	if nil != log.zapLogger {
		log.zapLogger.DPanic(args)
	} else {
		log.baseLog.Println(args)
	}
}

func (log *Logger) Fatal(args ...interface{}) {
	if nil != log.zapLogger {
		log.zapLogger.Fatal(args)
	} else {
		log.baseLog.Fatal(args)
	}
}
