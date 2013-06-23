package obcore

import (
	"io"
	"log"
)

//	An interface for log output. `NewLogger` provides the canonical implementation.
type Logger interface {
	// `Debugf` formats its arguments according to the `format`, analogous to `fmt.Printf`,
	// and records the text as a log message at Debug level.
	Debugf(format string, args ...interface{})

	// `Infof` is like `Debugf`, but at Info level.
	Infof(format string, args ...interface{})

	// `Warningf` is like `Debugf`, but at Warning level.
	Warningf(format string, args ...interface{})

	// `Error` records the specified `error` message at Error level,
	//	then should return the same specified `error` for more convenient in-place handling.
	Error(error) error

	// `Errorf` is like `Debugf`, but at Error level.
	Errorf(format string, args ...interface{})

	// `Criticalf` is like `Debugf`, but at Critical level.
	Criticalf(format string, args ...interface{})
}

//	Creates and returns a new `Logger`; `out` is optional and if `nil`, this disables logging.
func NewLogger(out io.Writer) Logger {
	var me logger
	if out != nil {
		me.Logger = log.New(out, "", log.LstdFlags)
	}
	return &me
}

//	The canonical implementation of the `Logger` interface, using a standard `log.Logger`.
type logger struct {
	*log.Logger
}

// Implements `Logger` interface.
func (me *logger) Debugf(format string, args ...interface{}) {
	if me.Logger != nil {
		me.Logger.Printf("[DEBUG]\t\t"+format+"\n", args...)
	}
}

// Implements `Logger` interface.
func (me *logger) Infof(format string, args ...interface{}) {
	if me.Logger != nil {
		me.Logger.Printf("[INFO]\t\t"+format+"\n", args...)
	}
}

// Implements `Logger` interface.
func (me *logger) Warningf(format string, args ...interface{}) {
	if me.Logger != nil {
		me.Logger.Printf("[WARNING]\t"+format+"\n", args...)
	}
}

// Implements `Logger` interface.
func (me *logger) Error(err error) error {
	me.Errorf(err.Error())
	return err
}

// Implements `Logger` interface.
func (me *logger) Errorf(format string, args ...interface{}) {
	if me.Logger != nil {
		me.Logger.Printf("[ERROR]\t\t"+format+"\n", args...)
	}
}

// Implements `Logger` interface.
func (me *logger) Criticalf(format string, args ...interface{}) {
	if me.Logger != nil {
		me.Logger.Printf("[CRITICAL]\t"+format+"\n", args...)
	}
}
