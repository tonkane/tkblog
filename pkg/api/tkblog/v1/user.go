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

type GetUserResponse UserInfo

type UserInfo struct {
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	PostCount int64 `json:"postCount"`
	CreateAt string `json:"createAt"`
	UpdateAt string `json:"updateAt"`
}

// user list 的请求参数结构
// 这里的 form ？表单类型？
type ListUserRequest struct {
	Offset int `form:"offset"`
	Limit int `form:"limit"`
}

// user list 的响应参数结构
type ListUserResponse struct {
	TotalCount int64 `json:"totalCount"`
	Users []*UserInfo `json:"users"`
}