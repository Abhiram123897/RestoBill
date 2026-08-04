package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log   *zap.Logger
	sugar *zap.SugaredLogger
)

func init() {
	var err error
	config := zap.NewProductionConfig()

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.StacktraceKey = ""
	config.EncoderConfig = encoderConfig

	log, err = config.Build(zap.AddCallerSkip(1))
	if err != nil {
		os.Exit(1)
	}

	sugar = log.Sugar()
}

func Info(message string, fields ...zap.Field) {
	log.Info(message, fields...)
}

func Debug(message string, fields ...zap.Field) {
	log.Debug(message, fields...)
}

func Error(message string, fields ...zap.Field) {
	log.Error(message, fields...)
}

func Fatal(message string, fields ...zap.Field) {
	log.Fatal(message, fields...)
}

func InfoF(message string, args ...interface{}) {
	sugar.Infof(message, args...)
}

func DebugF(message string, args ...interface{}) {
	sugar.Debugf(message, args...)
}

func ErrorF(message string, args ...interface{}) {
	sugar.Errorf(message, args...)
}

func FatalF(message string, args ...interface{}) {
	sugar.Fatalf(message, args...)
}
