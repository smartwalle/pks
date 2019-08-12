package grpc

import (
	"google.golang.org/grpc"
	"sync"
	"time"
)

type streamPool struct {
	sync.Mutex
	size   int
	conns  map[string][]*poolConn
	indexs map[string]int
}

func newStreamPool(size int) *streamPool {
	return &streamPool{
		size:   size,
		conns:  make(map[string][]*poolConn),
		indexs: make(map[string]int),
	}
}

func (p *streamPool) getConn(addr string, opts ...grpc.DialOption) (*poolConn, error) {
	p.Lock()
	defer p.Unlock()

	conns := p.conns[addr]
	index := p.indexs[addr]

	if len(conns) == 0 {
		for i := len(conns); i < p.size; i++ {
			cc, err := grpc.Dial(addr, opts...)
			if err != nil {
				return nil, err
			}
			pc := &poolConn{cc, time.Now().Unix()}
			conns = append(conns, pc)
		}

		p.conns[addr] = conns
	}
	conn := conns[index]

	index++
	if index >= p.size {
		index = 0
	}
	p.indexs[addr] = index
	return conn, nil
}
