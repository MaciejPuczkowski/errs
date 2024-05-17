package errs

import (
	"fmt"
)

// ErrorData is a data struct that holds the data of an error.
type ErrorData struct {
	TrackingData *TrackingData
	Message      string
	Args         map[string]any
}

// HasValidTracking checks if the ErrorData's TrackingData is not nil and IsValid is true.
//
// Returns:
// - bool: true if TrackingData is not nil and IsValid is true, false otherwise.
func (e ErrorData) HasValidTracking() bool {
	return e.TrackingData != nil && e.TrackingData.IsValid
}

// Error is a struct that holds the data of an error.
type Error struct {
	trackingData *TrackingData
	formatter    Formatter
	sep          string
	inner        *Error
	msg          string
	msgArgs      []any
	args         map[string]any
}

// Error implements error interface.
func (e Error) Error() string {
	errors := e.collect()
	return e.formatter.Format(errors)
}

// Msg sets the message of the error. May be used for additional information of the tracking point.
// eg. return Wrap(err).Msg("it happened during processing %s", "data")
func (e *Error) Msg(msg string, args ...any) *Error {
	e.msg = fmt.Sprintf(msg, args...)
	return e
}

// Arg sets an extra argument of the error. May be used to pass the parameters of error's context.
// eg. return Wrap(err).Msg("it happened during processing %s", "data").Arg("methodArg1", arg1)
func (e *Error) Arg(name string, arg any) *Error {
	e.args[name] = arg
	return e
}

func (e *Error) collect() []ErrorData {
	res := []ErrorData{{
		TrackingData: e.trackingData,
		Message:      fmt.Sprintf(e.msg, e.msgArgs...),
		Args:         e.args,
	}}
	if e.inner != nil {
		return append(res, e.inner.collect()...)
	}
	return res
}

// WithFormatter sets a formatter for the error. It copies all erros in the chain with the given formatter.
func (e *Error) WithFormatter(formatter Formatter) *Error {
	var inner *Error
	if e.inner != nil {
		inner = e.inner.WithFormatter(formatter)
	}
	return &Error{
		trackingData: e.trackingData,
		formatter:    formatter,
		sep:          e.sep,
		inner:        inner,
		msg:          e.msg,
		msgArgs:      e.msgArgs,
		args:         e.args,
	}
}
