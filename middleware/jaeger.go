package middleware

import (
	"context"
	"io"
	"koala/logger"
	"koala/meta"
	"koala/util"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/transport/zipkin"
	"google.golang.org/grpc/metadata"
)

var closeTrace io.Closer

// InitTrace 初始化分布式追踪, 不再用Trace时需要调用CloseTrace
func InitTrace(opts ...TraceOption) (err error) {
	cfg := &TraceConfig{
		AppName:    "UnknowService",
		ReportAddr: "http://localhost:9411/api/v1/spans",
		SampleType: "const",
		SampleRate: 1,
	}
	for _, opt := range opts {
		opt(cfg)
	}

	transport, err := zipkin.NewHTTPTransport(
		cfg.ReportAddr,
		zipkin.HTTPBatchSize(16),
		zipkin.HTTPLogger(jaeger.StdLogger),
	)
	if err != nil {
		return
	}

	jaegerCfg := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  cfg.SampleType,
			Param: float64(cfg.SampleRate),
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
	}

	r := jaeger.NewRemoteReporter(transport)
	tracer, close, err := jaegerCfg.New(cfg.AppName,
		config.Logger(jaeger.StdLogger),
		config.Reporter(r))
	if err != nil {
		return
	}

	closeTrace = close

	opentracing.SetGlobalTracer(tracer)
	return
}

// CloseTrace 关闭追踪
func CloseTrace() {
	if closeTrace != nil {
		closeTrace.Close()
		closeTrace = nil
	}
}

// TraceServerMiddleware Trace中间件
func TraceServerMiddleware(handle HandleFunc) HandleFunc {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		//从ctx获取grpc的metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			//没有的话,新建一个
			md = metadata.Pairs()
		}

		tracer := opentracing.GlobalTracer()
		parentSpanContext, err := tracer.Extract(opentracing.HTTPHeaders, metadataTextMap(md))
		if err != nil && err != opentracing.ErrSpanContextNotFound {
			logger.Warn(ctx, "trace extract failed, parsing trace information: %v", err)
		}

		serverMeta := meta.GetServerMeta(ctx)
		//开始追踪该方法
		serverSpan := tracer.StartSpan(
			serverMeta.Method,
			ext.RPCServerOption(parentSpanContext),
			ext.SpanKindRPCServer,
		)
		defer serverSpan.Finish()

		serverSpan.SetTag(util.TraceID, logger.GetTraceID(ctx))
		ctx = opentracing.ContextWithSpan(ctx, serverSpan)
		resp, err = handle(ctx, req)
		//记录错误
		if err != nil {
			ext.Error.Set(serverSpan, true)
			serverSpan.LogFields(log.String("event", "error"), log.String("message", err.Error()))
		}

		return
	}
}

// TraceRPCMiddleware Trace中间件, 用于追踪rpc的调用情况
func TraceRPCMiddleware(handle HandleFunc) HandleFunc {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		// 获取parent span
		var parentCtx opentracing.SpanContext
		if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
			parentCtx = parentSpan.Context()
		}

		rpcMeta := meta.GetClientRPCMeta(ctx)

		tracer := opentracing.GlobalTracer()
		clientSpan := tracer.StartSpan(rpcMeta.ServiceName,
			opentracing.ChildOf(parentCtx),
			ext.SpanKindRPCClient,
		)
		defer clientSpan.Finish()

		traceID := logger.GetTraceID(ctx)

		clientSpan.SetTag(string(ext.Component), "koala-rpc")
		clientSpan.SetTag(util.TraceID, traceID)

		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.Pairs()
		}
		if err := tracer.Inject(clientSpan.Context(), opentracing.HTTPHeaders, metadataTextMap(md)); err != nil {
			logger.Warn(ctx, "grpc_opentracing: failed serializing trace information: %v", err)
		}

		// 把md、traceID放到Context, grpc使用该Context把Context内容传给服务
		ctx = metadata.NewOutgoingContext(ctx, md)
		ctx = metadata.AppendToOutgoingContext(ctx, util.TraceID, traceID)

		ctx = opentracing.ContextWithSpan(ctx, clientSpan)
		resp, err = handle(ctx, req)
		//记录错误
		if err != nil {
			ext.Error.Set(clientSpan, true)
			clientSpan.LogFields(log.String("event", "error"), log.String("message", err.Error()))
		}

		return
	}
}
