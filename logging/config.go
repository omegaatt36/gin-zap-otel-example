//go:generate go-enum
package logging

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ENUM(
// Development = development
// Production = production
// )
type Env string

// TraceInjector is an interface for injecting trace information into the logger.
type TraceInjector interface {
	Inject(context.Context, *zap.Logger)
}

// TraceInjectorDefault is the default trace injector.
type traceInjectorDefault struct {
}

// NewDefaultTraceInjector returns a new default trace injector.
func NewDefaultTraceInjector() TraceInjector {
	return traceInjectorDefault{}
}

var _ TraceInjector = traceInjectorDefault{}

// Inject injects trace information into the logger.
func (traceInjectorDefault) Inject(ctx context.Context, l *zap.Logger) {
	// decuple with package tracing.
	if traceID := trace.SpanFromContext(ctx).SpanContext().TraceID(); traceID.IsValid() {
		*l = *l.With(zap.String("trace-id", traceID.String()))
	}

	if spanID := trace.SpanFromContext(ctx).SpanContext().SpanID(); spanID.IsValid() {
		*l = *l.With(zap.String("span-id", spanID.String()))
	}
}

// Config is the configuration for the logger.
type Config struct {
	Environment   Env
	Level         zapcore.Level
	TraceInjector TraceInjector
}
