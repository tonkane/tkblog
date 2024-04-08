package v1

// 指定请求接收的参数
type CreateUserRequest struct {
	Username string `json:"username" vaild:"alphanum,required,stringlength(1|50)"`
	Password string `json:"passwprd" vaild:"required,stringlength(6|18)"`
	Nickname string `json:"nickname" vaild:"required,stringlength(1|50)"`
	Email string `json:"email" vaild:"required,email"`
	Phone string `json:"phone" vaild:"required,stringlength(11|11)"`
}