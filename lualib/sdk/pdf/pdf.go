package pdf

import (
	"bytes"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	lua "github.com/vivi-app/lua"
	io2 "github.com/vivi-app/vivi/lualib/sdk/io"
	"github.com/vivi-app/vivi/luautil"
	"io"
)

func New(L *lua.LState) *lua.LTable {
	return luautil.NewTable(L, nil, map[string]lua.LGFunction{
		"from_images": fromImages,
	})
}

func fromImages(L *lua.LState) int {
	images := L.CheckTable(1)
	readers := make([]io.Reader, 0)

	images.ForEach(func(key, value lua.LValue) {
		reader, ok := value.(*lua.LUserData).Value.(io.Reader)
		if !ok {
			L.RaiseError("pdf.from_images: expected a table of io.Reader")
			return
		}

		readers = append(readers, reader)
	})

	pdf, err := convertImagesToPDF(readers)
	if err != nil {
		L.RaiseError(err.Error())
		return 0
	}

	io2.PushReadCloser(L, io.NopCloser(pdf))
	return 1
}

func convertImagesToPDF(images []io.Reader) (io.Reader, error) {
	var pdf bytes.Buffer
	err := api.ImportImages(nil, &pdf, images, nil, nil)
	if err != nil {
		return nil, err
	}

	return &pdf, nil
}
