package logger

import (
	"context"
	"io"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
)

var (
	loggerGlobal = New(defaultLevel, os.Stdout)
	defaultLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
)

func Init(level string) {
	var logLvl zap.AtomicLevel

	switch strings.ToLower(level) {
	case "debug":
		logLvl = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		logLvl = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		logLvl = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		logLvl = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		logLvl = defaultLevel
	}

	loggerGlobal = New(logLvl, os.Stdout)
}

func New(level zap.AtomicLevel, sink io.Writer, opts ...zap.Option) *zap.SugaredLogger {
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}),
		zapcore.AddSync(sink),
		level,
	)

	return zap.New(core, opts...).Sugar()
}

func GetLogger() *zap.SugaredLogger {
	return loggerGlobal
}

func Debug(args ...interface{}) {
	loggerGlobal.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	loggerGlobal.Debugf(template, args...)
}

func Info(args ...interface{}) {
	loggerGlobal.Info(args...)
}

func Infof(template string, args ...interface{}) {
	loggerGlobal.Infof(template, args...)
}

func Warn(args ...interface{}) {
	loggerGlobal.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	loggerGlobal.Warnf(template, args...)
}

func Error(args ...interface{}) {
	loggerGlobal.Error(args...)
}

func Errorf(ctx context.Context, method, template string, args ...interface{}) {
	withTraceID(ctx).Desugar().
		With(zap.String("method", method)).Sugar().Errorf(template, args...)
}

func Fatal(args ...interface{}) {
	loggerGlobal.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	loggerGlobal.Fatalf(template, args...)
}

func withTraceID(ctx context.Context) *zap.SugaredLogger {
	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		return loggerGlobal
	}

	if sc, ok := span.Context().(jaeger.SpanContext); ok {
		return loggerGlobal.Desugar().With(
			zap.Stringer("trace_id", sc.TraceID()),
			zap.Stringer("span_id", sc.SpanID()),
		).Sugar()
	}

	return loggerGlobal
}
