package cache

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/metafates/gache"
	"github.com/awirix/awirix/filesystem"
	"github.com/awirix/awirix/where"
	"net/http"
	"net/http/httputil"
	"path/filepath"
	"time"
)

var ErrCacheTooLarge = errors.New("cache too large")

var HTTP *httpCache

func init() {
	cache := gache.New[map[string][]byte](&gache.Options{
		Path:       filepath.Join(where.Cache(), "http.json"),
		Lifetime:   time.Hour * 5,
		FileSystem: &filesystem.GacheFs{},
	})

	HTTP = &httpCache{
		cache: cache,
	}
}

type httpCache struct {
	cache *gache.Cache[map[string][]byte]
}

func (h *httpCache) Get(r *http.Request) (*http.Response, bool) {
	data, expired, err := h.cache.Get()
	if err != nil {
		return nil, false
	}

	if expired {
		return nil, false
	}

	dumpedRequest, err := httputil.DumpRequest(r, true)
	if err != nil {
		return nil, false
	}

	dumpedResponse, ok := data[string(dumpedRequest)]
	if !ok {
		return nil, false
	}

	response, err := http.ReadResponse(
		bufio.NewReaderSize(bytes.NewBuffer(dumpedResponse), len(dumpedResponse)),
		r,
	)

	if err != nil {
		return nil, false
	}

	return response, true
}

func (h *httpCache) Set(r *http.Request, res *http.Response) error {
	if res.ContentLength > 1024*1024 {
		return ErrCacheTooLarge
	}

	data, expired, err := h.cache.Get()
	if err != nil {
		return err
	}

	if expired || data == nil {
		data = make(map[string][]byte)
	}

	dumpedRequest, err := httputil.DumpRequest(r, true)
	if err != nil {
		return err
	}

	dumpedResponse, err := httputil.DumpResponse(res, true)
	if err != nil {
		return err
	}

	data[string(dumpedRequest)] = dumpedResponse
	return h.cache.Set(data)
}
