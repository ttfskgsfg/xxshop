package main

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"xxshop-api/user-web/global"
	"xxshop-api/user-web/initialize"
	"xxshop-api/user-web/utils"
	"xxshop-api/user-web/utils/register/consul"
	myvalidator "xxshop-api/user-web/validator"
)

func main() {
	//1、初始化logger
	initialize.InitLogger()
	//2、初始化配置文件
	initialize.InitConfig()

	//3、初始化router
	Router := initialize.Routers()

	//4、初始化翻译
	if err := initialize.InitTrans("zh"); err != nil {
		panic(err)
	}

	//5、初始化srv的连接
	initialize.InitSrvConn()

	viper.AutomaticEnv()
	//如果是开发环境 本地端口号固定，线上环境则自动获取端口号
	debug := viper.GetBool("CCSHOP_DEBUG")
	if !debug {
		port, err := utils.GetFreePort()
		if err == nil {
			global.Serverconfig.Port = port
		}
	}

	//注册验证器到gin
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("mobile", myvalidator.ValidateMobile) //自定义判断器
		_ = v.RegisterTranslation("mobile", global.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "{0} 非法的手机号码!", true) // see universal-translator for details
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
	}
	//服务注册
	serviceId := fmt.Sprintf("%s", uuid.New())
	register_client := consul.NewRegistryClient(global.Serverconfig.ConsulInfo.Host, global.Serverconfig.ConsulInfo.Port) //配置更加灵活 可以改使用etcd
	err := register_client.Register(global.Serverconfig.Host, global.Serverconfig.Port, global.Serverconfig.Name, global.Serverconfig.Tags, serviceId)
	if err != nil {
		zap.S().Panic("服务注册失败：", err.Error())
	}
	zap.S().Infof("启动服务器, 端口： %d", global.Serverconfig.Port)

	if err := Router.Run(fmt.Sprintf(":%d", global.Serverconfig.Port)); err != nil {
		zap.S().Panic("启动失败:", err.Error())
	}
	//接收终止信号 优雅退出
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	err = register_client.DeRegister(serviceId)
	if err = register_client.DeRegister(serviceId); err != nil {
		zap.S().Info("注销失败:", err.Error())
	} else {
		zap.S().Info("注销成功:")
	}
}
