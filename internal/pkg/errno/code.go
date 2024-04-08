package errno

var (
	// 定义OK时的结构
	OK = &Errno{HTTP: 200, Code: "", Message: ""}

	// 定义未知错误
	InternalServeError = &Errno{500, "InternalError", "Internal server error."}

	// 定义路由不匹配错误
	ErrPageNotFound = &Errno{404, "ResourceNotFound.PageNotFound", "Page not found."}

	// 参数绑定错误
	ErrBind = &Errno{400, "InvalidParameter.BindError", "Error occurred while binding the request body to the struct."}

	// 验证失败
	ErrInvalidParameter = &Errno{400, "InvalidParameter", "Parmeter verification failed."}
)