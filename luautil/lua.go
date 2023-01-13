package luautil

import (
	"fmt"
	lua "github.com/awirix/lua"
	"reflect"
	"strconv"
)

func NewTable(L *lua.LState, fields map[string]lua.LValue, funcs map[string]lua.LGFunction) *lua.LTable {
	t := L.NewTable()
	for k, v := range fields {
		L.SetField(t, k, v)
	}

	if funcs != nil {
		L.SetFuncs(t, funcs)
	}

	return t
}

func FromLTable(table *lua.LTable) (map[string]any, error) {
	var err error
	result := make(map[string]any)
	table.ForEach(func(key lua.LValue, value lua.LValue) {
		if key.Type() == lua.LTString {
			result[string(key.(lua.LString))], err = FromLValue(value)
			if err != nil {
				return
			}
		}
	})

	return result, nil
}

func FromLValue(value lua.LValue) (any, error) {
	switch value.Type() {
	case lua.LTNil:
		return nil, nil
	case lua.LTBool:
		return bool(value.(lua.LBool)), nil
	case lua.LTNumber:
		return float64(value.(lua.LNumber)), nil
	case lua.LTString:
		return string(value.(lua.LString)), nil
	case lua.LTTable:
		return FromLTable(value.(*lua.LTable))
	case lua.LTUserData:
		return value.(*lua.LUserData).Value, nil
	default:
		return nil, fmt.Errorf("unsupported type: %s", value.Type())
	}
}

func ToLValue(L *lua.LState, value any) (lua.LValue, error) {
	switch value := value.(type) {
	case lua.LValue:
		return value, nil
	case nil:
		return lua.LNil, nil
	case bool:
		return lua.LBool(value), nil
	case uint, uint8, uint16, uint32, uint64, int, int8, int16, int32, int64, float32, float64:
		f, err := strconv.ParseFloat(fmt.Sprint(value), 64)
		if err != nil {
			return nil, err
		}

		return lua.LNumber(f), nil
	case string:
		return lua.LString(value), nil
	default:
		switch reflect.TypeOf(value).Kind() {
		case reflect.Slice:
			table := L.NewTable()

			// iterate over the slice
			slice := reflect.ValueOf(value)
			for i := 0; i < slice.Len(); i++ {
				lvalue, err := ToLValue(L, slice.Index(i).Interface())
				if err != nil {
					return nil, err
				}

				L.RawSetInt(table, i+1, lvalue)
			}

			return table, nil
		case reflect.Map:
			table := L.NewTable()

			// iterate over the map
			m := reflect.ValueOf(value)
			for _, key := range m.MapKeys() {
				lkey, err := ToLValue(L, key.Interface())
				if err != nil {
					return nil, err
				}

				lvalue, err := ToLValue(L, m.MapIndex(key).Interface())
				if err != nil {
					return nil, err
				}

				L.SetTable(table, lkey, lvalue)
			}

			return table, nil
		default:
			return nil, fmt.Errorf("unsupported type: %T", value)
		}
	}
}
