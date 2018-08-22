package logger

import (
  "fmt"
  "github.com/Sirupsen/logrus"
  r "github.com/gomodule/redigo/redis"
)

// Define different log levels for this logger:
var (
  // InfoLevel applies to general operational entries about what's going on inside the app.
  InfoLevel = logrus.InfoLevel.String()

  // WarnLevel applies to non-critical entries that deserve eyes.
  WarnLevel = logrus.WarnLevel.String()
  
  // ErrorLevel applies to errors that should definitely be noted.
  ErrorLevel = logrus.ErrorLevel.String()
)

// Lgr is a wrapper type around `*logrus.Logger`, providing Redis stream functionality.
type Lgr struct {
  Logger *logrus.Logger
  RedisPool *r.Pool
  Stream string
}

// Info logs to the InfoLevel and the Redis log stream.
func (l *Lgr) Info(args ...interface{}) {
  msg := fmt.Sprint(args...)
  l.Logger.Info(msg)
  l.newStreamEntry(msg, InfoLevel)
}

// Infof logs a formatted string to the InfoLevel and the Redis stream.
func (l *Lgr) Infof(format string, args ...interface{}) {
  msg := fmt.Sprintf(format, args...)
  l.Logger.Info(msg)
  l.newStreamEntry(msg, InfoLevel)
}

// Infoln is equivalent to `Info` but appends a new line to the message.
func (l *Lgr) Infoln(args ...interface{}) {
  msg := fmt.Sprintln(args...)
  l.Logger.Info(msg)
  l.newStreamEntry(msg, InfoLevel)
}

// InternalInfo logs to the InfoLevel (no Redis stream).
func (l *Lgr) InternalInfo(args ...interface{}) {
  l.Logger.Info(args...)
}

// InternalInfof logs a formatted string to the InfoLevel (no Redis stream).
func (l *Lgr) InternalInfof(format string, args ...interface{}) {
  l.Logger.Infof(format, args...)
}

// InternalInfoln is equivalent to `InternalInfo` but appends a new line to the message.
func (l *Lgr) InternalInfoln(args ...interface{}) {
  l.Logger.Infoln(args...)
}

// Warn logs to the WarnLevel and the Redis log stream.
func (l *Lgr) Warn(args ...interface{}) {
  msg := fmt.Sprint(args...)
  l.Logger.Warn(msg)
  l.newStreamEntry(msg, WarnLevel)
}

// Warnf logs a formatted string to the WarnLevel and the Redis stream.
func (l *Lgr) Warnf(format string, args ...interface{}) {
  msg := fmt.Sprintf(format, args...)
  l.Logger.Warn(msg)
  l.newStreamEntry(msg, WarnLevel)
}

// Warnln is equivalent to `Warn` but appends a new line to the message.
func (l *Lgr) Warnln(args ...interface{}) {
  msg := fmt.Sprintln(args...)
  l.Logger.Warn(msg)
  l.newStreamEntry(msg, WarnLevel)
}

// InternalWarn logs to the WarnLevel (no Redis stream).
func (l *Lgr) InternalWarn(args ...interface{}) {
  l.Logger.Warn(args...)
}

// InternalWarnf logs a formatted string to the WarnLevel (no Redis stream).
func (l *Lgr) InternalWarnf(format string, args ...interface{}) {
  l.Logger.Warnf(format, args...)
}

// InternalWarnln is equivalent to `InternalWarn` but appends a new line to the message.
func (l *Lgr) InternalWarnln(args ...interface{}) {
  l.Logger.Warnln(args...)
}

// Error logs to the ErrorLevel and the Redis log stream.
func (l *Lgr) Error(args ...interface{}) {
  msg := fmt.Sprint(args...)
  l.Logger.Error(msg)
  l.newStreamEntry(msg, ErrorLevel)
}

// Errorf logs a formatted string to the ErrorLevel and the Redis stream.
func (l *Lgr) Errorf(format string, args ...interface{}) {
  msg := fmt.Sprintf(format, args...)
  l.Logger.Error(msg)
  l.newStreamEntry(msg, ErrorLevel)
}

// Errorln is equivalent to `Error` but appends a new line to the message.
func (l *Lgr) Errorln(args ...interface{}) {
  msg := fmt.Sprintln(args...)
  l.Logger.Error(msg)
  l.newStreamEntry(msg, ErrorLevel)
}

// InternalError logs to the ErrorLevel (no Redis stream).
func (l *Lgr) InternalError(args ...interface{}) {
  l.Logger.Error(args...)
}

// InternalErrorf logs a formatted string to the ErrorLevel (no Redis stream).
func (l *Lgr) InternalErrorf(format string, args ...interface{}) {
  l.Logger.Errorf(format, args...)
}

// InternalErrorln is equivalent to `InternalError` but appends a new line to the message.
func (l *Lgr) InternalErrorln(args ...interface{}) {
  l.Logger.Errorln(args...)
}

// newStreamEntry adds a new message to the logger's redis stream.
func (l *Lgr) newStreamEntry(msg string, level string) {
  // Get new connection from Redis pool.
  conn := l.RedisPool.Get()
  defer conn.Close()

  // Add message to log stream.
  if _, err := conn.Do("XADD", l.Stream, "*", "msg", msg, "level", level); err != nil {
    l.InternalErrorf("error logging to stream: %s", err.Error())
  }
}