package app

import (
	"gin-zap-otel/logging"
	"gin-zap-otel/tracing"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// GinLogger returns a gin middleware for logging.
func GinLogger(logger *zap.Logger, tracer trace.Tracer, stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		requestCtx := &ctx
		if tracer != nil {
			ctxWithSpan, span := tracer.Start(
				tracing.NewContext(*requestCtx, tracer),
				c.Request.URL.Path,
			)
			defer span.End()

			requestCtx = &ctxWithSpan
		}

		ctxWithLogger := logging.NewContext(*requestCtx, logging.Get())
		c.Request = c.Request.WithContext(ctxWithLogger)

		defer func() {
			if err := recover(); err != nil {
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(
							strings.ToLower(se.Error()),
							"broken pipe",
						) || strings.Contains(
							strings.ToLower(
								se.Error()),
							"connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					c.Error(err.(error))
					c.Abort()
					return
				}

				if stack {
					logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()

		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		cost := time.Since(start)
		c.Next()

		logger.Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			// zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}
