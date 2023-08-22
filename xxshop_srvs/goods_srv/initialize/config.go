package initialize

import (
	"encoding/json"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"xxshop_srvs/goods_srv/global"
)

// 将生成环境和开发环境隔离开来
func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
	//刚才设置的环境变量 想要生效 我们必须得重启goland
}

func InitConfig() {
	//从配置文件中读取对应配置
	debug := GetEnvInfo("XXSHOP_DEBUG")
	configFilePrefix := "config"
	configFileName := fmt.Sprintf("goods_srv/%s-pro.yaml", configFilePrefix)
	if debug {
		configFileName = fmt.Sprintf("goods_srv/%s-debug.yaml", configFilePrefix)
	}
	v := viper.New()
	//文件的路径如何设置
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	//定义的时候不是实例化指针类型 所以要加取地址符号
	if err := v.Unmarshal(&global.Nacosconfig); err != nil {
		panic(err)
	}
	zap.S().Infof("配置信息: &v", global.Nacosconfig)
	//从nacos中读取信息
	sc := []constant.ServerConfig{
		{
			IpAddr: global.Nacosconfig.Host,
			Port:   global.Nacosconfig.Port,
		},
	}

	cc := constant.ClientConfig{
		NamespaceId:         global.Nacosconfig.Namespace, // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "tmp/nacos/log",
		CacheDir:            "tmp/nacos/cache",
		//RotateTime:          "1h",
		//MaxAge:              3,
		LogLevel: "debug",
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		panic(err)
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: global.Nacosconfig.DataId,
		Group:  global.Nacosconfig.Group})

	if err != nil {
		panic(err)
	}
	//fmt.Println(content) //字符串 - yaml
	//局部变量 出错原因
	//serverConfig := config.ServerConfig{}
	//想要将一个json字符串转换成struct，需要去设置这个struct的tag
	err = json.Unmarshal([]byte(content), &global.ServerConfig)
	if err != nil {
		zap.S().Fatalf("读取nacos配置失败 ： %s", err.Error())
	}
	fmt.Println(global.ServerConfig)
}
