package errs

var FormatterFileStack Formatter = NewFileStackFormatter()
var FormatterLogLine Formatter = NewLogLineFormatter()
var factory = NewErrorFactory()
