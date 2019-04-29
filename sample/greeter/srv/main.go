package main

import (
	"fmt"
	"github.com/micro/go-micro/server"
	"github.com/smartwalle/tx4go"
	"log"

	hello "github.com/micro/examples/greeter/srv/proto/hello"

	"context"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-plugins/registry/etcdv3"
	"github.com/smartwalle/pks"
	pks_client "github.com/smartwalle/pks/plugins/client/pks_grpc"
	pks_server "github.com/smartwalle/pks/plugins/server/pks_grpc"
	"time"
)

type Say struct{}

func (s *Say) Hello(ctx context.Context, req *hello.Request, rsp *hello.Response) error {
	log.Print("Received Say.Hello request")
	rsp.Msg = "Hello " + req.Name

	tx, _ := tx4go.Begin(ctx, func() {
		fmt.Println("confirm")
	}, func() {
		fmt.Println("cancel")
	})

	tx.Commit()

	return nil
}

func main() {
	var s = pks.New(
		micro.Server(pks_server.NewServer(server.Address("192.168.1.99:8921"))),
		micro.Client(pks_client.NewClient(client.PoolSize(10))),
		micro.RegisterTTL(time.Second*10),
		micro.RegisterInterval(time.Second*5),
		micro.Registry(etcdv3.NewRegistry()),
		micro.Name("hello-srv"),
	)

	tx4go.Init(s)

	hello.RegisterSayHandler(s.Server(), &Say{})

	s.Run()
}
