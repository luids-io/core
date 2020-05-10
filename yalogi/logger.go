// Package yalogi means "Yet Another Logger Interface" and provides a simple
// logger interface for use it in my projects.
//
// Feel free to use it in yours ;)
//
// This package is a work in progress and makes no API stability promises.
package yalogi

import (
	"bytes"
	"fmt"
	"log"
	"os"
)

// Level type can be used to classify the level of log messages
type Level int

// Constants for levels
const (
	Debug Level = iota
	Info
	Warning
	Error
	Fatal
)

// Logger is the main interface of the package
type Logger interface {
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})
}

// LogNull is an instance of a logger object that does nothing
var LogNull = &nullLogger{}

// nullLogger satisfies the interface with an implementation that does nothing.
type nullLogger struct{}

func (l *nullLogger) Debugf(template string, args ...interface{}) {}
func (l *nullLogger) Infof(template string, args ...interface{})  {}
func (l *nullLogger) Warnf(template string, args ...interface{})  {}
func (l *nullLogger) Errorf(template string, args ...interface{}) {}
func (l *nullLogger) Fatalf(template string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, template, args...)
	os.Exit(1)
}

type logWriter struct {
	intlog Logger
	level  Level
}

func (l logWriter) Write(p []byte) (int, error) {
	buff := bytes.NewBuffer(p)
	switch l.level {
	case Debug:
		l.intlog.Debugf(buff.String())
	case Info:
		l.intlog.Infof(buff.String())
	case Warning:
		l.intlog.Warnf(buff.String())
	case Error:
		l.intlog.Errorf(buff.String())
	case Fatal:
		l.intlog.Fatalf(buff.String())
	}
	return len(p), nil
}

// NewStandard creates a new instance of a golang standard log from an object
// that satisfaces the Logger interface. It is to provide compatibility in
// applications with packages that use the standard interface.
func NewStandard(l Logger, lvl Level) *log.Logger {
	tmp := &logWriter{intlog: l, level: lvl}
	return log.New(tmp, "", 0)
}
