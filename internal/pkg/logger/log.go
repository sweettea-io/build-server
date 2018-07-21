package logger

import (
  "fmt"
  "github.com/Sirupsen/logrus"
  r "github.com/gomodule/redigo/redis"
)

// Define different log level keys.
const infoLevel = "info"
const warnLevel = "warn"
const errorLevel = "error"

// Wrapper type around `*logrus.Logger`, providing Redis stream functionality.
type Lgr struct {
  Logger *logrus.Logger
  Redis *r.Pool
  Stream string
}

func (l *Lgr) Info(args ...interface{}) {
  msg := fmt.Sprint(args...)
  l.Logger.Info(msg)
  l.newStreamEntry(msg, infoLevel)
}

func (l *Lgr) Infof(format string, args ...interface{}) {
  msg := fmt.Sprintf(format, args...)
  l.Logger.Info(msg)
  l.newStreamEntry(msg, infoLevel)
}

func (l *Lgr) Infoln(args ...interface{}) {
  msg := fmt.Sprintln(args...)
  l.Logger.Info(msg)
  l.newStreamEntry(msg, infoLevel)
}

func (l *Lgr) InternalInfo(args ...interface{}) {
  l.Logger.Info(args...)
}

func (l *Lgr) InternalInfof(format string, args ...interface{}) {
  l.Logger.Infof(format, args...)
}

func (l *Lgr) InternalInfoln(args ...interface{}) {
  l.Logger.Infoln(args...)
}

func (l *Lgr) Warn(args ...interface{}) {
  msg := fmt.Sprint(args...)
  l.Logger.Warn(msg)
  l.newStreamEntry(msg, warnLevel)
}

func (l *Lgr) Warnf(format string, args ...interface{}) {
  msg := fmt.Sprintf(format, args...)
  l.Logger.Warn(msg)
  l.newStreamEntry(msg, warnLevel)
}

func (l *Lgr) Warnln(args ...interface{}) {
  msg := fmt.Sprintln(args...)
  l.Logger.Warn(msg)
  l.newStreamEntry(msg, warnLevel)
}

func (l *Lgr) InternalWarn(args ...interface{}) {
  l.Logger.Warn(args...)
}

func (l *Lgr) InternalWarnf(format string, args ...interface{}) {
  l.Logger.Warnf(format, args...)
}

func (l *Lgr) InternalWarnln(args ...interface{}) {
  l.Logger.Warnln(args...)
}

func (l *Lgr) Error(args ...interface{}) {
  msg := fmt.Sprint(args...)
  l.Logger.Error(msg)
  l.newStreamEntry(msg, errorLevel)
}

func (l *Lgr) Errorf(format string, args ...interface{}) {
  msg := fmt.Sprintf(format, args...)
  l.Logger.Error(msg)
  l.newStreamEntry(msg, errorLevel)
}

func (l *Lgr) Errorln(args ...interface{}) {
  msg := fmt.Sprintln(args...)
  l.Logger.Error(msg)
  l.newStreamEntry(msg, errorLevel)
}

func (l *Lgr) InternalError(args ...interface{}) {
  l.Logger.Error(args...)
}

func (l *Lgr) InternalErrorf(format string, args ...interface{}) {
  l.Logger.Errorf(format, args...)
}

func (l *Lgr) InternalErrorln(args ...interface{}) {
  l.Logger.Errorln(args...)
}

func (l *Lgr) newStreamEntry(msg string, level string) {
  // TODO: Put "stage" inside payload here somewhere.
  // TODO: Figure out redis xadd Go functionality
}