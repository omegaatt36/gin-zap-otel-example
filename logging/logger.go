//go:generate go-enum
package logging

import (
	"context"
	"os"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ctxKey is the type of value for the context key.
type ctxKey struct{}

// NewContext returns a new context with the logger instance.
func NewContext(parent context.Context, z *zap.Logger) context.Context {
	return context.WithValue(parent, ctxKey{}, z)
}

// FromContext returns the logger instance from the context.
func FromContext(ctx context.Context) *zap.Logger {
	c, ok := ctx.Value(ctxKey{}).(*zap.Logger)
	if !ok {
		return logger
	}

	logger := *c

	// decuple with package tracing.
	if traceID := trace.SpanFromContext(ctx).SpanContext().TraceID(); traceID.IsValid() {
		logger = *logger.With(zap.String("trace-id", traceID.String()))
	}

	if spanID := trace.SpanFromContext(ctx).SpanContext().SpanID(); spanID.IsValid() {
		logger = *logger.With(zap.String("span-id", spanID.String()))
	}

	return &logger
}

var logger *zap.Logger

// Get returns the logger instance.
func Get() *zap.Logger {
	return logger
}

// ENUM(
// Development = development
// Production = production
// )
type Env string

type Config struct {
	Environment Env
	Level       zapcore.Level
}

// Init initializes the logger.
func Init(cfg Config) {
	level := zapcore.DebugLevel
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	if cfg.Environment == EnvProduction {
		level = zapcore.InfoLevel
		encoderConfig = zap.NewProductionEncoderConfig()
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	if cfg.Level != level {
		level = cfg.Level
	}

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level),
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

type printWrapper func(args ...any)
type printfWrapper func(template string, args ...any)

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

// DebugCtx logs a message at level Debug on the logger associated with the context.
func DebugCtx(ctx context.Context, message any) {
	FromContext(ctx).Sugar().Debug(message)
}

// DebugfCtx logs a message at level Debug on the logger associated with the context.
func DebugfCtx(ctx context.Context, template string, args ...any) {
	FromContext(ctx).Sugar().Debugf(template, args...)
}

// DebugWithData logs a message at level Debug on the logger associated with the context.
func DebugWithData(message any, data any) {
	logger.With(zap.Any("data", data)).Sugar().Debug(message)
}

// DebugWithDataCtx logs a message at level Debug on the logger associated with the context.
func DebugWithDataCtx(ctx context.Context, message any, data any) {
	FromContext(ctx).With(zap.Any("data", data)).Sugar().Debug(message)
}

// InfoCtx logs a message at level Info on the logger associated with the context.
func InfoCtx(ctx context.Context, message any) {
	FromContext(ctx).Sugar().Info(message)
}

// InfofCtx logs a message at level Info on the logger associated with the context.
func InfofCtx(ctx context.Context, template string, args ...any) {
	FromContext(ctx).Sugar().Infof(template, args...)
}

// InfoWithData logs a message at level Info on the logger associated with the context.
func InfoWithData(message any, data any) {
	logger.With(zap.Any("data", data)).Sugar().Info(message)
}

// InfoWithDataCtx logs a message at level Info on the logger associated with the context.
func InfoWithDataCtx(ctx context.Context, message any, data any) {
	FromContext(ctx).With(zap.Any("data", data)).Sugar().Info(message)
}

// WarnCtx logs a message at level Warn on the logger associated with the context.
func WarnCtx(ctx context.Context, message any) {
	FromContext(ctx).Sugar().Warn(message)
}

// WarnfCtx logs a message at level Warn on the logger associated with the context.
func WarnfCtx(ctx context.Context, template string, args ...any) {
	FromContext(ctx).Sugar().Warnf(template, args...)
}

// WarnWithData logs a message at level Warn on the logger associated with the context.
func WarnWithData(message any, data any) {
	logger.With(zap.Any("data", data)).Sugar().Warn(message)
}

// WarnWithDataCtx logs a message at level Warn on the logger associated with the context.
func WarnWithDataCtx(ctx context.Context, message any, data any) {
	FromContext(ctx).With(zap.Any("data", data)).Sugar().Warn(message)
}

// ErrorCtx logs a message at level Error on the logger associated with the context.
func ErrorCtx(ctx context.Context, message any) {
	FromContext(ctx).Sugar().Error(message)
}

// ErrorfCtx logs a message at level Error on the logger associated with the context.
func ErrorfCtx(ctx context.Context, template string, args ...any) {
	FromContext(ctx).Sugar().Errorf(template, args...)
}

// ErrorWithData logs a message at level Error on the logger associated with the context.
func ErrorWithData(message any, data any) {
	logger.With(zap.Any("data", data)).Sugar().Error(message)
}

// ErrorWithDataCtx logs a message at level Error on the logger associated with the context.
func ErrorWithDataCtx(ctx context.Context, message any, data any) {
	FromContext(ctx).With(zap.Any("data", data)).Sugar().Error(message)
}

// FatalCtx logs a message at level Fatal on the logger associated with the context.
func FatalCtx(ctx context.Context, message any) {
	FromContext(ctx).Sugar().Fatal(message)
}

// FatalfCtx logs a message at level Fatal on the logger associated with the context.
func FatalfCtx(ctx context.Context, template string, args ...any) {
	FromContext(ctx).Sugar().Fatalf(template, args...)
}

// FatalWithData logs a message at level Fatal on the logger associated with the context.
func FatalWithData(message any, data any) {
	logger.With(zap.Any("data", data)).Sugar().Fatal(message)
}

// FatalWithDataCtx logs a message at level Fatal on the logger associated with the context.
func FatalWithDataCtx(ctx context.Context, message any, data any) {
	FromContext(ctx).With(zap.Any("data", data)).Sugar().Fatal(message)
}

// PanicCtx logs a message at level Panic on the logger associated with the context.
func PanicCtx(ctx context.Context, message any) {
	FromContext(ctx).Sugar().Panic(message)
}

// PanicfCtx logs a message at level Panic on the logger associated with the context.
func PanicfCtx(ctx context.Context, template string, args ...any) {
	FromContext(ctx).Sugar().Panicf(template, args...)
}

// PanicWithData logs a message at level Panic on the logger associated with the context.
func PanicWithData(message any, data any) {
	logger.With(zap.Any("data", data)).Sugar().Panic(message)
}

// PanicWithDataCtx logs a message at level Panic on the logger associated with the context.
func PanicWithDataCtx(ctx context.Context, message any, data any) {
	FromContext(ctx).With(zap.Any("data", data)).Sugar().Panic(message)
}
