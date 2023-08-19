package main

import (
	"context"
	"net/http"

	"gin-zap-otel/app"
	"gin-zap-otel/logging"
	"gin-zap-otel/tracing"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.uber.org/zap"
)

func run(ctx context.Context) {
	ctx, span := tracing.Start(ctx, "run")
	defer span.End()

	logging.Debug("running", zap.String("run", "running"))
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

	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		panic(errors.Wrap(err, "failed to initialize stdouttrace exporter"))
	}

	shutdownFn, err := tracing.Init(exporter, "gin-zap-otel")
	if err != nil {
		panic(err)
	}
	defer shutdownFn(context.Background())

	router := gin.New()

	router.Use(app.GinLogger(logging.Get(), tracing.Get(), false))

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

	router.GET("/panic", func(c *gin.Context) {
		panic("a")
	})

	router.Run(":8000")
}
