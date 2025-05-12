// Package slogresty 实现 resty.Logger 接口
package slogresty

import (
	"context"
	"log/slog"

	"github.com/go-resty/resty/v2"

	"bk.tencent.com/{{cookiecutter.project_name}}/{{cookiecutter.project_name}}/pkg/logging"
)

// Logger 用于实现 resty.Logger
type Logger struct {
	ctx context.Context
}

// New 实例化 Logger
func New(ctx context.Context) *Logger {
	return &Logger{ctx: ctx}
}

// Errorf ...
func (l *Logger) Errorf(format string, v ...any) {
	logging.Logf(l.ctx, slog.LevelError, format, v...)
}

// Warnf ...
func (l *Logger) Warnf(format string, v ...any) {
	logging.Logf(l.ctx, slog.LevelWarn, format, v...)
}

// Debugf ...
func (l *Logger) Debugf(format string, v ...any) {
	logging.Logf(l.ctx, slog.LevelDebug, format, v...)
}

var _ resty.Logger = (*Logger)(nil)
