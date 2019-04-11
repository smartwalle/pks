package pks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/metadata"
	"github.com/smartwalle/pks/pb"
	"github.com/smartwalle/xid"
	"path/filepath"
	"sync"
	"time"
)

const (
	kHeaderFromAddress = "X-From-Address"
	kHeaderFromService = "X-From-Service"
	kHeaderFromId      = "X-From-Id"
	kHeaderToService   = "X-To-Service"
	kHeaderToPath      = "X-To-Path"
	kHeaderDate        = "X-Date"
	kHeaderTraceId     = "X-Trace-Id"

	kTimeFormat = "Mon, 02 Jan 2006 15:04:05 GMT"
)

var (
	PathNotFoundErr = errors.New("request path not found")
)

type HandlerFunc func(req *Request, rsp *Response) error

// --------------------------------------------------------------------------------
type Service struct {
	ms           micro.Service
	mu           sync.RWMutex
	h            map[string]HandlerFunc
	acceptStream chan *Stream
}

// --------------------------------------------------------------------------------
func New(opts ...micro.Option) *Service {
	var s = &Service{}

	s.ms = micro.NewService(opts...)
	s.acceptStream = make(chan *Stream)

	return s
}

//func (this *Service) Option(opts ...micro.Option) {
//	if this.ms != nil {
//		this.ms.Init(opts...)
//	}
//}

func (this *Service) Options() micro.Options {
	return this.ms.Options()
}

func (this *Service) Service() micro.Service {
	return this.ms
}

//func (this *Service) Server() server.Server {
//	if this.ms == nil {
//		return nil
//	}
//	return this.ms.Server()
//}
//
//func (this *Service) Client() client.Client {
//	if this.ms == nil {
//		return nil
//	}
//	return this.ms.Client()
//}

func (this *Service) ServerAddress() string {
	if this.ms == nil {
		return ""
	}
	if this.ms.Server() == nil {
		return ""
	}
	return this.ms.Server().Options().Address
}

func (this *Service) ServerName() string {
	if this.ms == nil {
		return ""
	}
	if this.ms.Server() == nil {
		return ""
	}
	return this.ms.Server().Options().Name
}

func (this *Service) ServerId() string {
	if this.ms == nil {
		return ""
	}
	if this.ms.Server() == nil {
		return ""
	}
	return this.ms.Server().Options().Id
}

func (this *Service) ServerVersion() string {
	if this.ms == nil {
		return ""
	}
	if this.ms.Server() == nil {
		return ""
	}
	return this.ms.Server().Options().Version
}

func (this *Service) Run() error {
	if this.ms == nil {
		return nil
	}

	pb.RegisterRPCHandler(this.ms.Server(), this)

	return this.ms.Run()
}

// --------------------------------------------------------------------------------
func (this *Service) Handle(path string, h HandlerFunc) {
	this.mu.Lock()
	defer this.mu.Unlock()

	if this.h == nil {
		this.h = make(map[string]HandlerFunc)
	}
	this.h[filepath.Join(path)] = h
}

// --------------------------------------------------------------------------------
func (this *Service) SimpleRequest(ctx context.Context, in *pb.Param, out *pb.Param) error {
	// 处理请求参数信息
	var req = &Request{}
	req.s = this
	req.ctx = ctx
	req.Body = in.Body

	// 从 ctx 中取出 metadata，并将 metadata 转换为 header
	meta, _ := metadata.FromContext(ctx)
	req.Header = WithMetadata(meta)
	req.localAddress = this.ServerAddress()

	// 处理响应参数信息
	var rsp = &Response{}

	if this.h != nil {
		this.mu.RLock()
		var h = this.h[req.Path()]
		if h == nil {
			h = this.h[this.ServerName()]
		}
		this.mu.RUnlock()

		if h == nil {
			return PathNotFoundErr
		}

		if err := h(req, rsp); err != nil {
			return err
		}
	}

	out.Body = rsp.Body

	// 处理响应头信息
	var header = rsp.Header
	if header == nil {
		header = Header{}
	}

	// 添加默认响应头信息
	header.Add(kHeaderFromAddress, this.ServerAddress())
	header.Add(kHeaderFromService, this.ServerName())
	header.Add(kHeaderFromId, this.ServerId())
	header.Add(kHeaderDate, time.Now().Format(kTimeFormat))
	header.Add(kHeaderToPath, req.Path())
	header.Add(kHeaderTraceId, req.TraceId())
	out.Header = header

	return nil
}

func (this *Service) StreamRequest(ctx context.Context, stream pb.RPC_StreamRequestStream) error {
	var nStream = newStream(this, stream)
	nStream.ctx = ctx

	// 从 ctx 中取出 metadata，并将 metadata 转换为 header
	meta, _ := metadata.FromContext(ctx)
	nStream.header = WithMetadata(meta)

	this.acceptStream <- nStream

	var err = nStream.waitDone()
	return err
}

func (this *Service) AcceptStream() (*Stream, error) {
	s := <-this.acceptStream
	return s, nil
}

// --------------------------------------------------------------------------------
func (this *Service) Request(ctx context.Context, service, path string, header Header, data interface{}, opts ...client.CallOption) (rsp *Response, err error) {
	ctx = this.ctxWrapper(ctx, service, path, header)

	var reqData []byte
	switch bt := data.(type) {
	case []byte:
		reqData = bt
	default:
		if reqData, err = json.Marshal(data); err != nil {
			return nil, err
		}
	}

	// 处理请求参数信息
	var req = &pb.Param{}
	req.Body = reqData

	// 发起请求
	var ts = pb.NewRPCService(service, this.Service().Client())
	sRsp, err := ts.SimpleRequest(ctx, req, opts...)
	if err != nil {
		return nil, err
	}

	// 处理返回参数信息
	rsp = &Response{}
	rsp.Body = sRsp.Body

	// 转换响应头信息
	rsp.Header = sRsp.Header
	rsp.localAddress = this.ServerAddress()

	return rsp, err
}

func (this *Service) RequestAddress(ctx context.Context, address, path string, header Header, body []byte, opts ...client.CallOption) (rsp *Response, err error) {
	var nOpts = make([]client.CallOption, 0, len(opts)+1)
	nOpts = append(nOpts, client.WithAddress(address))
	for _, opt := range opts {
		if opt != nil {
			nOpts = append(nOpts, opt)
		}
	}
	return this.Request(ctx, "", path, header, body, nOpts...)
}

// --------------------------------------------------------------------------------
func (this *Service) RequestStream(ctx context.Context, service, path string, header Header, opts ...client.CallOption) (*Stream, error) {
	ctx = this.ctxWrapper(ctx, service, path, header)

	// 发起请求
	var ts = pb.NewRPCService(service, this.Service().Client())

	var stream, err = ts.StreamRequest(ctx, opts...)
	if err != nil {
		return nil, err
	}

	var nStream = newStream(this, stream)
	nStream.ctx = ctx

	// 从 ctx 中取出 metadata，并将 metadata 转换为 header，此处记录的是发起流请求时的 header 信息
	meta, _ := metadata.FromContext(ctx)
	nStream.header = WithMetadata(meta)

	return nStream, err
}

func (this *Service) RequestStreamWithAddress(ctx context.Context, address, path string, header Header, opts ...client.CallOption) (*Stream, error) {
	var nOpts = make([]client.CallOption, 0, len(opts)+1)
	nOpts = append(nOpts, client.WithAddress(address))
	for _, opt := range opts {
		if opt != nil {
			nOpts = append(nOpts, opt)
		}
	}
	return this.RequestStream(ctx, "", path, header, nOpts...)
}

func (this *Service) ctxWrapper(ctx context.Context, service, path string, header Header) context.Context {
	if header == nil {
		header = Header{}
	}

	meta, _ := metadata.FromContext(ctx)
	for key, value := range meta {
		if header.Exists(key) == false {
			header.Add(key, value)
		}
	}

	// 添加默认值
	header.Add(kHeaderFromAddress, this.ServerAddress())
	header.Add(kHeaderFromService, this.ServerName())
	header.Add(kHeaderFromId, this.ServerId())
	header.Add(kHeaderDate, time.Now().Format(kTimeFormat))
	if len(path) > 0 {
		header.Add(kHeaderToPath, path)
	}
	if len(service) > 0 {
		header.Add(kHeaderToService, service)
	}

	// 添加 trace id
	if header.Exists(kHeaderTraceId) == false {
		header.Add(kHeaderTraceId, fmt.Sprintf("%s-%s", this.ServerName(), xid.NewXID().Hex()))
	}

	// 以 meta 为数据构建新的 ctx
	return metadata.NewContext(ctx, metadata.Metadata(header))
}
