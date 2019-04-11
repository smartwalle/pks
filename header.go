package pks

import (
	"github.com/micro/go-micro/metadata"
	"net/http"
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
	h[http.CanonicalHeaderKey(key)] = value
}

func (h Header) Set(key, value string) {
	h[http.CanonicalHeaderKey(key)] = value
}

func (h Header) Get(key string) string {
	if h == nil {
		return ""
	}
	return h[http.CanonicalHeaderKey(key)]
}

func (h Header) Del(key string) {
	delete(h, http.CanonicalHeaderKey(key))
}

func (h Header) Exists(key string) bool {
	_, ok := h[http.CanonicalHeaderKey(key)]
	return ok
}
