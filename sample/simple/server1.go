package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-plugins/registry/etcdv3"
	"github.com/smartwalle/pks"
	pks_client "github.com/smartwalle/pks/plugins/client/grpc"
	pks_server "github.com/smartwalle/pks/plugins/server/grpc"
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
	s.Handle(func(ctx context.Context, req *pks.Request, rsp *pks.Response) error {
		fmt.Printf("-----收到来自 %s 的请求-----\n", req.FromService())
		fmt.Printf("IP: %s \n", req.FromAddress())
		return nil
	})

	s.Run()
}
