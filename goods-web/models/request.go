package models

import (
	"github.com/dgrijalva/jwt-go"
)

// 这里是payload部分
type CustomClaims struct {
	ID          uint
	NickName    string
	AuthorityId uint // 这个表示权限等级，1是普通用户，2是管理员
	jwt.StandardClaims
}
