package pks

import (
	"context"
	"encoding/json"
	"github.com/micro/go-micro/client"
	"strings"
)

const (
	kMicroFromService = "micro-from-service"
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

func (this *base) FromServices() []string {
	var v = this.Header.Get(kMicroFromService)
	if v != "" {
		return strings.Split(v, ",")
	}
	return nil
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

func (this *base) Unmarshal(v interface{}) error {
	return json.Unmarshal(this.Body, v)
}

// --------------------------------------------------------------------------------
type Request struct {
	base
	t   *Service
	ctx context.Context
}

func (this *Request) Context() context.Context {
	return this.ctx
}

func (this *Request) TraceId() string {
	return this.Header.Get(kHeaderTraceId)
}

func (this *Request) Request(ctx context.Context, path string, header Header, data interface{}, opts ...client.CallOption) (rsp *Response, err error) {
	if this.t != nil {
		var nOpts = make([]client.CallOption, 0, len(opts)+1)
		nOpts = append(nOpts, client.WithAddress(this.FromAddress()))

		for _, opt := range opts {
			if opt != nil {
				nOpts = append(nOpts, opt)
			}
		}

		return this.t.Request(ctx, this.FromService(), path, header, data, nOpts...)
	}
	return nil, PathNotFoundErr
}

// --------------------------------------------------------------------------------
type Response struct {
	base
}
