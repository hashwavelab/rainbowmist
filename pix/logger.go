package pix

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	BASE_PATH = "logs/"
)

func NewLogger(path string) *zap.Logger {
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   BASE_PATH + path,
		MaxSize:    10, // megabytes
		MaxBackups: 10,
		MaxAge:     1, // days
	})
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(config),
		w,
		zap.InfoLevel,
	)
	return zap.New(core)
}
