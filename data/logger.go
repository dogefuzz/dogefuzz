package data

import (
	"context"
	"time"

	"go.uber.org/zap"
	gormlogger "gorm.io/gorm/logger"
)

type DataLogger struct {
	Logger   *zap.Logger
	LogLevel gormlogger.LogLevel
}

func (l DataLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return DataLogger{
		Logger:   l.Logger,
		LogLevel: level,
	}
}

func (l DataLogger) Info(ctx context.Context, message string, args ...interface{}) {
	if l.LogLevel < gormlogger.Info {
		return
	}
	l.Logger.Sugar().Infof(message, args...)
}

func (l DataLogger) Warn(ctx context.Context, message string, args ...interface{}) {
	if l.LogLevel < gormlogger.Warn {
		return
	}
	l.Logger.Sugar().Warnf(message, args...)
}

func (l DataLogger) Error(ctx context.Context, message string, args ...interface{}) {
	if l.LogLevel < gormlogger.Error {
		return
	}
	l.Logger.Sugar().Errorf(message, args...)
}

func (l DataLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	// TODO
}
