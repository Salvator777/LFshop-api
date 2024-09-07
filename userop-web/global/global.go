package global

import (
	"LFshop-api/userop-web/config"
	"LFshop-api/userop-web/proto"

	ut "github.com/go-playground/universal-translator"
)

// 定义一些全局变量

var (
	Trans ut.Translator

	ServerConfig *config.ServerConfig = &config.ServerConfig{}

	NacosConfig *config.NacosConfig = &config.NacosConfig{}

	MessageClient proto.MessageClient

	AddressClient proto.AddressClient

	UserFavClient proto.UserFavClient
)
