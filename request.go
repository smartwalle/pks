package pks

import (
	"context"
	"github.com/micro/go-micro/client"
)

// --------------------------------------------------------------------------------
type base struct {
	localAddress string
	Header       Header
	Body         []byte
}

func (this *base) FromService() string {
	return this.Header.Get(kHeaderFromService)
}

func (this *base) Path() string {
	return this.Header.Get(kHeaderToPath)
}

func (this *base) FromAddress() string {
	return this.Header.Get(kHeaderFromAddress)
}

func (this *base) LocalAddress() string {
	return this.localAddress
}

// --------------------------------------------------------------------------------
type Request struct {
	base
	s   *Service
	ctx context.Context
}

func (this *Request) Context() context.Context {
	return this.ctx
}

func (this *Request) TraceId() string {
	return this.Header.Get(kHeaderTraceId)
}

func (this *Request) Request(ctx context.Context, path string, header Header, data interface{}, opts ...client.CallOption) (rsp *Response, err error) {
	if this.s != nil {
		var nOpts = make([]client.CallOption, 0, len(opts)+1)
		nOpts = append(nOpts, client.WithAddress(this.FromAddress()))

		for _, opt := range opts {
			if opt != nil {
				nOpts = append(nOpts, opt)
			}
		}

		return this.s.Request(ctx, this.FromService(), path, header, data, nOpts...)
	}
	return nil, PathNotFoundErr
}

// --------------------------------------------------------------------------------
type Response struct {
	base
}
