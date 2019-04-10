package pks_grpc

import (
	"bytes"
)

type buffer struct {
	*bytes.Buffer
}

func (b *buffer) Close() error {
	b.Buffer.Reset()
	return nil
}
