package contracts

// Login 登陆
type Login struct {
	Account  string `json:"account" form:"account"`   // 账号
	Password string `json:"password" form:"password"` // 密码
}

// Register 注册信息
type Register struct {
	PerfectInfo
	Account         string `json:"account" form:"account" binding:"required,email"`             // 账号
	Password        string `json:"password" form:"password" binding:"required"`                 // 密码
	ConfirmPassword string `json:"confirm_password" form:"confirm_password" binding:"required"` // 确认密码
}

// PerfectInfo 完善信息
type PerfectInfo struct {
	Phone  string `json:"phone" form:"phone"`
	Name   string `json:"name" form:"name" binding:"required"`
	Avatar string `json:"avatar" form:"avatar"`
}

// LoginSuccess 登陆成功
type LoginSuccess struct {
	PerfectInfo
	Account string   `json:"account" ` // 账号
	Auth    []string `json:"auth"`     // 权限列表
	Token   string   `json:"token"`
}
