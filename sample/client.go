package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-plugins/registry/etcdv3"
	"github.com/smartwalle/pks"
	pks_client "github.com/smartwalle/pks/plugins/client/pks_grpc"
)

func main() {
	var c = pks.New(
		micro.Client(pks_client.NewClient(client.PoolSize(10))),
		micro.Registry(etcdv3.NewRegistry()),
		micro.Name("c"),
	)

	fmt.Println(c.Request(context.Background(), "s", "p", nil, nil))
}
