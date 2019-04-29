package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/server"
	"github.com/smartwalle/pks"
	"time"

	hello "github.com/micro/examples/grpc/greeter/srv/proto/hello"
	"github.com/micro/go-plugins/registry/etcdv3"
	pks_client "github.com/smartwalle/pks/plugins/client/pks_grpc"
	pks_server "github.com/smartwalle/pks/plugins/server/pks_grpc"
)

func main() {
	var s = pks.New(
		micro.Server(pks_server.NewServer(server.Address("192.168.1.99:8911"))),
		micro.Client(pks_client.NewClient(client.PoolSize(10))),
		micro.RegisterTTL(time.Second*10),
		micro.RegisterInterval(time.Second*5),
		micro.Registry(etcdv3.NewRegistry()),
		micro.Name("hello-cli"),
	)

	cl := hello.NewSayService("hello-srv", s.Client())

	time.AfterFunc(time.Second*3, func() {
		rsp, _ := cl.Hello(context.Background(), &hello.Request{
			Name: "John",
		})

		fmt.Println(rsp.Msg)
	})

	s.Run()
}
