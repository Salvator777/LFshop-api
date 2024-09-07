package global

import (
	"LFshop-api/user-web/config"
	"LFshop-api/user-web/proto"

	ut "github.com/go-playground/universal-translator"
)

// 定义一些全局变量

var (
	ServerConfig *config.ServerConfig = new(config.ServerConfig)

	Trans ut.Translator

	UserSrvClient proto.UserClient

	NacosConfig *config.NacosConfig = &config.NacosConfig{}
)
