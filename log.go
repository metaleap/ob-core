package obcore

import (
	"io"
	"log"
)

//	An interface for log output. ObLogger provides the canonical implementation
type Logger interface {
	// Debugf formats its arguments according to the format, analogous to fmt.Printf,
	// and records the text as a log message at Debug level
	Debugf(format string, args ...interface{})

	// Infof is like Debugf, but at Info level
	Infof(format string, args ...interface{})

	// Warningf is like Debugf, but at Warning level
	Warningf(format string, args ...interface{})

	// Error is like Debugf, but at Error level
	Error(err error) error

	// Errorf is like Debugf, but at Error level
	Errorf(format string, args ...interface{})

	// Criticalf is like Debugf, but at Critical level
	Criticalf(format string, args ...interface{})
}

//	The canonical implementation of the Logger interface
type ObLogger struct {
	logger *log.Logger
}

//	Creates and returns a new ObLogger with the specified Out io.Writer
func NewLogger(out io.Writer) (me *ObLogger) {
	me = &ObLogger{}
	if out != nil {
		me.logger = log.New(out, "", log.LstdFlags)
	}
	return
}

// Debugf formats its arguments according to the format, analogous to fmt.Printf,
// and records the text as a log message at Debug level
func (me *ObLogger) Debugf(format string, args ...interface{}) {
	if me.logger != nil {
		me.logger.Printf("[DEBUG]\t\t"+format+"\n", args...)
	}
}

// Infof is like Debugf, but at Info level
func (me *ObLogger) Infof(format string, args ...interface{}) {
	if me.logger != nil {
		me.logger.Printf("[INFO]\t\t"+format+"\n", args...)
	}
}

// Warningf is like Debugf, but at Warning level
func (me *ObLogger) Warningf(format string, args ...interface{}) {
	if me.logger != nil {
		me.logger.Printf("[WARNING]\t"+format+"\n", args...)
	}
}

// Error is like Debugf, but at Error level. Returns err.
func (me *ObLogger) Error(err error) error {
	me.Errorf(err.Error())
	return err
}

// Errorf is like Debugf, but at Error level
func (me *ObLogger) Errorf(format string, args ...interface{}) {
	if me.logger != nil {
		me.logger.Printf("[ERROR]\t\t"+format+"\n", args...)
	}
}

// Criticalf is like Debugf, but at Critical level
func (me *ObLogger) Criticalf(format string, args ...interface{}) {
	if me.logger != nil {
		me.logger.Printf("[CRITICAL]\t"+format+"\n", args...)
	}
}
