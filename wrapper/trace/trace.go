package trace

import (
	"context"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-micro/registry"
	"github.com/smartwalle/xid"
)

const (
	kTraceId = "trace-id"
)

func Wrap() client.Option {
	return client.WrapCall(wrap)
}

func wrap(cf client.CallFunc) client.CallFunc {
	return func(ctx context.Context, node *registry.Node, req client.Request, rsp interface{}, opts client.CallOptions) error {
		md, ok := metadata.FromContext(ctx)
		if ok == false {
			md = metadata.Metadata{}
		}

		_, ok = md[kTraceId]
		if ok == false {
			md[kTraceId] = xid.NewXID().Hex()
		}

		return cf(ctx, node, req, rsp, opts)
	}
}

func Id(ctx context.Context) string {
	md, ok := metadata.FromContext(ctx)
	if ok == false {
		return ""
	}

	return md[kTraceId]
}
