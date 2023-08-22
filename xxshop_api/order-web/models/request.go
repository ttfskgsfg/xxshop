package models

import (
	"github.com/dgrijalva/jwt-go"
)

// 自定义的struct
type CustomClaims struct {
	ID          uint
	NickName    string
	AuthorityId uint //不同角色设计
	jwt.StandardClaims
}
