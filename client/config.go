package client

import (
	"fmt"
	"koala/config"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

var koalaConf = &KoalaClientConfig{}

// KoalaClientConfig client程序的rpc设置
// 字段
//   Log: 日志配置
//   Discovery: 服务发现配置
//   Prometheus: Prometheus集成配置
//   Limit: 限流器配置
//   Trace: 分布式链路追踪配置
//   Balancer: 负载均衡
type KoalaClientConfig struct {
	Log        config.LogConfig        `toml:"log"`
	Discovery  config.RegisterConfig   `toml:"discovery"`
	Prometheus config.PrometheusConfig `toml:"prometheus"`
	Limit      config.LimiterConfig    `toml:"limit"`
	Trace      config.TraceConfig      `toml:"trace"`
	Balancer   LoadBalanceConfig       `toml:"loadbalance"`
}

// LoadBalanceConfig 负载均衡配置
// 字段
//   Name: 名称
type LoadBalanceConfig struct {
	Name string `toml:"name"`
}

func initConfig() (err error) {
	appRootDir := filepath.Dir(os.Args[0])
	fpath := filepath.Join(appRootDir, "..", "conf", "rpc_client.toml")
	_, err = toml.DecodeFile(fpath, koalaConf)
	if err != nil {
		fmt.Printf("加载配置文件失败, err=%v\n", err)
		return
	}
	fmt.Printf("读取配置成功, 配置文件[%s]\n", fpath)
	return
}
