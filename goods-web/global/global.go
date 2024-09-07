package global

import (
	"LFshop-api/goods-web/config"
	"LFshop-api/goods-web/proto"

	ut "github.com/go-playground/universal-translator"
)

// 定义一些全局变量

var (
	ServerConfig *config.ServerConfig = new(config.ServerConfig)

	Trans ut.Translator

	GoodsSrvClient proto.GoodsClient

	NacosConfig *config.NacosConfig = &config.NacosConfig{}
)
