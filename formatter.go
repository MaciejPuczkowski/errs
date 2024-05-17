package errs

import (
	"fmt"
	"strings"
)

type FileStackFormatter struct {
	separator string
}

func (f *FileStackFormatter) Format(errData []ErrorData) string {
	return f.formatList(errData)
}

func (f *FileStackFormatter) formatList(errors []ErrorData) string {
	strs := make([]string, 0)
	for _, err := range errors {
		fl := f.formatLine(&err)
		if fl != "" {
			strs = append(strs, fl)
		}
	}
	return strings.Join(strs, f.separator)
}

func (f *FileStackFormatter) formatLine(err *ErrorData) string {
	args := make([]string, 0)
	for k, v := range err.Args {
		args = append(args, f.formatArg(k, v))
	}
	fileLine := "(error getting file's line)"
	message := ""
	params := ""
	if err.HasValidTracking() {
		fileLine = fmt.Sprintf("%s:%d", err.TrackingData.FileName, err.TrackingData.Line)
	}
	lineTerms := []string{fileLine}
	if err.Message != "" {
		message = err.Message
		lineTerms = append(lineTerms, message)
	}
	if len(args) > 0 {
		params = fmt.Sprintf("params: %s", strings.Join(args, ", "))
		lineTerms = append(lineTerms, params)
	}
	return strings.Join(lineTerms, ": ")
}
func (f *FileStackFormatter) formatArg(k string, value any) string {
	return fmt.Sprintf("%s=%v", k, value)
}

func NewFileStackFormatter() *FileStackFormatter {
	return &FileStackFormatter{
		separator: "\n\t",
	}
}

type LogLineFormatter struct {
}

func NewLogLineFormatter() *LogLineFormatter {
	return &LogLineFormatter{}
}

func (f *LogLineFormatter) Format(errData []ErrorData) string {
	args := make(map[string]any)
	var last *ErrorData = nil
	for _, err := range errData {
		for k, v := range err.Args {
			if _, ok := args[k]; ok {
				args[k] = v
			} else {
				args[k] = v
			}
		}
		last = &err
	}
	argsStr := ""
	if len(args) > 0 {
		argsStr = fmt.Sprintf(": %s", f.formatArgs(args))
	}
	return fmt.Sprintf("%s%s", last.Message, argsStr)
}

func (f *LogLineFormatter) formatArgs(args map[string]any) string {
	strs := make([]string, 0)
	for k, v := range args {
		strs = append(strs, f.formatArg(k, v))
	}
	return strings.Join(strs, ", ")
}
func (f *LogLineFormatter) formatArg(k string, value any) string {
	return fmt.Sprintf("%s=%v", k, value)
}
