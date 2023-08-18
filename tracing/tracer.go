package tracing

import (
	"context"
	"log"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

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

func Get() trace.Tracer {
	return tracer
}

// Init initializes the OpenTelemetry tracing with span exporter
func Init() (func(context.Context) error, error) {
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize stdouttrace exporter")
	}

	bsp := sdktrace.NewBatchSpanProcessor(exporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSpanProcessor(bsp),
	)
	defer tracerProvider.ForceFlush(context.Background())
	otel.SetTracerProvider(tracerProvider)

	otel.SetTextMapPropagator(propagation.TraceContext{})

	tracer = tracerProvider.Tracer("kryptogo.com/trace")

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

// // Middleware returns a gin middleware that traces requests.
// func Middleware() gin.HandlerFunc {
// 	return otelgin.Middleware("kryptogo.com/trace", otelgin.WithTracerProvider(otel.GetTracerProvider()))
// }

// GetTraceSpanID returns the trace span id from the context.
func GetTraceSpanID(ctx context.Context) (string, string) {
	spanContext := trace.SpanContextFromContext(ctx)
	if spanContext.IsValid() {
		return spanContext.TraceID().String(), spanContext.SpanID().String()
	}

	return "", ""
}

// AsyncOp is a helper function to run an operation in a new goroutine with the same trace context. This is useful for detaching cancellation from the context but preserve the trace context.
func AsyncOp(ctx context.Context, spanName string, op func(ctx context.Context)) {
	t := FromContext(ctx)
	if t != nil {
		go func() {
			spanContext := trace.SpanContextFromContext(ctx)
			detachedCtx := trace.ContextWithSpanContext(context.Background(), spanContext)
			detachedCtx, span := t.Start(detachedCtx, spanName)
			defer span.End()
			op(detachedCtx)
		}()
	} else {
		go op(context.Background())
	}
}
