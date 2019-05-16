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

type Say struct{}

func (s *Say) Hello(ctx context.Context, req *hello.Request, rsp *hello.Response) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "srv-hello-begin")
	span.LogKV("srv-key", "srv-value")
	span.Finish()

	fmt.Print("Received Say.Hello request")
	rsp.Msg = "Hello " + req.Name
	return nil
}

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

	closer, err := cfg.InitGlobalTracer("hello-srv", jaegercfg.Logger(jLogger), jaegercfg.Metrics(jMetricsFactory))

	if err != nil {
		fmt.Println(err)
		return
	}
	defer closer.Close()

	var s = pks.New(
		micro.Server(pks_server.NewServer(server.Address("192.168.1.99:8921"))),
		micro.Client(pks_client.NewClient(client.PoolSize(10))),
		micro.RegisterTTL(time.Second*10),
		micro.RegisterInterval(time.Second*5),
		micro.Registry(etcdv3.NewRegistry()),
		micro.Name("hello-srv"),
		micro.WrapHandler(wo.NewHandlerWrapper()),
	)

	hello.RegisterSayHandler(s.Server(), &Say{})

	s.Run()
}
