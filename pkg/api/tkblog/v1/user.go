package v1

// 指定请求接收的参数
type CreateUserRequest struct {
	Username string `json:"username" vaild:"alphanum,required,stringlength(1|50)"`
	Password string `json:"password" vaild:"required,stringlength(6|18)"`
	Nickname string `json:"nickname" vaild:"required,stringlength(1|50)"`
	Email string `json:"email" vaild:"required,email"`
	Phone string `json:"phone" vaild:"required,stringlength(11|11)"`
}

// login request
type LoginRequest struct {
	Username string `json:"username" valid:"alphanum,required,stringlength(1|50)"`
	Password string `json:"password" valid:"required,stringlength(6|18)"`
}

// login response
type LoginResponse struct {
	Token string `json:"token"`
}

// change pwd request
type ChangePwdRequest struct {
	OldPwd string `json:"oldpwd" valid:"required,stringlength(6|18)"`
	NewPwd string `json:"newpwd" valid:"required,stringlength(6|18)"`
}