package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	zap *zap.Logger
}

func New(level string) (*Logger, error) {
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(level)); err != nil {
		zapLevel = zapcore.InfoLevel
	}

	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(zapLevel),
		Development: false,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:      "ts",
			LevelKey:     "level",
			MessageKey:   "msg",
			CallerKey:    "caller",
			EncodeLevel:  zapcore.LowercaseLevelEncoder,
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	zapLogger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return &Logger{zap: zapLogger}, nil
}

func (l *Logger) Info(msg string, args ...any) {
	l.zap.Sugar().Infow(msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	l.zap.Sugar().Errorw(msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	l.zap.Sugar().Warnw(msg, args...)
}

func (l *Logger) Debug(msg string, args ...any) {
	l.zap.Sugar().Debugw(msg, args...)
}
