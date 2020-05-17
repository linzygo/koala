package middleware

// TraceConfig 追踪设置
// 字段
//   AppName: 要在追踪系统上显示的名字
//   ReportAddr: Jaeger接口url
//   SampleType: 采样类型
//   SampleRate: 采样率
type TraceConfig struct {
	AppName    string
	ReportAddr string
	SampleType string
	SampleRate int
}

// TraceOption 修改TraceConfig字段
type TraceOption func(cfg *TraceConfig)

// WithAppName 修改TraceConfig.AppName
func WithAppName(service string) TraceOption {
	return func(cfg *TraceConfig) {
		cfg.AppName = service
	}
}

// WithReportAddr 修改TraceConfig.ReportAddr
func WithReportAddr(addr string) TraceOption {
	return func(cfg *TraceConfig) {
		cfg.ReportAddr = addr
	}
}

// WithSampleType 修改TraceConfig.SampleType
func WithSampleType(sampleType string) TraceOption {
	return func(cfg *TraceConfig) {
		cfg.SampleType = sampleType
	}
}

// WithSampleRate 修改TraceConfig.SampleRate
func WithSampleRate(sampleRate int) TraceOption {
	return func(cfg *TraceConfig) {
		cfg.SampleRate = sampleRate
	}
}
