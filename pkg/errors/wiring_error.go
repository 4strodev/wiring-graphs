package errors

import "fmt"

type wiringErrorCode int

const (
	E_CIRCULAR_DEPENDENCY wiringErrorCode = iota
	E_INVALID_RESOLVER
	E_REDECLARED_DEPENDENCY
	E_DEPENDENCY_NOT_FOUND
	E_TYPE_ERROR
)

type WiringError struct {
	err  error
	code wiringErrorCode
}

func Errorf(code wiringErrorCode, format string, values ...any) *WiringError {
	error := fmt.Errorf(format, values...)
	return &WiringError{
		err:  error,
		code: code,
	}
}

func (e *WiringError) Error() string {
	return e.err.Error()
}

func (e *WiringError) Unwrap() error {
	return e.err
}

func (e WiringError) Code() wiringErrorCode {
	return e.code
}
