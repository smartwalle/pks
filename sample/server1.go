package main

import (
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
		micro.Name("s1"),
	)

	s.Handle("p", func(req *pks.Request, rsp *pks.Response) error {
		fmt.Println("Handle Request", req.TraceId())
		fmt.Println("Handle Request", req.Header)

		var r, err = s.Request(req.Context(), "s2", "q", nil, nil)
		fmt.Println("Req", r, err)

		return nil
	})

	s.Run()
}
