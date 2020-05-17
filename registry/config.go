package registry

import "time"

// Config 注册中心相关属性
// 字段
//   Addrs: 注册中心地址, ip:port
//   Timeout: 与注册中心的通信超时时间
//   RegistryPath: 在注册中心上的key, 一般采用的格式: /com.lzy/app/golang/service_A/127.0.0.1:8801
//   HeartBeat: 与注册中心的心跳, 也可以认为服务的租期, 如果注册中心与注册的机器超过此时间, 则RegistryPath(key)在注册中心会消失
type Config struct {
	Addrs        []string
	Timeout      time.Duration
	RegistryPath string
	HeartBeat    int64
}

// Option 用于修改Config字段
type Option func(cfg *Config)

// WithAddrs 修改Config.Addrs属性
func WithAddrs(Addrs []string) Option {
	return func(cfg *Config) {
		cfg.Addrs = Addrs
	}
}

// WithTimeout 修改Config.Timeout属性
func WithTimeout(Timeout time.Duration) Option {
	return func(cfg *Config) {
		cfg.Timeout = Timeout
	}
}

// WithRegistryPath 修改Config.RegistryPath属性
func WithRegistryPath(RegistryPath string) Option {
	return func(cfg *Config) {
		cfg.RegistryPath = RegistryPath
	}
}

// WithHeartBeat 修改Config.HeartBeat属性
func WithHeartBeat(HeartBeat int64) Option {
	return func(cfg *Config) {
		cfg.HeartBeat = HeartBeat
	}
}
