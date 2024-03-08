// apiserver服务的错误码 以11打头
package code

// 员工模块00
const (
	// 找不到用户
	ErrUserNotFound int = iota + 110001
	// 用户已存在
	ErrUserAlreadyExist
)

// 医生模块00
