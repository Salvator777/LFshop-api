package main

import (
	"LFshop-api/order-web/global"
	"LFshop-api/order-web/initialize"
	"LFshop-api/order-web/utils"
	"LFshop-api/order-web/utils/register/consul"
	myvalidator "LFshop-api/user-web/validator"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
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

	// 5.初始化其他服务
	initialize.InitSrvConn()

	// 注册自定义的手机验证器，和对应的报错信息
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("Phone", myvalidator.ValidatePhone)
		_ = v.RegisterTranslation("Phone", global.Trans, func(ut ut.Translator) error {
			return ut.Add("Phone", "{0} 非法的手机号码!", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("Phone", fe.Field())
			return t
		})
	}

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
		zap.S().Panic("order-web服务注册失败:", err.Error())
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
