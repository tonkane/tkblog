package errno

import "fmt"

// 定义错误结构
type Errno struct {
	HTTP int
	Code string
	Message string
}

// 获取文字信息
func (err *Errno) Error() string{
	return err.Message
}

// 设置 message 信息
func (err *Errno) SetMessage(format string, args ...interface{}) *Errno {
	err.Message = fmt.Sprintf(format, args...)
	return err
}

// 从 err 中获取错误码和信息
func Decode(err error) (int, string, string) {
	if err == nil {
		return OK.HTTP, OK.Code, OK.Message
	}

	// err.(type) 断言机制
	switch typed := err.(type) {
	case *Errno:
		return typed.HTTP, typed.Code, typed.Message
	default:
	}

	return InternalServeError.HTTP, InternalServeError.Code, InternalServeError.Message
}

