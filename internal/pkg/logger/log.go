package logger

import (
  "github.com/Sirupsen/logrus"
  r "github.com/gomodule/redigo/redis"
  "fmt"
)

// Define different log level keys.
const infoLevel = "info"
const warnLevel = "warn"
const errorLevel = "error"

// Wrapper type around `*logrus.Logger`, providing Redis stream functionality.
type lgr struct {
  Logger *logrus.Logger
  redis *r.Pool
  stream string
}

func (l *lgr) Info(args ...interface{}) {
  msg := fmt.Sprint(args...)
  l.Logger.Info(msg)
  l.newStreamEntry(msg, infoLevel)
}

func (l *lgr) Infof(format string, args ...interface{}) {
  msg := fmt.Sprintf(format, args...)
  l.Logger.Info(msg)
  l.newStreamEntry(msg, infoLevel)
}

func (l *lgr) Infoln(args ...interface{}) {
  msg := fmt.Sprintln(args...)
  l.Logger.Info(msg)
  l.newStreamEntry(msg, infoLevel)
}

func (l *lgr) InternalInfo(args ...interface{}) {
  l.Logger.Info(args...)
}

func (l *lgr) InternalInfof(format string, args ...interface{}) {
  l.Logger.Infof(format, args...)
}

func (l *lgr) InternalInfoln(args ...interface{}) {
  l.Logger.Infoln(args...)
}

func (l *lgr) Warn(args ...interface{}) {
  msg := fmt.Sprint(args...)
  l.Logger.Warn(msg)
  l.newStreamEntry(msg, warnLevel)
}

func (l *lgr) Warnf(format string, args ...interface{}) {
  msg := fmt.Sprintf(format, args...)
  l.Logger.Warn(msg)
  l.newStreamEntry(msg, warnLevel)
}

func (l *lgr) Warnln(args ...interface{}) {
  msg := fmt.Sprintln(args...)
  l.Logger.Warn(msg)
  l.newStreamEntry(msg, warnLevel)
}

func (l *lgr) InternalWarn(args ...interface{}) {
  l.Logger.Warn(args...)
}

func (l *lgr) InternalWarnf(format string, args ...interface{}) {
  l.Logger.Warnf(format, args...)
}

func (l *lgr) InternalWarnln(args ...interface{}) {
  l.Logger.Warnln(args...)
}

func (l *lgr) Error(args ...interface{}) {
  msg := fmt.Sprint(args...)
  l.Logger.Error(msg)
  l.newStreamEntry(msg, errorLevel)
}

func (l *lgr) Errorf(format string, args ...interface{}) {
  msg := fmt.Sprintf(format, args...)
  l.Logger.Error(msg)
  l.newStreamEntry(msg, errorLevel)
}

func (l *lgr) Errorln(args ...interface{}) {
  msg := fmt.Sprintln(args...)
  l.Logger.Error(msg)
  l.newStreamEntry(msg, errorLevel)
}

func (l *lgr) InternalError(args ...interface{}) {
  l.Logger.Error(args...)
}

func (l *lgr) InternalErrorf(format string, args ...interface{}) {
  l.Logger.Errorf(format, args...)
}

func (l *lgr) InternalErrorln(args ...interface{}) {
  l.Logger.Errorln(args...)
}

func (l *lgr) newStreamEntry(msg string, level string) {
  // TODO: Put "stage" inside payload here somewhere.
  // TODO: Figure out redis xadd Go functionality
}

// New creates and returns a pointer to a new `lgr` instance.
func New(pool *r.Pool, stream string) *lgr {
  return &lgr{
    Logger: logrus.New(),
    redis: pool,
    stream: stream,
  }
}