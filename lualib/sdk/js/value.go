package js

import (
	"github.com/awirix/awirix/luautil"
	lua "github.com/awirix/lua"
	"github.com/robertkrimen/otto"
)

const vmValueTypeName = vmTypeName + "_value"

func pushVMValue(L *lua.LState, value *otto.Value) {
	ud := L.NewUserData()
	ud.Value = value
	L.SetMetatable(ud, L.GetTypeMetatable(vmValueTypeName))
	L.Push(ud)
}

func checkVMValue(L *lua.LState, n int) *otto.Value {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*otto.Value); ok {
		return v
	}

	L.ArgError(n, "js_vm_value expected")
	return nil
}

func vmValueExport(L *lua.LState) int {
	value := checkVMValue(L, 1)
	nativeValue, err := value.Export()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	lvalue, err := luautil.ToLValue(L, nativeValue)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lvalue)
	return 1
}

func vmValueIsNull(L *lua.LState) int {
	value := checkVMValue(L, 1)
	L.Push(lua.LBool(value.IsNull()))
	return 1
}

func vmValueIsUndefined(L *lua.LState) int {
	value := checkVMValue(L, 1)
	L.Push(lua.LBool(value.IsUndefined()))
	return 1
}

func vmValueIsNumber(L *lua.LState) int {
	value := checkVMValue(L, 1)
	L.Push(lua.LBool(value.IsNumber()))
	return 1
}

func vmValueIsBoolean(L *lua.LState) int {
	value := checkVMValue(L, 1)
	L.Push(lua.LBool(value.IsBoolean()))
	return 1
}

func vmValueIsString(L *lua.LState) int {
	value := checkVMValue(L, 1)
	L.Push(lua.LBool(value.IsString()))
	return 1
}

func vmValueIsObject(L *lua.LState) int {
	value := checkVMValue(L, 1)
	L.Push(lua.LBool(value.IsObject()))
	return 1
}

func vmValueIsNaN(L *lua.LState) int {
	value := checkVMValue(L, 1)
	L.Push(lua.LBool(value.IsNaN()))
	return 1
}

func vmValueIsFunction(L *lua.LState) int {
	value := checkVMValue(L, 1)
	L.Push(lua.LBool(value.IsFunction()))
	return 1
}

func vmValueString(L *lua.LState) int {
	value := checkVMValue(L, 1)
	L.Push(lua.LString(value.String()))
	return 1
}
