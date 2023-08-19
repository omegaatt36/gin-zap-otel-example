package main

import (
	"context"
	"net/http"
	"time"

	"gin-zap-otel/app"
	"gin-zap-otel/logging"
	"gin-zap-otel/tracing"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.uber.org/zap"
)

const serviceName = "gin-zap-otel"

func run(ctx context.Context) {
	ctx, span := tracing.Start(ctx, "run")
	defer span.End()

	logging.Debug("running")
	logging.FromContext(ctx).Info("span test 1 with trace")
	logging.FromContext(ctx).Info("span test 1 repeats trace and span id")
}

func run2(ctx context.Context) {
	ctx, span := tracing.Start(ctx, "run2")
	defer span.End()

	logging.FromContext(ctx).Info("span test 2 with new trace")
	run3(ctx)
}

func run3(ctx context.Context) {
	ctx, span := tracing.Start(ctx, "run3")
	defer span.End()

	logging.FromContext(ctx).Info("span test 3 repeates trace id from 2 with new span id")
}

func main() {
	logging.Init(logging.Config{
		Environment: logging.EnvDevelopment,
		Level:       zap.DebugLevel,
	})

	// exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(
		jaeger.WithEndpoint("http://localhost:14268/api/traces"),
	))
	if err != nil {
		panic(errors.Wrap(err, "failed to initialize stdouttrace exporter"))
	}

	shutdownFn, err := tracing.Init(exporter, serviceName)
	if err != nil {
		panic(err)
	}
	defer shutdownFn(context.Background())

	router := gin.New()

	router.Use(app.GinLogger(logging.Get(), tracing.Get(), true))

	router.GET("/", func(c *gin.Context) {
		run(c.Request.Context())
		run2(c.Request.Context())
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, World!",
		})
	})

	router.GET("/error", func(c *gin.Context) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
	})

	router.GET("/routine", func(c *gin.Context) {
		tracing.AsyncFn(c.Request.Context(), "routine", func(ctx context.Context) {
			time.Sleep(time.Second * 5)
			logging.FromContext(ctx).Info("after 5 seconds")
		})

		time.Sleep(time.Second * 2)

		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, World!",
		})
	})

	router.GET("/panic", func(c *gin.Context) {
		panic("a")
	})

	router.Run(":8000")
}
