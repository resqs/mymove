package logging

import (
	"net/http"

	"github.com/felixge/httpsnoop"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/transcom/mymove/pkg/app"
	"github.com/transcom/mymove/pkg/auth"
)

// Config configures a Zap logger based on the environment string and debugLevel
func Config(env string, debugLogging bool) (*zap.Logger, error) {
	var loggerConfig zap.Config

	if env != "development" {
		loggerConfig = zap.NewProductionConfig()
	} else {
		loggerConfig = zap.NewDevelopmentConfig()
	}

	loggerConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	if debugLogging {
		debug := zap.NewAtomicLevel()
		debug.SetLevel(zap.DebugLevel)
		loggerConfig.Level = debug
	}

	return loggerConfig.Build()
}

// LogRequestMiddleware generates an HTTP/HTTPS request logs using Zap
func LogRequestMiddleware(inner http.Handler) http.Handler {
	mw := func(w http.ResponseWriter, r *http.Request) {
		var protocol string
		var userID string

		if r.TLS == nil {
			protocol = "http"
		} else {
			protocol = "https"
		}
		metrics := httpsnoop.CaptureMetrics(inner, w, r)
		userUUID, ok := auth.GetUserID(r.Context())
		if ok {
			userID = userUUID.String()
		}
		zap.L().Info("Request",
			zap.String("accepted-language", r.Header.Get("accepted-language")),
			zap.Int64("content-length", r.ContentLength),
			zap.Duration("duration", metrics.Duration),
			zap.String("host", r.Host),
			zap.String("method", r.Method),
			zap.String("protocol", protocol),
			zap.String("protocol-version", r.Proto),
			zap.String("referer", r.Header.Get("referer")),
			zap.Int64("resp-size-bytes", metrics.Written),
			zap.Int("resp-status", metrics.Code),
			zap.String("source", r.RemoteAddr),
			zap.String("url", r.URL.String()),
			zap.String("user-agent", r.UserAgent()),
			zap.String("x-amzn-trace-id", r.Header.Get("x-amzn-trace-id")),
			zap.String("x-forwarded-for", r.Header.Get("x-forwarded-for")),
			zap.String("x-forwarded-host", r.Header.Get("x-forwarded-host")),
			zap.String("x-forwarded-proto", r.Header.Get("x-forwarded-proto")),
			zap.String("user-id", userID),
			zap.String("app", app.GetAppFromContext(r)),
		)

	}
	return http.HandlerFunc(mw)
}
