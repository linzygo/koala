package prometheus

import (
	"context"

	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc/status"
)

// Metrics 服务端采样打点
// 字段
//   requestCounter: 请求数量
//   errCounter: 请求错误数
//   costSummary: 请求耗时
type Metrics struct {
	requestCounter *prom.CounterVec
	errCounter     *prom.CounterVec
	costSummary    *prom.SummaryVec
}

// NewServerMetrics 生成服务器metrics实例
func NewServerMetrics() *Metrics {
	return &Metrics{
		requestCounter: promauto.NewCounterVec(
			prom.CounterOpts{
				Name: "koala_server_request_total",
				Help: "请求数量",
			},
			[]string{"service", "method"},
		),
		errCounter: promauto.NewCounterVec(
			prom.CounterOpts{
				Name: "koala_server_handle_err_total",
				Help: "请求错误数",
			},
			[]string{"service", "method", "grpc_code"},
		),
		costSummary: promauto.NewSummaryVec(
			prom.SummaryOpts{
				Name:       "koala_proc_cost",
				Help:       "请求耗时",
				Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
			},
			[]string{"service", "method"},
		),
	}
}

// NewRPCMetrics 生成rpc客户端调用metrics实例
func NewRPCMetrics() *Metrics {
	return &Metrics{
		requestCounter: promauto.NewCounterVec(
			prom.CounterOpts{
				Name: "koala_rpc_call_total",
				Help: "rpc调用服务数",
			},
			[]string{"service", "method"},
		),
		errCounter: promauto.NewCounterVec(
			prom.CounterOpts{
				Name: "koala_rpc_call_err_total",
				Help: "rpc调用错误数",
			},
			[]string{"service", "method", "grpc_code"},
		),
		costSummary: promauto.NewSummaryVec(
			prom.SummaryOpts{
				Name:       "koala_rpc_call_cost",
				Help:       "rpc调用耗时",
				Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
			},
			[]string{"service", "method"},
		),
	}
}

// IncRequest 增加请求数量
func (m *Metrics) IncRequest(ctx context.Context, service, method string) {
	m.requestCounter.WithLabelValues(service, method).Inc()
}

// IncRequestErr 增加请求错误数
func (m *Metrics) IncRequestErr(ctx context.Context, service, name string, err error) {
	st, _ := status.FromError(err)
	m.errCounter.WithLabelValues(service, name, st.Code().String()).Inc()
}

// Cost 记录耗时时间
func (m *Metrics) Cost(ctx context.Context, service, name string, us int64) {
	m.costSummary.WithLabelValues(service, name).Observe(float64(us))
}
