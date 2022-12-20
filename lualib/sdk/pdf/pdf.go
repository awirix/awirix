package pdf

import (
	"bytes"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/vivi-app/lua"
	"github.com/vivi-app/vivi/luautil"
	"io"
	"strings"
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
		if value.Type() != lua.LTString {
			L.RaiseError("pdf.from_images: expected a table of strings")
			return
		}

		readers = append(
			readers,
			strings.NewReader(value.String()),
		)
	})

	pdf, err := convertImagesToPDF(readers)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	contents, err := io.ReadAll(pdf)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LString(contents))
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
