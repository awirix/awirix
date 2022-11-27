package html

import (
	"github.com/PuerkitoBio/goquery"
	lua "github.com/yuin/gopher-lua"
	"strings"
)

const documentTypeName = "document"

func registerDocumentType(L *lua.LState) {
	mt := L.NewTypeMetatable(documentTypeName)
	L.SetGlobal(documentTypeName, mt)
}

func checkDocument(L *lua.LState) *goquery.Document {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*goquery.Document); ok {
		return v
	}
	L.ArgError(1, "document expected")
	return nil
}

func parse(L *lua.LState) int {
	value := L.CheckString(1)
	reader := strings.NewReader(value)
	document, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	ud := L.NewUserData()
	ud.Value = document
	L.SetMetatable(ud, L.GetTypeMetatable(documentTypeName))
	L.Push(ud)

	return 1
}
