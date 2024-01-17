package deps

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"userCRUD/internal/common/constants"
)

type Logger interface {
	Info(ctx context.Context, msg string, keysAndValues ...interface{})
	Debug(ctx context.Context, msg string, keysAndValues ...interface{})
	Error(ctx context.Context, msg string, keysAndValues ...interface{})
}

type ZapLogger struct {
	*zap.Logger
}

func NewZapLogger() *ZapLogger {
	config := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		OutputPaths:      []string{"stdout", "./log.log"},
		ErrorOutputPaths: []string{"stdout", "./logerrors.log"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			CallerKey:      "caller",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.RFC3339TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}

	logger, _ := config.Build(
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)

	return &ZapLogger{
		logger,
	}
}

func (l *ZapLogger) Info(ctx context.Context, msg string, keysAndValues ...interface{}) {
	l.Sugar().Infow(msg, withTraceId(ctx, keysAndValues)...)
}

func (l *ZapLogger) Debug(ctx context.Context, msg string, keysAndValues ...interface{}) {
	l.Sugar().Debugw(msg, withTraceId(ctx, keysAndValues)...)
}

func (l *ZapLogger) Error(ctx context.Context, msg string, keysAndValues ...interface{}) {
	l.Sugar().Errorw(msg, withTraceId(ctx, keysAndValues)...)
}

func withTraceId(ctx context.Context, keysAndValues []interface{}) []interface{} {
	if traceID, ok := ctx.Value(constants.TraceId).(string); ok {
		return append([]interface{}{constants.TraceId, traceID}, keysAndValues...)
	}
	return keysAndValues
}
