package http

import (
	"github.com/awirix/awirix/cache"
	"github.com/awirix/awirix/extensions"
	"github.com/awirix/lua"
	"net/http"
	"path/filepath"
)

var (
	cachers = make(map[string]*cache.HTTPCache)
)

func getExtensionWrapper(L *lua.LState) (wrapper extensions.ExtensionWrapper, ok bool) {
	wrapper, ok = L.Context().Value("extension").(extensions.ExtensionWrapper)
	return
}

func cacheGet(L *lua.LState, request *http.Request) (response *http.Response, ok bool) {
	wrapper, ok := getExtensionWrapper(L)
	if !ok {
		return nil, false
	}

	id := wrapper.Passport().ID
	if cacher, ok := cachers[id]; ok {
		return cacher.Get(request)
	}

	cachers[id] = cache.NewHTTPCache(filepath.Join(wrapper.Cache(), "http.json"))
	return nil, false
}

func cacheSet(L *lua.LState, request *http.Request, response *http.Response) (ok bool) {
	wrapper, ok := getExtensionWrapper(L)
	if !ok {
		return false
	}

	id := wrapper.Passport().ID
	if cacher, ok := cachers[id]; ok {
		return cacher.Set(request, response) == nil
	}

	cachers[id] = cache.NewHTTPCache(filepath.Join(wrapper.Cache(), "http.json"))
	return
}
