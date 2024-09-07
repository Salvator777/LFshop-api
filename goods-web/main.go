package main

import (
	"LFshop-api/goods-web/global"
	"LFshop-api/goods-web/initialize"
	"LFshop-api/goods-web/utils"
	"LFshop-api/goods-web/utils/register/consul"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	// 1.初始化全局logger
	initialize.InitLogger()

	// 2.初始化config
	initialize.InitConfig()

	// 3.初始化router
	Router := initialize.Routers()

	// 4.初始化翻译器
	err := initialize.InitTrans("zh")
	if err != nil {
		zap.S().Warnf("翻译器初始化错误：%s", err.Error())
	}

	// 5.初始化用户服务user_srv
	initialize.InitSrvConn()

	// 6.初始化sentinel
	initialize.InitSentinel()

	//如果是本地开发环境端口号固定，线上环境启动获取端口号
	debug := viper.GetBool("LFSHOP_DEBUG")
	if !debug {
		port, err := utils.GetFreePort()
		if err == nil {
			global.ServerConfig.Port = port
		}
	}

	register_client := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	serviceId := uuid.NewString()
	err = register_client.Register(global.ServerConfig.Host, global.ServerConfig.Port, global.ServerConfig.Name, global.ServerConfig.Tags, serviceId)
	if err != nil {
		zap.S().Panic("goods-web服务注册失败:", err.Error())
	}

	// 通过zap.S()使用全局的logger的sugar
	// 并发安全，这个全局的sugar都可用
	zap.S().Info("启动服务器，端口：", global.ServerConfig.Port)

	go func() {
		if err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
			zap.S().Panic("启动失败:", err.Error())
		}
	}()
	//接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err = register_client.DeRegister(serviceId); err != nil {
		zap.S().Info("注销失败:", err.Error())
	} else {
		zap.S().Info("注销成功")
	}
}
