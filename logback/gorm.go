package logback

import (
	"context"
	"errors"
	"runtime"
	"strings"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Gorm(l *zap.Logger, level string) logger.Interface {
	lvl := logger.Error
	switch strings.ToLower(level) {
	case "debug", "info":
		lvl = logger.Info
	case "warn":
		lvl = logger.Warn
	case "error":
		lvl = logger.Error
	}

	return &gormLog{
		zlog:                      l,
		level:                     lvl,
		slowThreshold:             200 * time.Millisecond,
		skipCallerLookup:          false,
		ignoreRecordNotFoundError: true,
	}
}

type gormLog struct {
	zlog                      *zap.Logger
	level                     logger.LogLevel
	slowThreshold             time.Duration
	skipCallerLookup          bool
	ignoreRecordNotFoundError bool
}

func (gl *gormLog) Tracef(tmpl string, args ...any) {
	if gl.level >= logger.Info {
		gl.zlog.Sugar().Infof(tmpl, args...)
	}
}

func (gl *gormLog) Debugf(tmpl string, args ...any) {
	if gl.level >= logger.Info {
		gl.zlog.Sugar().Infof(tmpl, args...)
	}
}

func (gl *gormLog) Infof(tmpl string, args ...any) {
	if gl.level >= logger.Info {
		gl.zlog.Sugar().Infof(tmpl, args...)
	}
}

func (gl *gormLog) Warnf(tmpl string, args ...any) {
	if gl.level >= logger.Warn {
		gl.zlog.Sugar().Warnf(tmpl, args...)
	}
}

func (gl *gormLog) Errorf(tmpl string, args ...any) {
	if gl.level >= logger.Error {
		gl.zlog.Sugar().Errorf(tmpl, args...)
	}
}

func (gl *gormLog) LogMode(level logger.LogLevel) logger.Interface {
	return &gormLog{
		zlog:                      gl.zlog,
		level:                     level,
		slowThreshold:             gl.slowThreshold,
		skipCallerLookup:          gl.skipCallerLookup,
		ignoreRecordNotFoundError: gl.ignoreRecordNotFoundError,
	}
}

func (gl *gormLog) Info(_ context.Context, str string, args ...interface{}) {
	if gl.level < logger.Info {
		return
	}
	gl.gormGenCaller().Sugar().Infof(str, args...)
}

func (gl *gormLog) Warn(_ context.Context, str string, args ...interface{}) {
	if gl.level < logger.Warn {
		return
	}
	gl.gormGenCaller().Sugar().Warnf(str, args...)
}

func (gl *gormLog) Error(_ context.Context, str string, args ...interface{}) {
	if gl.level < logger.Error {
		return
	}
	gl.gormGenCaller().Sugar().Errorf(str, args...)
}

func (gl *gormLog) Trace(_ context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if gl.level < logger.Silent {
		return
	}
	elapsed := time.Since(begin)
	switch {
	case err != nil && gl.level >= logger.Error && (!gl.ignoreRecordNotFoundError || !errors.Is(err, gorm.ErrRecordNotFound)):
		sql, rows := fc()
		gl.gormGenCaller().Sugar().Errorf("[elapsed %s, rows %d] %s. error: %s", elapsed, rows, sql, err)
	case gl.slowThreshold != 0 && elapsed > gl.slowThreshold && gl.level >= logger.Warn:
		sql, rows := fc()
		if err != nil {
			gl.gormGenCaller().Sugar().Warnf("[elapsed %s, rows %d] %s, error: %s", elapsed, rows, sql, err.Error())
		} else {
			gl.gormGenCaller().Sugar().Warnf("[elapsed %s, rows %d] %s", elapsed, rows, sql)
		}
	case gl.level >= logger.Info:
		sql, rows := fc()
		if err != nil {
			gl.gormGenCaller().Sugar().Infof("[elapsed %s, rows %d] %s, error: %s", elapsed, rows, sql, err.Error())
		} else {
			gl.gormGenCaller().Sugar().Infof("[elapsed %s, rows %d] %s", elapsed, rows, sql)
		}
	}
}

// gormGenCaller 只适配了 [gorm.io/gen] 代码生成的场景，
// 并没有在纯 [gorm] 的场景下使用与测试。
//
// [gorm]: https://github.com/go-gorm/gorm
// [gorm.io/gen]: https://github.com/go-gorm/gen
func (gl *gormLog) gormGenCaller() *zap.Logger {
	var flag bool
	for i := 2; i < 8; i++ {
		_, file, _, ok := runtime.Caller(i)
		if !ok {
			break
		}
		if strings.HasSuffix(file, "/do.go") && strings.Contains(file, "gorm.io/gen") {
			flag = true
			continue
		}
		gen := strings.HasSuffix(file, ".gen.go")
		if flag && gen {
			return gl.zlog.WithOptions(zap.AddCallerSkip(i))
		}
		if flag && !gen {
			return gl.zlog.WithOptions(zap.AddCallerSkip(i - 1))
		}
	}

	return gl.zlog
}
