package global

import (
	ut "github.com/go-playground/universal-translator"
	config "xxshop-api/user-web/config"
	"xxshop-api/user-web/proto"
)

var (
	//要在其他地方使用 所以是指针类型
	Serverconfig *config.ServerConfig = &config.ServerConfig{}

	Nacosconfig *config.NacosConfig = &config.NacosConfig{}

	Trans ut.Translator

	UserSrvClient proto.UserClient //调用srv端传来的接口
)
