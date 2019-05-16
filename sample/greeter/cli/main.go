package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/server"
	wo "github.com/micro/go-plugins/wrapper/trace/opentracing"
	"github.com/opentracing/opentracing-go"
	"github.com/smartwalle/pks"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-lib/metrics"
	"time"

	"github.com/micro/go-plugins/registry/etcdv3"
	pks_client "github.com/smartwalle/pks/plugins/client/pks_grpc"
	pks_server "github.com/smartwalle/pks/plugins/server/pks_grpc"
	hello "github.com/smartwalle/pks/sample/greeter/srv/proto/hello"

	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
)

func main() {
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  "192.168.1.99:5775",
		},
	}

	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory

	closer, err := cfg.InitGlobalTracer("hello-cli", jaegercfg.Logger(jLogger), jaegercfg.Metrics(jMetricsFactory))

	if err != nil {
		fmt.Println(err)
		return
	}
	defer closer.Close()

	var s = pks.New(
		micro.Server(pks_server.NewServer(server.Address("192.168.1.99:8911"))),
		micro.Client(pks_client.NewClient(client.PoolSize(10), client.Wrap(wo.NewClientWrapper()))),
		micro.RegisterTTL(time.Second*10),
		micro.RegisterInterval(time.Second*5),
		micro.Registry(etcdv3.NewRegistry()),
		micro.Name("hello-cli"),
		micro.WrapHandler(wo.NewHandlerWrapper()),
	)

	cl := hello.NewSayService("hello-srv", s.Client())

	time.AfterFunc(time.Second*3, func() {
		span, ctx := opentracing.StartSpanFromContext(context.Background(), "cli-begin")
		span.LogKV("cli-key", "cli-value")
		span.Finish()

		rsp, _ := cl.Hello(ctx, &hello.Request{
			Name: "John",
		})

		fmt.Println(rsp.Msg)
	})

	span, ctx := opentracing.StartSpanFromContext(context.Background(), "cli-main")
	span.LogKV("cli-main-key", "cli-main-value")
	span.Finish()

	a1(ctx)

	s.Run()
}

func a1(ctx context.Context) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "a1")
	span.LogKV("cli-a1-key", "cli-a1-value")
	span.Finish()

	a2(ctx)
}

func a2(ctx context.Context) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "a2")
	span.LogKV("cli-a2-key", "cli-a2-value")
	a3(ctx)
	span.Finish()

}

func a3(ctx context.Context) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "a3")
	span.LogKV("cli-a3-key", "cli-a3-value")
	span.Finish()
}
