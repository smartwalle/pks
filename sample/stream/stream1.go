package main

import (
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
		micro.Name("st1"),
	)

	go func() {
		for {
			var stream, err = s.AcceptStream()
			if err != nil {
				log4go.Errorln("处理流请求时发生错误:", err)
				continue
			}

			log4go.Infoln("-----建立新的流-----")
			log4go.Infof("Path: %s, TraceId: %s \n", stream.Path(), stream.TraceId())
			log4go.Infoln("流请求头")
			for key, value := range stream.Header() {
				log4go.Infoln(key, value)
			}

			stream.Handle(func(s *pks.Stream, req *pks.Request, err error) error {
				if err != nil {
					log4go.Errorln("接收流消息时发生错误:", err)
					return err
				}

				log4go.Infoln("-----收到新的流消息-----")
				log4go.Infoln("流消息请求头")
				for key, value := range req.Header {
					log4go.Infoln(key, value)
				}

				log4go.Infoln("流消息内容:", string(req.Body))
				return nil
			})
		}
	}()

	s.Run()
}
