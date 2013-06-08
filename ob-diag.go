package obcore

import (
	"fmt"
	"io"
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

	// Errorf is like Debugf, but at Error level
	Errorf(format string, args ...interface{})

	// Criticalf is like Debugf, but at Critical level
	Criticalf(format string, args ...interface{})

	//	Fatal is like Criticalf, then panics
	Fatal(err error)
}

//	The canonical implementation of the Logger interface
type ObLogger struct {
	//	Unless nil, all logging methods write to Out
	Out io.Writer
}

//	Creates and returns a new ObLogger with the specified Out io.Writer
func NewLogger(out io.Writer) (me *ObLogger) {
	me = &ObLogger{Out: out}
	return
}

// Debugf formats its arguments according to the format, analogous to fmt.Printf,
// and records the text as a log message at Debug level
func (me *ObLogger) Debugf(format string, args ...interface{}) {
	if me.Out != nil {
		fmt.Fprintf(me.Out, "[DEBUG]\t\t"+format+"\n", args...)
	}
}

// Infof is like Debugf, but at Info level
func (me *ObLogger) Infof(format string, args ...interface{}) {
	if me.Out != nil {
		fmt.Fprintf(me.Out, "[INFO]\t\t"+format+"\n", args...)
	}
}

// Warningf is like Debugf, but at Warning level
func (me *ObLogger) Warningf(format string, args ...interface{}) {
	if me.Out != nil {
		fmt.Fprintf(me.Out, "[WARNING]\t"+format+"\n", args...)
	}
}

// Errorf is like Debugf, but at Error level
func (me *ObLogger) Errorf(format string, args ...interface{}) {
	if me.Out != nil {
		fmt.Fprintf(me.Out, "[ERROR]\t\t"+format+"\n", args...)
	}
}

// Criticalf is like Debugf, but at Critical level
func (me *ObLogger) Criticalf(format string, args ...interface{}) {
	if me.Out != nil {
		fmt.Fprintf(me.Out, "[CRITICAL]\t"+format+"\n", args...)
	}
}

//	Fatal is like Criticalf, then panics
func (me *ObLogger) Fatal(err error) {
	me.Criticalf("FATAL: %+v", err)
	panic(err)
}
