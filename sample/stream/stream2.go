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
	var s = pks.NewService(
		micro.Server(pks_server.NewServer()),
		micro.Client(pks_client.NewClientWithSteamPoolSize(15, client.PoolSize(10))),
		micro.RegisterTTL(time.Minute*5),
		micro.RegisterInterval(time.Second*5),
		micro.Registry(etcdv3.NewRegistry()),
		micro.Name("st2"),
	)

	for i := 0; i < 10000; i++ {
		time.Sleep(time.Second * 1)

		var h = pks.Header{}
		h.Add("ST2-Id", "ST2")
		var stream, err = s.RequestStream(context.Background(), "st1", h)
		if err != nil {
			fmt.Println("请求建立流时发生错误:", err)
			continue
		}

		fmt.Println("建立流成功")

		//stream.Handle(func(s *pks.Stream, req *pks.Request, err error) error {
		//	if err != nil {
		//		fmt.Println("err ", err)
		//		return err
		//	}
		//	fmt.Println("eee", err)
		//	return nil
		//})
		//
		//go func() {
		//	time.Sleep(time.Second * 5)
		//	fmt.Println(stream.Write(nil, []byte("ee")))
		//}()

		stream.Handle(&Stream2Handler{})

		//h = pks.Header{}
		//h.Add("PKG-Id", "ST2_PKG1")
		//
		//for ii:=0; ii<1000; ii++ {
		//	fmt.Println(stream.Write(h, []byte("hhhhh")))
		//	time.Sleep(time.Second * 2)
		//}

		stream.Write(nil, []byte("x"))
	}

	select {}
}

type Stream2Handler struct {
}

func (this *Stream2Handler) OnMessage(s *pks.Stream, req *pks.Request) {
	fmt.Println("message", string(req.Body))
}

func (this *Stream2Handler) OnClose(s *pks.Stream, err error) {
	fmt.Println("close", err)
}
