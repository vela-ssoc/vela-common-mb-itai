package logback

import "go.uber.org/zap"

// Logger 日志接口定义
type Logger interface {
	Trace(...any)
	Debug(...any)
	Info(...any)
	Warn(...any)
	Error(...any)
	Tracef(string, ...any)
	Debugf(string, ...any)
	Infof(string, ...any)
	Warnf(string, ...any)
	Errorf(string, ...any)
	Replace(*zap.Logger)
}
