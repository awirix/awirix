package html

import (
	"github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	"github.com/awirix/lua"
	"github.com/cixtor/readability"
	"strings"
)

const selectionTypeName = documentTypeName + "_selection"

func checkSelection(L *lua.LState, n int) *goquery.Selection {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*goquery.Selection); ok {
		return v
	}
	L.ArgError(n, "selection expected")
	return nil
}

func pushSelection(L *lua.LState, selection *goquery.Selection) {
	ud := L.NewUserData()
	ud.Value = selection
	L.SetMetatable(ud, L.GetTypeMetatable(selectionTypeName))
	L.Push(ud)
}

func selectionFind(L *lua.LState) int {
	selection := checkSelection(L, 1)
	selector := L.CheckString(2)

	selection = selection.Find(selector)
	pushSelection(L, selection)
	return 1
}

func selectionEach(L *lua.LState) int {
	selection := checkSelection(L, 1)
	callback := L.CheckFunction(2)

	pushSelection(L, selection.Each(func(i int, selection *goquery.Selection) {
		L.Push(callback)
		pushSelection(L, selection)
		L.Push(lua.LNumber(i))

		if err := L.PCall(2, lua.MultRet, nil); err != nil {
			L.RaiseError(err.Error())
			return
		}
	}))

	return 1
}

func selectionText(L *lua.LState) int {
	selection := checkSelection(L, 1)
	L.Push(lua.LString(selection.Text()))
	return 1
}

func selectionHtml(L *lua.LState) int {
	selection := checkSelection(L, 1)
	html, err := selection.Html()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LString(html))
	return 1
}

func selectionFirst(L *lua.LState) int {
	selection := checkSelection(L, 1)
	pushSelection(L, selection.First())
	return 1
}

func selectionLast(L *lua.LState) int {
	selection := checkSelection(L, 1)
	pushSelection(L, selection.Last())
	return 1
}

func selectionParent(L *lua.LState) int {
	selection := checkSelection(L, 1)
	pushSelection(L, selection.Parent())
	return 1
}

func selectionEq(L *lua.LState) int {
	selection := checkSelection(L, 1)
	index := L.CheckInt(2)
	pushSelection(L, selection.Eq(index))
	return 1
}

func selectionAttr(L *lua.LState) int {
	selection := checkSelection(L, 1)
	key := L.CheckString(2)
	attr, exists := selection.Attr(key)
	L.Push(lua.LString(attr))
	L.Push(lua.LBool(exists))

	return 2
}

func selectionAttrOr(L *lua.LState) int {
	selection := checkSelection(L, 1)
	key := L.CheckString(2)
	defaultValue := L.CheckString(3)
	attr, exists := selection.Attr(key)
	if !exists {
		L.Push(lua.LString(defaultValue))
		return 1
	}

	L.Push(lua.LString(attr))
	return 1
}

func selectionHasClass(L *lua.LState) int {
	selection := checkSelection(L, 1)
	class := L.CheckString(2)
	L.Push(lua.LBool(selection.HasClass(class)))
	return 1
}

func selectionAddClass(L *lua.LState) int {
	selection := checkSelection(L, 1)
	class := L.CheckString(2)
	selection.AddClass(class)
	return 0
}

func selectionRemoveClass(L *lua.LState) int {
	selection := checkSelection(L, 1)
	class := L.CheckString(2)
	selection.RemoveClass(class)
	return 0
}

func selectionToggleClass(L *lua.LState) int {
	selection := checkSelection(L, 1)
	class := L.CheckString(2)
	selection.ToggleClass(class)
	return 0
}

func selectionNext(L *lua.LState) int {
	selection := checkSelection(L, 1)
	pushSelection(L, selection.Next())
	return 1
}

func selectionNextAll(L *lua.LState) int {
	selection := checkSelection(L, 1)
	pushSelection(L, selection.NextAll())
	return 1
}

func selectionNextUntil(L *lua.LState) int {
	selection := checkSelection(L, 1)
	selector := L.CheckString(2)
	pushSelection(L, selection.NextUntil(selector))
	return 1
}

func selectionPrev(L *lua.LState) int {
	selection := checkSelection(L, 1)
	pushSelection(L, selection.Prev())
	return 1
}

func selectionPrevAll(L *lua.LState) int {
	selection := checkSelection(L, 1)
	pushSelection(L, selection.PrevAll())
	return 1
}

func selectionPrevUntil(L *lua.LState) int {
	selection := checkSelection(L, 1)
	selector := L.CheckString(2)
	pushSelection(L, selection.PrevUntil(selector))
	return 1
}

func selectionSiblings(L *lua.LState) int {
	selection := checkSelection(L, 1)
	pushSelection(L, selection.Siblings())
	return 1
}

func selectionChildren(L *lua.LState) int {
	selection := checkSelection(L, 1)
	pushSelection(L, selection.Children())
	return 1
}

func selectionContents(L *lua.LState) int {
	selection := checkSelection(L, 1)
	pushSelection(L, selection.Contents())
	return 1
}

func selectionFilter(L *lua.LState) int {
	selection := checkSelection(L, 1)
	selector := L.CheckString(2)
	pushSelection(L, selection.Filter(selector))
	return 1
}

func selectionRemove(L *lua.LState) int {
	selection := checkSelection(L, 1)
	selector := L.CheckString(2)
	pushSelection(L, selection.Not(selector))
	return 1
}

func selectionIs(L *lua.LState) int {
	selection := checkSelection(L, 1)
	selector := L.CheckString(2)
	L.Push(lua.LBool(selection.Is(selector)))
	return 1
}

func selectionFindSelection(L *lua.LState) int {
	selection := checkSelection(L, 1)
	other := checkSelection(L, 2)
	pushSelection(L, selection.FindSelection(other))
	return 1
}

func selectionMap(L *lua.LState) int {
	selection := checkSelection(L, 1)
	callback := L.CheckFunction(2)

	table := L.NewTable()

	selection.Each(func(i int, selection *goquery.Selection) {
		L.Push(callback)
		pushSelection(L, selection)
		if err := L.PCall(1, lua.MultRet, nil); err != nil {
			L.RaiseError(err.Error())
			return
		}

		value := L.Get(-1)
		L.Pop(1)

		table.RawSetInt(i, value)
	})

	L.Push(table)
	return 1
}

func selectionClosest(L *lua.LState) int {
	selection := checkSelection(L, 1)
	selector := L.CheckString(2)
	pushSelection(L, selection.Closest(selector))
	return 1
}

func selectionParents(L *lua.LState) int {
	selection := checkSelection(L, 1)
	pushSelection(L, selection.Parents())
	return 1
}

func selectionParentsUntil(L *lua.LState) int {
	selection := checkSelection(L, 1)
	selector := L.CheckString(2)
	pushSelection(L, selection.ParentsUntil(selector))
	return 1
}

func selectionSlice(L *lua.LState) int {
	selection := checkSelection(L, 1)
	start := L.CheckInt(2)
	end := L.CheckInt(3)
	pushSelection(L, selection.Slice(start, end))
	return 1
}

func selectionTerminate(L *lua.LState) int {
	selection := checkSelection(L, 1)
	pushSelection(L, selection.End())
	return 1
}

func selectionAdd(L *lua.LState) int {
	selection := checkSelection(L, 1)
	selector := L.CheckString(2)
	pushSelection(L, selection.Add(selector))
	return 1
}

func selectionLength(L *lua.LState) int {
	selection := checkSelection(L, 1)
	L.Push(lua.LNumber(selection.Length()))
	return 1
}

func selectionAddSelection(L *lua.LState) int {
	selection := checkSelection(L, 1)
	other := checkSelection(L, 2)
	pushSelection(L, selection.AddSelection(other))
	return 1
}

func selectionAddBack(L *lua.LState) int {
	selection := checkSelection(L, 1)
	pushSelection(L, selection.AddBack())
	return 1
}

func selectionMarkdown(L *lua.LState) int {
	selection := checkSelection(L, 1)
	converter := md.NewConverter("", true, nil)

	L.Push(lua.LString(converter.Convert(selection)))
	return 1
}

func selectionSimplified(L *lua.LState) int {
	selection := checkSelection(L, 1)

	html, err := selection.Html()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	article, err := readability.New().Parse(strings.NewReader(html), "https://example.com")
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	document := goquery.NewDocumentFromNode(article.Node)

	pushSelection(L, document.Selection)
	return 1
}
