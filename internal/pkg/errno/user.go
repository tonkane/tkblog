package errno

var (
	ErrUserAlreadyExist = &Errno{400, "FailedOperation.UserAlreadyExist", "User already exist."}
)