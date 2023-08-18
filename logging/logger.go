//go:generate go-enum
package logging

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ctxKey is the type of value for the context key.
type ctxKey struct{}

// NewContext returns a new context with the logger instance.
func NewContext(parent context.Context, z *zap.Logger) context.Context {
	return context.WithValue(parent, ctxKey{}, z)
}

func FromContext(ctx context.Context) *zap.Logger {
	c, _ := ctx.Value(ctxKey{}).(*zap.Logger)
	return c
}

// ENUM(
// Default = DEFAULT
// Debug = DEBUG
// Info = INFO
// Notice = NOTICE
// Warning = WARNING
// Error = ERROR
// Critical = CRITICAL
// Alert = ALERT
// Emergency = EMERGENCY
// )
type severity string

var logger *zap.Logger

// Get returns the logger instance.
func Get() *zap.Logger {
	return logger
}

func init() {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	// if !config.IsLocal() {
	// encoderConfig = zap.NewProductionEncoderConfig()
	// encoder = zapcore.NewJSONEncoder(encoderConfig)
	// }

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
	)

	logger = zap.New(core, zap.AddCaller())

	Info = logger.Sugar().Info
	Infof = logger.Sugar().Infof
	Debug = logger.Sugar().Debug
	Debugf = logger.Sugar().Debugf
	Error = logger.Sugar().Error
	Errorf = logger.Sugar().Errorf
	Fatal = logger.Sugar().Fatal
	Fatalf = logger.Sugar().Fatalf
	Panic = logger.Sugar().Panic
	Panicf = logger.Sugar().Panicf
	Warn = logger.Sugar().Warn
	Warnf = logger.Sugar().Warnf
}

type printWrapper func(args ...interface{})
type printfWrapper func(template string, args ...interface{})

var (
	Info   printWrapper
	Infof  printfWrapper
	Debug  printWrapper
	Debugf printfWrapper
	Error  printWrapper
	Errorf printfWrapper
	Fatal  printWrapper
	Fatalf printfWrapper
	Panic  printWrapper
	Panicf printfWrapper
	Warn   printWrapper
	Warnf  printfWrapper
)
