package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/server"
	"github.com/smartwalle/pks"
	"github.com/smartwalle/tx4go"
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

	tx4go.Init(s)

	time.AfterFunc(time.Second*3, func() {
		tx, err := tx4go.Begin(context.Background(), func() {
			fmt.Println("confirm")
		}, func() {
			fmt.Println("cancel")
		})

		rsp, err := cl.Hello(tx.Context(), &hello.Request{
			Name: "John",
		})

		if err != nil {
			tx.Rollback()
			return
		}

		tx.Commit()

		fmt.Println(rsp.Msg)
	})

	s.Run()
}
