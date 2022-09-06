package orm

import (
	"context"
	"errors"
	"fmt"
	core "github.com/cheivin/dio-core"
	"gorm.io/gorm/logger"
	"time"
)

type gormLogger struct {
	log                       core.Log
	slowThreshold             time.Duration
	ignoreRecordNotFoundError bool
	level                     logger.LogLevel
}

func newLogger(log core.Log, property LogProperty) logger.Interface {
	return &gormLogger{
		log:                       log.Skip(4),
		level:                     property.LogLevel(),
		slowThreshold:             property.SlowThreshold,
		ignoreRecordNotFoundError: property.IgnoreRecordNotFoundError,
	}
}

// LogMode log mode
func (l *gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return &gormLogger{
		slowThreshold:             l.slowThreshold,
		ignoreRecordNotFoundError: l.ignoreRecordNotFoundError,
		level:                     level,
		log:                       l.log,
	}
}

// Info print info
func (l gormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= logger.Info {
		l.log.Info(ctx, fmt.Sprintf("%s\n "+msg, data...))
	}
}

// Warn print warn messages
func (l gormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= logger.Warn {
		l.log.Warn(ctx, fmt.Sprintf("%s\n "+msg, data...))
	}
}

// Error print error messages
func (l gormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= logger.Error {
		l.log.Error(ctx, fmt.Sprintf("%s\n "+msg, data...))
	}
}

// Trace print sql message
func (l gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.level > logger.Silent {
		elapsed := time.Since(begin)
		switch {
		case err != nil && l.level >= logger.Error && (!errors.Is(err, logger.ErrRecordNotFound) || !l.ignoreRecordNotFoundError):
			sql, rows := fc()
			if rows == -1 {
				l.log.Error(ctx, err.Error(), "Cost", float64(elapsed.Nanoseconds())/1e6, "SQL", sql)
			} else {
				l.log.Error(ctx, err.Error(), "Cost", float64(elapsed.Nanoseconds())/1e6, "Rows", rows, "SQL", sql)
			}
		case l.slowThreshold != 0 && elapsed > l.slowThreshold && l.level >= logger.Warn:
			sql, rows := fc()
			slowLog := fmt.Sprintf("SLOW SQL >= %v", l.slowThreshold)
			if rows == -1 {
				l.log.Warn(ctx, slowLog, "Cost", float64(elapsed.Nanoseconds())/1e6, "SQL", sql)
			} else {
				l.log.Warn(ctx, slowLog, "Cost", float64(elapsed.Nanoseconds())/1e6, "Rows", rows, "SQL", sql)
			}
		case l.level == logger.Info:
			sql, rows := fc()
			if rows == -1 {
				l.log.Debug(ctx, "", "Cost", float64(elapsed.Nanoseconds())/1e6, "SQL", sql)
			} else {
				l.log.Debug(ctx, "", "Cost", float64(elapsed.Nanoseconds())/1e6, "Rows", rows, "SQL", sql)
			}
		}
	}
}
