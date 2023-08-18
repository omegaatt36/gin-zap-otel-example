package main

import (
	"context"
	"net/http"

	"gin-zap-otel/app"
	"gin-zap-otel/logging"
	"gin-zap-otel/tracing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func run(ctx context.Context) {
	ctx, span := tracing.Start(ctx, "run")
	defer span.End()

	logging.Debug("run~~~~~~~~", zap.String("run", "running"))
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
	router := gin.New()

	shutdownFn, err := tracing.Init()
	if err != nil {
		panic(err)
	}
	defer shutdownFn(context.Background())

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
