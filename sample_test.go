package errs

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func init() {
	factory.SetFormatter(FormatterLogLine)
}

type formattingTest struct {
	formatter Formatter
	expected  string
}

func test(t *testing.T, err error, formattingTests []formattingTest) {
	for _, f := range formattingTests {
		err = err.(*Error).WithFormatter(f.formatter)
		if err == nil || err.Error() != f.expected {
			t.Errorf("For formatter %v expected test, got %v", reflect.TypeOf(f.formatter).Elem().Name(), err)
		}
	}
}

func Test_Sample1(t *testing.T) {
	f1 := func(s string) error {
		return New(s)
	}
	err := f1("test")
	formattingTests := []formattingTest{
		{
			formatter: FormatterLogLine,
			expected:  "test",
		},
	}
	test(t, err, formattingTests)
}

func Test_Sample2(t *testing.T) {
	f1 := func() error {
		return New("test %s %d", "arg", 1)
	}
	err := f1()
	if err == nil || err.Error() != "test arg 1" {
		t.Errorf("Expected test, got %v", err)
	}
}

func Test_Sample3(t *testing.T) {
	f1 := func() error {
		return New("test %s %d", "arg", 1)
	}
	err := f1()
	if err == nil || err.Error() != "test arg 1" {
		t.Errorf("Expected test, got %v", err)
	}
}

func Test_Sample4(t *testing.T) {
	var err error
	err = New("test %s %d", "arg", 1)
	err = Wrap(err)
	err = Wrap(err)
	if err.Error() != "test arg 1" {
		t.Errorf("Expected test, got \n%v", err)
	}
}

func Test_Sample5(t *testing.T) {
	f1 := func() error {
		return New("test %s %d", "arg", 1)
	}
	err := f1()
	err = Wrap(err).Msg("new %d", 2)
	err = Wrap(err)
	if err.Error() != "test arg 1" {
		t.Errorf("Expected test, got \n%v", err)
	}
}
func Test_Sample6(t *testing.T) {
	f1 := func() error {
		return fmt.Errorf("test")
	}
	err := f1()
	err = Wrap(err)
	err = Wrap(err).Msg("new %d", 2)
	err = Wrap(err).Arg("arg", "yes")
	if err.Error() != "test: arg=yes" {
		t.Errorf("Expected test, got %v", err)
	}
}

func Test_Sample7(t *testing.T) {
	f1 := func() error {
		return fmt.Errorf("test")
	}
	err := f1()
	err = Wrap(err)
	err = Wrap(err).Msg("new %d", 2).Arg("arg", "no")
	err = Wrap(err).Arg("arg", "yes")
	if err.Error() != "test: arg=no" {
		t.Errorf("Expected test, got %v", err)
	}
}

func Test_Sample8(t *testing.T) {
	f1 := func() error {
		return fmt.Errorf("test")
	}
	err := f1()
	err = Wrap(err)
	err = Wrap(err).Msg("new %d", 2).Arg("arg", "no")
	err = Wrap(err).Arg("arg", "yes")
	if err.Error() != "test: arg=no" {
		t.Errorf("Expected test, got %v", err)
	}
}

func Test_errors_is_should_work(t *testing.T) {
	var predefinedErr error = errors.New("predefined")
	f1 := func() error {
		return Wrap(predefinedErr)
	}
	err := f1()
	err = Wrap(err)
	err = Wrap(err).Msg("new %d", 2)
	err = Wrap(err).Arg("arg", "yes")
	if errors.Is(err, predefinedErr) {
		t.Errorf("Expected test, got %v", err)
	}
}

func Test_errors_is_should_work2(t *testing.T) {
	var predefinedErr error = errors.New("predefined")
	var predefinedErr2 error = Wrap(predefinedErr)
	f1 := func() error {
		return Wrap(predefinedErr2)
	}
	err := f1()
	err = Wrap(err)
	err = Wrap(err).Msg("new %d", 2)
	err = Wrap(err).Arg("arg", "yes")
	if errors.Is(err, predefinedErr2) {
		t.Errorf("Expected test, got %v", err)
	}
}

type TError struct {
	msg string
}

func (e *TError) Error() string {
	return e.msg
}
func (e *TError) Error2() string {
	return e.msg
}

type ITError interface {
	Error2() string
}

func Test_errors_as_should_work2(t *testing.T) {
	var predefinedErr error = &TError{msg: "predefined"}
	f1 := func() error {
		return Wrap(predefinedErr)
	}
	err := f1()
	err = Wrap(err)
	err = Wrap(err).Msg("new %d", 2)
	err = Wrap(err).Arg("arg", "yes")
	var te ITError
	if errors.As(err, &te) {
		t.Errorf("Expected test, got %v", err)
	}
}
