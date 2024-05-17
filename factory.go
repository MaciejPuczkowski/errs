package errs

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
)

// ErrCallerFailed is returned in TrackingData.Err when runtime.Caller regurn false on 4th element.
// It means that for some reason it was not possible to get the caller and tracking data may not be
// valid.
var ErrCallerFailed = errors.New("failed to get caller")

// New creates a new error with the given message and args.
func New(t string, args ...any) *Error {
	return factory.newErrorTr(trackingData(), t, args...)
}

// Wrap wraps the given error into a new error. It also makes a tracking point for the error's chain.
// Error returned directly from the function is transparent. It is enough to do return Wrap(err), instead
// of return err, to track the path of how the error was returned. It also creates a new object that envelopes
// the given one.
func Wrap(err error) *Error {
	return factory.wrapTr(trackingData(), err)
}

// SetFormatter sets the default formatter for the errors.
func SetFormatter(formatter Formatter) {
	factory.formatter = formatter
}

func trackingData() *TrackingData {
	var tr TrackingData
	tr.PC, tr.FileName, tr.Line, tr.IsValid = runtime.Caller(2)
	if !tr.IsValid {
		tr.Err = ErrCallerFailed
	} else {
		basepath, err := os.Getwd()
		if err != nil {
			tr.Err = err
		} else {
			tr.FileName, err = filepath.Rel(basepath, tr.FileName)
			if err != nil {
				tr.Err = err
			}
		}
	}
	return &tr
}

// TrackingData is a data struct that holds the data of a tracking point.
type TrackingData struct {
	PC       uintptr
	FileName string
	Line     int
	IsValid  bool
	Err      error
}

// Formatter is an interface for formatting errors. A client may create a custome formatter
// for the errors. errData contains chain of wrapped errors from top to bottom.
// The last element is the initial error.
type Formatter interface {
	Format(errData []ErrorData) string
}

// ErrorFactory is a factory for creating errors. A client doesn't have to use it directly.
// functions New and Wrap uses default factory. It could be useful to create own factory
// if a client wants to use more than one custom formatter at the same time.
// If a client wants to change formatter for the default factory it is enough to use SetFormatter
// function.
type ErrorFactory struct {
	formatter Formatter
}

// NewErrorFactory creates a new factory with default formatter (FileStackFormatter).
func NewErrorFactory() *ErrorFactory {
	return &ErrorFactory{
		formatter: NewFileStackFormatter(),
	}
}

// SetFormatter sets a formatter for the errors.
func (ef *ErrorFactory) SetFormatter(t Formatter) {
	ef.formatter = t
}

// NewError creates a new error with the given message and args.
func (ef *ErrorFactory) NewError(t string, args ...any) *Error {
	var tr TrackingData
	tr.PC, tr.FileName, tr.Line, tr.IsValid = runtime.Caller(1)
	return ef.newErrorTr(&tr, t, args...)
}

func (ef *ErrorFactory) newErrorTr(tr *TrackingData, t string, args ...any) *Error {
	err := &Error{
		msg:          t,
		msgArgs:      args,
		args:         map[string]any{},
		trackingData: tr,
		formatter:    ef.formatter,
	}
	return err
}

// Wrap wraps the given error into a new error. It also makes a tracking point for the error's chain.
func (ef *ErrorFactory) Wrap(err error) *Error {
	var tr TrackingData
	tr.PC, tr.FileName, tr.Line, tr.IsValid = runtime.Caller(1)
	return ef.wrapTr(&tr, err)
}

func (ef *ErrorFactory) wrapTr(tr *TrackingData, err error) *Error {
	msg := ""
	var inner *Error = nil
	switch e := err.(type) {
	case Error:
		inner = &e
	case *Error:
		inner = e
	default:
		msg = e.Error()
	}

	newErr := &Error{
		msg:          msg,
		inner:        inner,
		args:         map[string]any{},
		trackingData: tr,
		formatter:    ef.formatter,
	}
	return newErr
}
