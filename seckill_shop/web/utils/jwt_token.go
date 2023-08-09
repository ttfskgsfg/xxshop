package utils

import (
	"github.com/dgrijalva/jwt-go"
	"time"
	"github.com/pkg/errors"
)

type UserToken struct {
	jwt.StandardClaims
	// 自定义的用户信息
	UserName string `json:"user_name"`

}

// 前端用户token过期时间
var FrontUserExpireDuration = time.Hour
var FrontUserSecretKey = []byte("front_user_token")


// 管理端用户token过期时间
var AdminUserExpireDuration = time.Hour * 2
var AdminUserSecretKey = []byte("admin_user_token")


// 生成token
func GenToken(UserName string,expireDuration time.Duration,secret_key []byte)  (string,error){

	user := UserToken{
		jwt.StandardClaims{
			// 现在 + 加上传的过期时间
			ExpiresAt:time.Now().Add(expireDuration).Unix(),
			Issuer:"micro_gin_vue",
		},
		UserName,

	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,user)
	return token.SignedString(secret_key)

}

// 认证token
func AuthToken(tokenString string,secretKey []byte) (*UserToken, error){


	// 解析token
	token,err := jwt.ParseWithClaims(tokenString,&UserToken{}, func(token *jwt.Token) (key interface{}, err error) {
		return secretKey,nil
	})

	if err != nil {
		return nil,err
	}

	clasims,is_ok := token.Claims.(*UserToken)

	// 验证token
	if is_ok && token.Valid { // 正常的
		return clasims,nil
	}

	return nil,errors.New("token valid err")

}