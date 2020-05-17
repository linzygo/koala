package client

import (
	"context"
	"fmt"
	"koala/loadbalance"
	"koala/logger"
	"koala/middleware"
	"koala/registry"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/time/rate"
)

var koalaClient = &KoalaClient{}

// KoalaClient 客户端
// 字段
//   discoveryInst: 服务发现
//   balancer: 负载均衡器
//   limiter: 限流器
type KoalaClient struct {
	discoveryInst registry.Registry
	balancer      loadbalance.LoadBalance
	limiter       *rate.Limiter
}

// InitClient 初始化整个client
func InitClient(clientName string) (err error) {
	err = initConfig()
	if err != nil {
		return
	}

	// 初始化日志
	logger.Start(logger.WithType(koalaConf.Log.Type),
		logger.WithLevel(koalaConf.Log.Level),
		logger.WithModuleName(clientName),
		logger.WithChanSize(koalaConf.Log.ChanSize),
		logger.WithLogDir(koalaConf.Log.Path),
		logger.WithFileSize(koalaConf.Log.FileSize),
	)
	defer func() {
		// 有错误发生, 需要Stop, 避免程序退出导致日志来不及打印
		if err != nil {
			logger.Stop()
		}
	}()

	ctx := context.TODO()

	logger.Info(ctx, "rpc客户端配置: %#v", koalaConf)

	// 初始化服务发现
	err = initDiscovery()
	if err != nil {
		logger.Error(ctx, "初始化服务发现失败[%#v], err=%v", koalaConf.Discovery, err)
		return
	}

	// 初始化负载均衡器
	koalaClient.balancer = loadbalance.NewLoadBalance(koalaConf.Balancer.Name)

	// 初始化Prometheus采样打点
	initPrometheus()

	// 初始化限流器
	initLimiter()

	// 初始化追踪系统
	initTrace(clientName)

	return
}

func initDiscovery() (err error) {
	if !koalaConf.Discovery.SwitchOn {
		return
	}

	ctx, cancle := context.WithTimeout(context.TODO(), time.Second*60)
	defer cancle()

	koalaClient.discoveryInst, err = registry.InitPlugin(ctx,
		koalaConf.Discovery.RegisterName,
		registry.WithAddrs(koalaConf.Discovery.RegisterAddrs),
		registry.WithRegistryPath(koalaConf.Discovery.RegisterPath),
		registry.WithTimeout(time.Duration(koalaConf.Discovery.Timeout)*time.Second),
		registry.WithHeartBeat(koalaConf.Discovery.Heartbeat),
	)
	if err != nil {
		return
	}
	return
}

func initPrometheus() (err error) {
	if koalaConf.Prometheus.SwitchOn {
		go func() {
			addr := fmt.Sprintf(":%d", koalaConf.Prometheus.Port)
			http.Handle("/metrics", promhttp.Handler())
			runErr := http.ListenAndServe(addr, nil)
			if runErr != nil {
				logger.Warn(context.TODO(), "Prometheus采样打点失败, err=%v", runErr)
			}
		}()
	}

	return
}

func initLimiter() (err error) {
	if koalaConf.Limit.SwitchOn {
		koalaClient.limiter = rate.NewLimiter(rate.Limit(koalaConf.Limit.QPS), koalaConf.Limit.QPS)
	}
	return
}

func initTrace(clientName string) (err error) {
	if koalaConf.Trace.SwitchOn {
		err = middleware.InitTrace(middleware.WithAppName(clientName),
			middleware.WithReportAddr(koalaConf.Trace.ReportAddr),
			middleware.WithSampleType(koalaConf.Trace.SampleType),
			middleware.WithSampleRate(koalaConf.Trace.SampleRate),
		)
	}
	return
}
