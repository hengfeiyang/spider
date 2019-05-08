package util

import "fmt"

// Error 通用错误
type Error struct {
	code    int
	message string
}

func (e Error) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.code, e.message)
}

//Errorf 报错
func Errorf(format string, v ...interface{}) *Error {
	u := new(Error)
	u.message = fmt.Sprintf(format, v...)
	return u
}

// NewError 创建新的错误
func NewError(code int, message interface{}) *Error {
	u := new(Error)
	u.code = code
	switch message.(type) {
	case error:
		u.message = message.(error).Error()
	case string:
		u.message = message.(string)
	default:
		u.message = fmt.Sprintf("%v", message)
	}
	return u
}
