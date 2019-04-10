package pks

import "github.com/micro/go-micro/metadata"

type Header metadata.Metadata

func (h Header) Add(key, value string) {
	h[key] = value
}

func (h Header) Set(key, value string) {
	h[key] = value
}

func (h Header) Get(key string) string {
	if h == nil {
		return ""
	}
	return h[key]
}

func (h Header) Del(key string) {
	delete(h, key)
}
