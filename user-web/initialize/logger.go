package initialize

import "go.uber.org/zap"

func InitLogger() {
	// 自定义全局的logger
	logger, _ := zap.NewProduction()
	zap.ReplaceGlobals(logger)
}
