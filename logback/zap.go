package logback

import "go.uber.org/zap"

func Zap(l *zap.Logger) Logger {
	return &sugarLog{sugar: l.WithOptions(zap.AddCallerSkip(1)).Sugar()}
}

type sugarLog struct {
	sugar *zap.SugaredLogger
}

func (sg *sugarLog) Trace(v ...any)            { sg.sugar.Debug(v...) }
func (sg *sugarLog) Debug(v ...any)            { sg.sugar.Debug(v...) }
func (sg *sugarLog) Info(v ...any)             { sg.sugar.Info(v...) }
func (sg *sugarLog) Warn(v ...any)             { sg.sugar.Warn(v...) }
func (sg *sugarLog) Error(v ...any)            { sg.sugar.Error(v...) }
func (sg *sugarLog) Tracef(s string, v ...any) { sg.sugar.Debugf(s, v...) }
func (sg *sugarLog) Debugf(s string, v ...any) { sg.sugar.Debugf(s, v...) }
func (sg *sugarLog) Infof(s string, v ...any)  { sg.sugar.Infof(s, v...) }
func (sg *sugarLog) Warnf(s string, v ...any)  { sg.sugar.Warnf(s, v...) }
func (sg *sugarLog) Errorf(s string, v ...any) { sg.sugar.Errorf(s, v...) }
func (sg *sugarLog) Replace(l *zap.Logger)     { sg.sugar = l.WithOptions(zap.AddCallerSkip(1)).Sugar() }
