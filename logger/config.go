package logger

var config = &Config{
	Type:       LogTypeConsole,
	Level:      LogLevelInfo,
	ModuleName: "service",
	ChanSize:   50000,
	FileSize:   1 << 20,
	LogDir:     "./log",
}

// Config 日志配置
// 字段
//   Type: 日志输出方式, 可以同时包含file和console, Type=file,console
//   Level: 日志级别
//   ModuleName: 使用日志的模块名称, 服务端程序一般使用服务名
//   ChanSize: 日志队列大小, 由于日志输出到控制台或文件时比较耗时, 先放到队列, 再在goroutine取出来处理
//   FileSize: 日志文件大小, 日志文件超过此大小时, 新建一个文件
//   LogDir: 日志所在路径, 输出到文件时需要
type Config struct {
	Type       LogType
	Level      LogLevel
	ModuleName string
	ChanSize   int
	FileSize   int64
	LogDir     string
}

// Option 设置Config字段的函数
type Option func(*Config)

// WithType 设置Config.Type
// 参数
//   typeStr: 日志类型, 请看constant定义
// 返回值
//   Option: 修改Type字段的函数
func WithType(typeStr string) Option {
	return func(cfg *Config) {
		cfg.Type = getLogType(typeStr)
	}
}

// WithLevel 设置Config.Level
// 参数
//   levelStr: 日志级别对应的字符串, 请看constant定义
// 返回值
//   Option: 修改Level字段的函数
func WithLevel(levelStr string) Option {
	return func(cfg *Config) {
		cfg.Level = getLogLevel(levelStr)
	}
}

// WithModuleName 设置Config.ModuleName
// 参数
//   module: 模块名称(服务名称)
// 返回值
//   Option: 修改ModuleName字段的函数
func WithModuleName(module string) Option {
	return func(cfg *Config) {
		cfg.ModuleName = module
	}
}

// WithChanSize 设置Config.ChanSize
// 参数
//   size: 队列大小
// 返回值
//   Option: 修改ChanSize字段的函数
func WithChanSize(size int) Option {
	return func(cfg *Config) {
		cfg.ChanSize = size
	}
}

// WithFileSize 设置Config.FileSize
// 参数
//   size: 文件大小
// 返回值
//   Option: 修改FileSize字段的函数
func WithFileSize(size int64) Option {
	return func(cfg *Config) {
		cfg.FileSize = size
	}
}

// WithLogDir 设置Config.LogPath
// 参数
//   dir: 日志目录
// 返回值
//   Option: 修改LogDir字段的函数
func WithLogDir(dir string) Option {
	return func(cfg *Config) {
		cfg.LogDir = dir
	}
}
