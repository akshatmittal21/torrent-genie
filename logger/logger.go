package logger

import (
	"os"
	"path/filepath"

	"github.com/akshatmittal21/torrent-genie/constants"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var sugar *zap.SugaredLogger
var ljLog *lumberjack.Logger

func init() {
	InitLogger(constants.LogPath, zapcore.InfoLevel)
}

func Rotate() {
	ljLog.Rotate()
}

// InitLogger :logger initialisation
func InitLogger(logPath string, loglevel zapcore.Level) {

	os.MkdirAll(filepath.Dir(logPath), os.ModePerm)
	ljLog = &lumberjack.Logger{
		Filename: logPath,
	}
	w := zapcore.AddSync(ljLog)
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()),
		w,
		loglevel,
	)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	defer logger.Sync()
	sugar = logger.Sugar()
}

// Info : info sugar logging
func Info(args ...interface{}) {
	sugar.Info(args)
}

// Error : error sugar logging
func Error(args ...interface{}) {
	sugar.Error(args)
}
