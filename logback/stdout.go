package logback

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Stdout() Logger {
	prod := zap.NewProductionEncoderConfig()
	prod.EncodeTime = zapcore.ISO8601TimeEncoder
	prod.EncodeLevel = zapcore.CapitalColorLevelEncoder

	encoder := zapcore.NewConsoleEncoder(prod)
	syncer := zapcore.AddSync(os.Stdout)
	level := zapcore.DebugLevel

	core := zapcore.NewCore(encoder, syncer, level)

	opts := []zap.Option{
		zap.WithCaller(true),
		zap.AddStacktrace(zapcore.ErrorLevel),
		zap.AddCallerSkip(1),
	}
	lg := zap.New(core, opts...)

	return &sugarLog{sugar: lg.Sugar()}
}
