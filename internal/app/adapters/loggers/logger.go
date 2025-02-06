package loggers

import (
	"go.uber.org/zap/zapcore"
	"os"
	"sync"

	"go.uber.org/zap"
)

var (
	loggerInstance *Logger
	once           sync.Once
)

type Logger struct {
	logger *zap.Logger
}

func GetLogger() *Logger {
	once.Do(func() {
		loggerInstance = &Logger{
			logger: NewLogger("development"),
		}
	})
	return loggerInstance
}

func NewLogger(env string) *zap.Logger {
	var logger *zap.Logger

	if env == "production" {
		logger, _ = zap.NewProduction()
	} else {
		encoderConfig := zapcore.EncoderConfig{
			TimeKey:      "time",
			LevelKey:     "level",
			MessageKey:   "msg",
			EncodeLevel:  zapcore.CapitalColorLevelEncoder, // Màu
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		}

		core := zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			zapcore.DebugLevel,
		)

		logger = zap.New(core)
	}

	return logger
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}

func (l *Logger) Sync() error {
	return l.logger.Sync()
}
