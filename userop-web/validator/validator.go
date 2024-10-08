package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

// 验证手机号，使用正则表达式判断是否合法
func ValidatePhone(fl validator.FieldLevel) bool {
	Phone := fl.Field().String()
	ok, _ := regexp.MatchString(`^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`, Phone)
	if !ok {
		return false
	}
	return true
}
