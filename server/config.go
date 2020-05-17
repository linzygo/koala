package server

import (
	"fmt"
	"koala/config"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

var koalaConf = &KoalaConfig{
	Service: ServiceConfig{
		Name: "service",
		Port: 8081,
	},
}

// KoalaConfig Koala服务端配置
// 字段
//   ServiceConf: 服务相关的配置
//   Log: 日志配置
//   Register: 服务注册配置
//   Prometheus: 与Prometheus集成相关的配置
//   Limit: 限流器配置
//   Trace: 分布式链路追踪配置
type KoalaConfig struct {
	Service    ServiceConfig           `toml:"service"`
	Log        config.LogConfig        `toml:"log"`
	Register   config.RegisterConfig   `toml:"register"`
	Prometheus config.PrometheusConfig `toml:"prometheus"`
	Limit      config.LimiterConfig    `toml:"limit"`
	Trace      config.TraceConfig      `toml:"trace"`
}

// ServiceConfig 服务的相关配置
// 字段
//   Name: 服务名称
//   Port: 服务监听的端口
type ServiceConfig struct {
	Name string `toml:"name"`
	Port int    `toml:"port"`
}

// initConfig 读取配置文件初始化配置
// 参数
//   service: 服务名称
// 返回值
//   err: error
func initConfig(service string) (err error) {
	appRootDir := filepath.Dir(os.Args[0])
	fpath := filepath.Join(appRootDir, "..", "conf", fmt.Sprintf("%s.toml", service))
	_, err = toml.DecodeFile(fpath, koalaConf)
	if err != nil {
		fmt.Printf("加载配置文件失败, err=%v\n", err)
		return
	}
	fmt.Printf("读取配置成功, 配置文件[%s]\n", fpath)
	return
}
