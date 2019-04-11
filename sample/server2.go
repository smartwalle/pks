package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-plugins/registry/etcdv3"
	"github.com/smartwalle/pks"
	pks_client "github.com/smartwalle/pks/plugins/client/pks_grpc"
	pks_server "github.com/smartwalle/pks/plugins/server/pks_grpc"
	"time"
)

func main() {
	var s = pks.New(
		micro.Server(pks_server.NewServer()),
		micro.Client(pks_client.NewClient(client.PoolSize(10))),
		micro.RegisterTTL(time.Second*5),
		micro.RegisterInterval(time.Second*5),
		micro.Registry(etcdv3.NewRegistry()),
		micro.Name("s2"),
	)

	s.Handle("q", func(req *pks.Request, rsp *pks.Response) error {
		fmt.Println(req.TraceId())
		fmt.Println(req.Header)
		return nil
	})

	time.AfterFunc(time.Second*2, func() {
		fmt.Println(s.Request(context.Background(), "s1", "p", nil, nil))
	})

	s.Run()
}
