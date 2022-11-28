package passport

import (
	"github.com/vivi-app/vivi/passport"
	"github.com/vivi-app/vivi/util"
	lua "github.com/yuin/gopher-lua"
)

func New(L *lua.LState) *lua.LTable {
	registerPassportType(L)

	thePassport := getPassportFromContext(L)

	return util.NewTable(L, map[string]lua.LValue{
		"version": lua.LString(thePassport.Version.String()),
		"name":    lua.LString(thePassport.Name),
	}, map[string]lua.LGFunction{
		"config": passportConfig,
	})
}

const passportTypeName = "passport"

func registerPassportType(L *lua.LState) {
	mt := L.NewTypeMetatable(passportTypeName)
	L.SetGlobal(passportTypeName, mt)
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), passportMethods))
}

func pushPassport(L *lua.LState, passport *passport.Passport) {
	ud := L.NewUserData()
	ud.Value = passport
	L.SetMetatable(ud, L.GetTypeMetatable(passportTypeName))
	L.Push(ud)
}

func checkPassport(L *lua.LState, n int) *passport.Passport {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*passport.Passport); ok {
		return v
	}
	L.ArgError(n, "passport expected")
	return nil
}

func newPassport(L *lua.LState) int {
	thePassport := getPassportFromContext(L)
	pushPassport(L, thePassport)
	return 1
}

func getPassportFromContext(L *lua.LState) *passport.Passport {
	return L.Context().Value("passport").(*passport.Passport)
}

var passportMethods = map[string]lua.LGFunction{}

func passportConfig(L *lua.LState) int {
	thePassport := getPassportFromContext(L)
	key := L.CheckString(1)

	section, ok := thePassport.Config[key]
	if !ok {
		L.RaiseError("passport config section %s not found", key)
		return 1
	}

	value := section.Value()
	lvalue, err := util.ToLValue(L, value)
	if err != nil {
		L.RaiseError(err.Error())
		return 1
	}

	L.Push(lvalue)
	return 0
}
