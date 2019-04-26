package main

import (
	"context"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-plugins/registry/etcdv3"
	"github.com/smartwalle/log4go"
	"github.com/smartwalle/pks"
	pks_client "github.com/smartwalle/pks/plugins/client/pks_grpc"
	pks_server "github.com/smartwalle/pks/plugins/server/pks_grpc"
	"time"
)

func main() {
	var s = pks.New(
		micro.Server(pks_server.NewServer()),
		micro.Client(pks_client.NewClient(client.PoolSize(10))),
		micro.RegisterTTL(time.Second*10),
		micro.RegisterInterval(time.Second*5),
		micro.Registry(etcdv3.NewRegistry()),
		micro.Name("s1"),
	)

	// 默认
	s.Handle("", func(ctx context.Context, req *pks.Request, rsp *pks.Response) error {
		log4go.Infof("-----收到来自 %s 的请求-----\n", req.FromService())
		log4go.Infof("IP: %s, TraceId: %s \n", req.FromAddress(), req.TraceId())
		return nil
	})

	s.Handle("h1", func(ctx context.Context, req *pks.Request, rsp *pks.Response) error {
		log4go.Infof("-----收到来自 %s 的请求-----\n", req.FromService())
		log4go.Infof("IP: %s, TraceId: %s \n", req.FromAddress(), req.TraceId())
		log4go.Infoln("请求头")
		for key, value := range req.Header {
			log4go.Infoln(key, value)
		}
		return nil
	})

	s.Run()
}
