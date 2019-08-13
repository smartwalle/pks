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

func (this *base) FromAddress() string {
	return this.Header.Get(kHeaderFromAddress)
}

func (this *base) LocalAddress() string {
	return this.localAddress
}

// --------------------------------------------------------------------------------
type Request struct {
	base
	s *Service
}

func (this *Request) Request(ctx context.Context, path string, header Header, data []byte, opts ...client.CallOption) (rsp *Response, err error) {
	var nOpts = make([]client.CallOption, 0, len(opts)+1)
	nOpts = append(nOpts, client.WithAddress(this.FromAddress()))

	for _, opt := range opts {
		if opt != nil {
			nOpts = append(nOpts, opt)
		}
	}

	return this.s.Request(ctx, this.FromService(), header, data, nOpts...)
}

// --------------------------------------------------------------------------------
type Response struct {
	base
}
