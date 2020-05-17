package server

import (
	"context"
	"fmt"
	"koala/logger"
	"koala/middleware"
	"koala/registry"
	_ "koala/registry/etcd" // 为了注册etcd插件
	"koala/util"
	"net"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
)

var koalaServer = &KoalaServer{
	Server: grpc.NewServer(),
}

// KoalaServer 封装服务启动与一些额外处理
// 字段
//   Server: *grpc.Server, rpc服务器
//   limiter: 限流器
//   userMiddlewares: 用户的中间件
type KoalaServer struct {
	*grpc.Server
	limiter         *rate.Limiter
	userMiddlewares []middleware.Middleware
}

// InitServer 初始化服务器
// 参数
//   service: 服务名
// 返回值
//   err: error
func InitServer(service string) (err error) {
	err = initConfig(service)
	if err != nil {
		return
	}

	// 初始化日志
	logger.Start(logger.WithType(koalaConf.Log.Type),
		logger.WithLevel(koalaConf.Log.Level),
		logger.WithModuleName(service),
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

	// 初始化注册中心
	err = initRegister(service)
	if err != nil {
		logger.Error(context.TODO(), "初始化注册中心失败[%#v], err=%v", koalaConf.Register, err)
		return
	}

	ctx := context.TODO()
	logger.Info(ctx, "初始化注册中心成功")

	// 初始化限流器
	if koalaConf.Limit.SwitchOn {
		koalaServer.limiter = rate.NewLimiter(rate.Limit(koalaConf.Limit.QPS), koalaConf.Limit.QPS)
	}

	// 初始化追踪系统
	err = initTrace(service)
	if err != nil {
		logger.Warn(ctx, "初始化分布式追踪[%#v]失败, err=%v", koalaConf.Trace, err)
		err = nil // 初始化追踪系统失败不会影响业务处理
	}

	return
}

// GetServer 获取grpc server
func GetServer() *grpc.Server {
	return koalaServer.Server
}

// Run 启动服务
func Run() (err error) {
	defer closeServer()

	ctx := context.TODO()
	if koalaConf.Prometheus.SwitchOn {
		go func() {
			addr := fmt.Sprintf(":%d", koalaConf.Prometheus.Port)
			http.Handle("/metrics", promhttp.Handler())
			runErr := http.ListenAndServe(addr, nil)
			if runErr != nil {
				logger.Warn(ctx, "Prometheus采样打点失败, err=%v", runErr)
			}
		}()
	}

	addr := fmt.Sprintf(":%d", koalaConf.Service.Port)
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Error(ctx, "监听端口失败, err=%v", err)
		return
	}
	err = koalaServer.Serve(listen)
	if err != nil {
		logger.Error(ctx, "启动server失败, err=%v", err)
		return
	}

	return
}

// Use 用户添加中间件
// 参数
//   m: Middleware不定参
func Use(m ...middleware.Middleware) {
	koalaServer.userMiddlewares = append(koalaServer.userMiddlewares, m...)
}

// BuildServerMiddleware 构建中间件
func BuildServerMiddleware(handle middleware.HandleFunc) middleware.HandleFunc {
	mids := []middleware.Middleware{}

	mids = append(mids, middleware.AccessLogMiddleware)

	if koalaConf.Prometheus.SwitchOn {
		mids = append(mids, middleware.PrometheusServerMiddleware)
	}

	if koalaConf.Limit.SwitchOn {
		mids = append(mids, middleware.NewLimiterMiddleware(koalaServer.limiter))
	}

	if koalaConf.Trace.SwitchOn {
		mids = append(mids, middleware.TraceServerMiddleware)
	}

	if len(koalaServer.userMiddlewares) != 0 {
		mids = append(mids, koalaServer.userMiddlewares...)
	}
	if len(mids) > 0 {
		m := middleware.Chain(middleware.PrepareMiddleware, mids...)
		return m(handle)
	}

	// 没有则返回自己
	return handle
}

func closeServer() {
	middleware.CloseTrace()
	logger.Stop()
}

func initRegister(service string) (err error) {
	if !koalaConf.Register.SwitchOn {
		return
	}

	ctx, cancle := context.WithTimeout(context.TODO(), time.Second*60)
	defer cancle()

	regInst, err := registry.InitPlugin(ctx,
		koalaConf.Register.RegisterName,
		registry.WithAddrs(koalaConf.Register.RegisterAddrs),
		registry.WithRegistryPath(koalaConf.Register.RegisterPath),
		registry.WithTimeout(time.Duration(koalaConf.Register.Timeout)*time.Second),
		registry.WithHeartBeat(koalaConf.Register.Heartbeat),
	)
	if err != nil {
		return
	}

	ip, err := util.GetLocalIP()
	if err != nil {
		return
	}

	regService := &registry.Service{
		Name: service,
	}
	regService.Nodes = append(regService.Nodes, &registry.Node{
		IP:   ip,
		Port: koalaConf.Service.Port,
	})
	err = regInst.Register(ctx, regService)

	return
}

func initTrace(service string) (err error) {
	if !koalaConf.Trace.SwitchOn {
		return
	}
	err = middleware.InitTrace(middleware.WithAppName(service),
		middleware.WithReportAddr(koalaConf.Trace.ReportAddr),
		middleware.WithSampleType(koalaConf.Trace.SampleType),
		middleware.WithSampleRate(koalaConf.Trace.SampleRate),
	)
	return
}
