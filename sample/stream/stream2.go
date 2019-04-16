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
		micro.RegisterTTL(time.Second*5),
		micro.RegisterInterval(time.Second*5),
		micro.Registry(etcdv3.NewRegistry()),
		micro.Name("st2"),
	)

	var h = pks.Header{}
	h.Add("ST2-Id", "ST2")
	var stream, err = s.RequestStream(context.Background(), "st1", "p", h)
	if err != nil {
		log4go.Errorln("请求建立流时发生错误:", err)
		return
	}

	log4go.Infoln("建立流成功, TraceId:", stream.TraceId())

	stream.Handle(func(s *pks.Stream, req *pks.Request, err error) error {
		return nil
	})

	h = pks.Header{}
	h.Add("PKG-Id", "ST2_PKG1")
	stream.Write(h, []byte("hhhhh"))

	select {}
}
