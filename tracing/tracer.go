package tracing

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// InfluenceTraceIDHeader is the header name for the trace id.
const InfluenceTraceIDHeader = "X-Influence-Trace-ID"

// ctxKey is the type of value for the context key.
type ctxKey struct{}

// NewContext returns a new context with the tracer instance.
func NewContext(parent context.Context, t trace.Tracer) context.Context {
	return context.WithValue(parent, ctxKey{}, t)
}

// FromContext returns the tracer instance from the context.
func FromContext(ctx context.Context) trace.Tracer {
	t, ok := ctx.Value(ctxKey{}).(trace.Tracer)
	if !ok {
		log.Println("failed to get tracer from context, using noop tracer")
		return tracer
	}

	return t
}

// tracer is the global tracer used by the app.
var tracer trace.Tracer

// Get returns the tracer instance.
func Get() trace.Tracer {
	return tracer
}

// Init initializes the OpenTelemetry tracing with span exporter
func Init(exporter sdktrace.SpanExporter, tracerName string) (func(context.Context) error, error) {
	bsp := sdktrace.NewBatchSpanProcessor(exporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSpanProcessor(bsp),
	)
	defer tracerProvider.ForceFlush(context.Background())

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	tracer = tracerProvider.Tracer(tracerName)

	return tracerProvider.Shutdown, nil
}

// Start starts a new span with the given name and returns the span and the context. If no tracer is set, a noop tracer is used.
func Start(ctx context.Context, spanName string) (context.Context, trace.Span) {
	t := FromContext(ctx)
	if t != nil {
		return t.Start(ctx, spanName)
	}

	return trace.NewNoopTracerProvider().Tracer("noop").Start(ctx, spanName)
}

// GetTraceIDAndSpanID returns the trace id and span id from the context.
func GetTraceIDAndSpanID(ctx context.Context) (string, string) {
	spanContext := trace.SpanContextFromContext(ctx)
	if spanContext.IsValid() {
		return spanContext.TraceID().String(), spanContext.SpanID().String()
	}

	return "", ""
}

// AsyncFn starts a new span with the given name and calls the given function in a new goroutine.
// If no tracer is set, a noop tracer is used.
func AsyncFn(ctx context.Context, spanName string, fn func(ctx context.Context)) {
	t := FromContext(ctx)
	if t != nil {
		go func() {
			spanContext := trace.SpanContextFromContext(ctx)
			detachedCtx := trace.ContextWithSpanContext(context.Background(), spanContext)
			detachedCtx, span := t.Start(detachedCtx, spanName)
			defer span.End()
			fn(detachedCtx)
		}()
	} else {
		go fn(ctx)
	}
}
