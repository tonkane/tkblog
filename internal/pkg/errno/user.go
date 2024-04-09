package errno

var (
	ErrUserAlreadyExist = &Errno{400, "FailedOperation.UserAlreadyExist", "User already exist."}

	ErrUserNotFound = &Errno{404, "ResourceNotFound.UserNotFound", "User not found."}

	ErrPwdIncorrect = &Errno{401, "InvalidParameter.PwdIncorrect", "Password was incorrect."}
)