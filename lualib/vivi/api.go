package vivi

import (
	"fmt"
	lua "github.com/yuin/gopher-lua"
)

func Watch(L *lua.LState) int {
	url := L.CheckString(1)

	fmt.Printf("watching %s\n", url)

	return 0
}
