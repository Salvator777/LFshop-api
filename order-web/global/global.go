package global

import (
	"LFshop-api/order-web/config"
	"LFshop-api/order-web/proto"

	ut "github.com/go-playground/universal-translator"
)

// 定义一些全局变量

var (
	ServerConfig *config.ServerConfig = new(config.ServerConfig)

	Trans ut.Translator

	OrderSrvClient proto.OrderClient

	NacosConfig *config.NacosConfig = &config.NacosConfig{}
)
