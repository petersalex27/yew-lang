package errors

import "fmt"

const SystemErrorHead string = " Error (System): "

type sysError string

func MkSystemError(msg string, args ...any) sysError {
	return sysError(fmt.Sprintf(msg, args...))
}

func (e sysError) Error() string {
	return SystemErrorHead + string(e)
}