package logger

import (
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"time"
)

var logg *zap.Logger

// New New
func New(name string) error {
	config := getConfig()

	defaultFields := zap.Fields(
		zap.String("app", name),
	)

	var err error

	logg, err = config.Build(
		defaultFields,
	)
	if err != nil {
		return err
	}
	defer logg.Sync()

	return nil
}

func getConfig() zap.Config {
	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:          "json",
		EncoderConfig:     getEncoderConfig(),
		OutputPaths:       []string{"stderr"},
		ErrorOutputPaths:  []string{"stderr"},
		DisableStacktrace: false,
		DisableCaller:     true,
	}

	return config
}

func getEncoderConfig() zapcore.EncoderConfig {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	return encoderConfig
}

// APILogger is a middleware that logs the start and end of each request, along
// with some useful data about what was requested, what the response status was,
// and how long it took to return.
func APILogger() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				logg.Info("Served",
					zap.String("proto", r.Proto),
					zap.String("path", r.URL.Path),
					zap.String("method", r.Method),
					zap.Duration("lat", time.Since(t1)),
					zap.Int("status", ww.Status()),
					zap.Int("size", ww.BytesWritten()),
					zap.String("reqId", middleware.GetReqID(r.Context())))
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}

// Info Info
func Info(msg string, fields ...zap.Field) {

	logg.Info(msg, fields...)
}

// Warning Warning
func Warning(msg string, fields ...zap.Field) {
	logg.Warn(msg, fields...)
}

// Error Error
func Error(msg string, fields ...zap.Field) {
	logg.Error(msg, fields...)
}

// Debug Debug
func Debug(msg string, fields ...zap.Field) {
	logg.Debug(msg, fields...)
}

// Fatal Fatal
func Fatal(msg string, fields ...zap.Field) {
	logg.Fatal(msg, fields...)
}
