package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// A Level is a logging priority. Higher levels are more important.
type LogLevel int8

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel LogLevel = iota - 1
	// InfoLevel is the default logging priority.
	InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel
	// DPanicLevel logs are particularly important errors. In development the
	// logger panics after writing the message.
	DPanicLevel
	// PanicLevel logs a message, then panics.
	PanicLevel
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel
)

type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
	Rotate() error
}

type log struct {
	sugar *zap.SugaredLogger
	ljLog *lumberjack.Logger
}

func (l *log) Rotate() error {
	return l.ljLog.Rotate()
}

// Init :logger initialization
func Init(logPath string, loglevel LogLevel) (Logger, error) {
	err := os.MkdirAll(filepath.Dir(logPath), os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("error creating log directory: %v", err)
	}
	ljLog := &lumberjack.Logger{
		Filename: logPath,
	}
	w := zapcore.AddSync(ljLog)
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()),
		w,
		zapcore.Level(loglevel),
	)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	defer logger.Sync()
	return &log{
		sugar: logger.Sugar(),
		ljLog: ljLog,
	}, nil
}

// Info : info sugar logging
func (l *log) Info(args ...interface{}) {
	l.sugar.Info(args)
}

// Error : error sugar logging
func (l *log) Error(args ...interface{}) {
	l.sugar.Error(args)
}
