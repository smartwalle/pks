package main

import (
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
		micro.Client(pks_client.NewClient(client.PoolSize(10))),
		micro.RegisterTTL(time.Minute*5),
		micro.RegisterInterval(time.Second*5),
		micro.Registry(etcdv3.NewRegistry()),
		micro.Name("st1"),
	)

	var i = 0

	go func() {
		for {
			var stream, err = s.AcceptStream()
			if err != nil {
				fmt.Println("处理流请求时发生错误:", err)
				continue
			}
			i++

			fmt.Println("=======", i)

			fmt.Println("-----建立新的流-----")
			//fmt.Println("流请求头")
			//for key, value := range stream.Header() {
			//	fmt.Println(key, value)
			//}

			//stream.Handle(func(s *pks.Stream, req *pks.Request, err error) error {
			//	if err != nil {
			//		fmt.Println("接收流消息时发生错误:", err)
			//		return err
			//	}
			//
			//	//fmt.Println("-----收到新的流消息-----")
			//	//fmt.Println("流消息请求头")
			//	//for key, value := range req.Header {
			//	//	fmt.Println(key, value)
			//	//}
			//	//
			//	//fmt.Println("流消息内容:", string(req.Body))
			//	fmt.Println("ell")
			//	return nil
			//})
			//
			//stream.Write(nil, []byte("dde"))
			//time.Sleep(time.Second * 5)
			//stream.Close()

			stream.Handle(&Stream1Handler{})
		}
	}()

	s.Run()
}

type Stream1Handler struct {
}

func (this *Stream1Handler) OnMessage(s *pks.Stream, req *pks.Request) {
	fmt.Println("message", string(req.Body))

	s.Write(nil, []byte("e"))
}

func (this *Stream1Handler) OnClose(s *pks.Stream, err error) {
	fmt.Println("close", err)
}
