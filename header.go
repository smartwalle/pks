package pks

import (
	"github.com/micro/go-micro/metadata"
	"strings"
)

type Header map[string]string

func WithMetadata(m metadata.Metadata) Header {
	var h = Header{}
	for key, value := range m {
		h.Set(key, value)
	}
	return h
}

func (h Header) Add(key, value string) {
	h[strings.ToLower(key)] = value
}

func (h Header) Set(key, value string) {
	h[strings.ToLower(key)] = value
}

func (h Header) Get(key string) string {
	if h == nil {
		return ""
	}
	return h[strings.ToLower(key)]
}

func (h Header) Del(key string) {
	delete(h, strings.ToLower(key))
}

func (h Header) Exists(key string) bool {
	_, ok := h[strings.ToLower(key)]
	return ok
}
