package config

// LogConfig 日志的相关配置
// 字段
//   Type: 日志输出类别, console、file或console,file
//   Level: 输出的日志级别
//   Path: 文件日志根目录
//   FileSize: 文件日志文件大小
//   ChanSize: 日志队列大小
type LogConfig struct {
	Type     string `toml:"type"`
	Level    string `toml:"level"`
	Path     string `toml:"path"`
	FileSize int64  `toml:"file_size"`
	ChanSize int    `toml:"chan_size"`
}

// RegisterConfig 注册中心的相关配置
// 字段
//   SwitchOn: 开关
//   RegisterName: 使用的注册中心名称, 例如etcd或consul
//   RegisterAddrs: 注册中心的地址
//   RegisterPath: 在注册中心的key
//   Timeout: 与注册中心超时时间
//   Heartbeat: 服务租期
type RegisterConfig struct {
	SwitchOn      bool     `toml:"switch_on"`
	RegisterName  string   `toml:"register_name"`
	RegisterAddrs []string `toml:"register_addrs"`
	RegisterPath  string   `toml:"register_path"`
	Timeout       int64    `toml:"timeout"`
	Heartbeat     int64    `toml:"heartbeat"`
}

// PrometheusConfig 集成Prometheus的相关配置
// 字段
//   Port: 提供给Prometheus访问的http端口
//   SwitchOn: 开关, =true时表示打开, 即启用, =false表示关闭
type PrometheusConfig struct {
	Port     int  `toml:"port"`
	SwitchOn bool `toml:"switch_on"`
}

// LimiterConfig 限流器配置
// 字段
//   QPS: QPS上限
//   SwitchOn: 开关
type LimiterConfig struct {
	QPS      int  `toml:"qps"`
	SwitchOn bool `toml:"switch_on"`
}

// TraceConfig 分布式链追踪设置
// 字段
//   SwitchOn: 开关
//   ReportAddr: 追踪系统接口url
//   SampleType: 采样类型
//   SampleRate: 采样率
type TraceConfig struct {
	SwitchOn   bool   `toml:"switch_on"`
	ReportAddr string `toml:"report_addr"`
	SampleType string `toml:"sample_type"`
	SampleRate int    `toml:"sample_rate"`
}
