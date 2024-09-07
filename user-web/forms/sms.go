package forms

// 验证码在业务中有很多地方会用到，需要区别不同的业务：比如登录验证码、找回密码验证码等
// 1表示注册，2表示找回密码
type SendSmsForm struct {
	Phone string `form:"Phone" json:"Phone" binding:"required,Phone"` //手机号码格式有规范可寻， 自定义validator
	Type  uint   `form:"type" json:"type" binding:"required,oneof=1 2"`
}
