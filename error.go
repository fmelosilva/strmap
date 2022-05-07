package strmap

import (
	"fmt"
)

type ParserError struct {
	msg string
}

func (m *ParserError) Error() string {
	return m.msg
}

func ParseErrorf(format string, args ...interface{}) *ParserError {
	return &ParserError{
		msg: fmt.Sprintf(format, args...),
	}
}

type OverflowError struct {
	msg string
}

func (m *OverflowError) Error() string {
	return m.msg
}

func OverflowErrorf(format string, args ...interface{}) *OverflowError {
	return &OverflowError{
		msg: fmt.Sprintf(format, args...),
	}
}
