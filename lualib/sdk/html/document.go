package html

import (
	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	lua "github.com/vivi-app/lua"
	"strings"
)

const documentTypeName = "html_document"

var documentMethods = map[string]lua.LGFunction{
	"find":     documentFind,
	"html":     documentHtml,
	"markdown": documentMarkdown,
}

func registerDocumentType(L *lua.LState) {
	mt := L.NewTypeMetatable(documentTypeName)
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), documentMethods))
}

func checkDocument(L *lua.LState, n int) *goquery.Document {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*goquery.Document); ok {
		return v
	}
	L.ArgError(1, "document expected")
	return nil
}

func pushDocument(L *lua.LState, document *goquery.Document) {
	ud := L.NewUserData()
	ud.Value = document
	L.SetMetatable(ud, L.GetTypeMetatable(documentTypeName))
	L.Push(ud)
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

	pushDocument(L, document)
	return 1
}

func documentFind(L *lua.LState) int {
	document := checkDocument(L, 1)
	selector := L.CheckString(2)

	selection := document.Find(selector)
	pushSelection(L, selection)
	return 1
}

func documentHtml(L *lua.LState) int {
	document := checkDocument(L, 1)
	html, err := document.Html()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LString(html))
	return 1
}

func documentSelection(L *lua.LState) int {
	document := checkDocument(L, 1)
	selection := document.Selection
	pushSelection(L, selection)
	return 1
}

func documentMarkdown(L *lua.LState) int {
	document := checkDocument(L, 1)
	converter := md.NewConverter("", true, nil)

	L.Push(lua.LString(converter.Convert(document.Selection)))
	return 1
}
