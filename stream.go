package pks

import (
	"context"
	"github.com/smartwalle/pks/pb"
	"sync"
)

type StreamHandler interface {
	OnMessage(*Stream, *Request)

	OnClose(*Stream, error)
}

type Stream struct {
	s         *Service
	stream    pb.RPC_StreamRequestStream
	h         StreamHandler
	header    Header
	ctx       context.Context
	done      chan error
	mu        sync.RWMutex
	once      sync.Once
	closeOnce sync.Once
	data      map[string]interface{}
}

func newStream(s *Service, stream pb.RPC_StreamRequestStream) *Stream {
	var ns = &Stream{}
	ns.s = s
	ns.stream = stream
	ns.done = make(chan error)
	return ns
}

func (this *Stream) Context() context.Context {
	return this.ctx
}

func (this *Stream) Header() Header {
	return this.header
}

func (this *Stream) FromService() string {
	return this.header.Get(kHeaderFromService)
}

func (this *Stream) FromAddress() string {
	return this.header.Get(kHeaderFromAddress)
}

func (this *Stream) waitDone() error {
	return <-this.done
}

func (this *Stream) read() {
	var err error
	defer func() {
		this.close(err)
	}()

	var param *pb.Param
	for {
		param, err = this.stream.Recv()
		if err != nil {
			return
		}

		if this.h != nil {
			var req = &Request{}
			req.s = this.s

			if param != nil {
				req.Body = param.Body

				// 转换请求头信息
				// 将建立流的请求头合并到流消息中
				//req.Header = this.header
				//for k, v := range param.Header {
				//	req.Header.Add(k, v)
				//}

				// 不合并建立流时的请求头
				req.Header = param.Header
			}

			req.localAddress = this.s.ServerAddress()

			if this.h != nil {
				this.h.OnMessage(this, req)
			}
		}
	}
}

func (this *Stream) Write(h Header, data []byte) error {
	var header = h
	if header == nil {
		header = Header{}
	}
	// 添加默认请求头信息
	header.Add(kHeaderFromAddress, this.s.ServerAddress())
	header.Add(kHeaderFromService, this.s.ServerName())
	header.Add(kHeaderFromId, this.s.ServerId())

	var out = &pb.Param{}
	out.Body = data
	out.Header = header
	return this.stream.Send(out)
}

func (this *Stream) Close() error {
	return this.close(nil)
}

func (this *Stream) close(err error) error {
	this.closeOnce.Do(func() {
		select {
		case this.done <- err:
		default:
		}

		close(this.done)
		this.done = nil

		if this.h != nil {
			this.h.OnClose(this, err)
		}
	})
	return this.stream.Close()
}

func (this *Stream) Handle(h StreamHandler) {
	this.h = h

	this.once.Do(func() {
		go this.read()
	})
}

func (this *Stream) Set(key string, value interface{}) {
	this.mu.Lock()
	defer this.mu.Unlock()
	if value != nil {
		if this.data == nil {
			this.data = make(map[string]interface{})
		}
		this.data[key] = value
	}
}

func (this *Stream) Get(key string) interface{} {
	this.mu.RLock()
	defer this.mu.RUnlock()
	if this.data == nil {
		return nil
	}
	return this.data[key]
}

func (this *Stream) Del(key string) {
	this.mu.Lock()
	defer this.mu.Unlock()
	if this.data == nil {
		return
	}
	delete(this.data, key)
}
