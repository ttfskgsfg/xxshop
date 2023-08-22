package initialize

import "go.uber.org/zap"

func InitLogger() {
	logger, _ := zap.NewDevelopment() //全局使用zap的信息
	zap.ReplaceGlobals(logger)
}
