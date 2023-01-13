package functional

import (
	"github.com/awirix/awirix/luadoc"
	"github.com/awirix/awirix/luautil"
	"github.com/awirix/lua"
	"github.com/samber/lo"
	"golang.org/x/exp/slices"
)

func Lib() *luadoc.Lib {
	predicate := &luadoc.Param{
		Name:        "predicate",
		Description: "Predicate function",
		Type: luadoc.Func{
			Params: []*luadoc.Param{
				{
					Name:        "value",
					Description: "Value to check",
					Type:        luadoc.Any,
				},
			},
			Returns: []*luadoc.Param{
				{
					Name:        "result",
					Description: "Result of predicate",
					Type:        luadoc.Boolean,
				},
			},
		}.AsType(),
	}

	predicateWithIndex := &luadoc.Param{
		Name:        "predicate",
		Description: "Predicate function",
		Type: luadoc.Func{
			Params: []*luadoc.Param{
				{
					Name:        "value",
					Description: "Value to check",
					Type:        luadoc.Any,
				},
				{
					Name:        "index",
					Description: "Index of value",
					Type:        luadoc.Number,
				},
			},
			Returns: []*luadoc.Param{
				{
					Name:        "result",
					Description: "Result of predicate",
					Type:        luadoc.Boolean,
				},
			},
		}.AsType(),
	}

	return &luadoc.Lib{
		Name:        "functional",
		Description: "Functional helpers",
		Funcs: []*luadoc.Func{
			{
				Name:        "chunk",
				Description: "Split list into chunks",
				Value:       chunk,
				Params: []*luadoc.Param{
					{
						Name:        "list",
						Description: "List to split",
						Type:        luadoc.List(luadoc.Any),
					},
					{
						Name:        "size",
						Description: "Chunk size",
						Type:        luadoc.Number,
						Optional:    true,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "chunks",
						Description: "List of chunks",
						Type:        luadoc.List(luadoc.List(luadoc.Any)),
					},
				},
			},
			{
				Name:        "drop",
				Description: "Drop n elements from the beginning list",
				Value:       drop,
				Params: []*luadoc.Param{
					{
						Name:        "list",
						Description: "List to drop from",
						Type:        luadoc.List(luadoc.Any),
					},
					{
						Name:        "n",
						Description: "Number of elements to drop",
						Type:        luadoc.Number,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "dropped",
						Description: "List of dropped elements",
						Type:        luadoc.List(luadoc.Any),
					},
				},
			},
			{
				Name:        "drop_right",
				Description: "Drop n elements from the end of list",
				Value:       dropRight,
				Params: []*luadoc.Param{
					{
						Name:        "list",
						Description: "List to drop from",
						Type:        luadoc.List(luadoc.Any),
					},
					{
						Name:        "n",
						Description: "Number of elements to drop",
						Type:        luadoc.Number,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "dropped",
						Description: "List of dropped elements",
						Type:        luadoc.List(luadoc.Any),
					},
				},
			},
			{
				Name:        "drop_while",
				Description: "Drop elements from the beginning of list while predicate is true",
				Value:       dropWhile,
				Params: []*luadoc.Param{
					{
						Name:        "list",
						Description: "List to drop from",
						Type:        luadoc.List(luadoc.Any),
					},
					predicate,
				},
			},
			{
				Name:        "drop_right_while",
				Description: "Drop elements from the end of list while predicate is true",
				Value:       dropRightWhile,
				Params: []*luadoc.Param{
					{
						Name:        "list",
						Description: "List to drop from",
						Type:        luadoc.List(luadoc.Any),
					},
					predicate,
				},
			},
			{
				Name:        "filter",
				Description: "Filter list by predicate",
				Value:       filter,
				Params: []*luadoc.Param{
					{
						Name:        "list",
						Description: "List to filter",
						Type:        luadoc.List(luadoc.Any),
					},
					predicateWithIndex,
				},
				Returns: []*luadoc.Param{
					{
						Name:        "filtered",
						Description: "List of filtered elements",
						Type:        luadoc.List(luadoc.Any),
					},
				},
			},
			{
				Name:        "find",
				Description: "Find first element in list that satisfies predicate.",
				Value:       find,
				Params: []*luadoc.Param{
					{
						Name:        "list",
						Description: "List to search",
						Type:        luadoc.List(luadoc.Any),
					},
					predicate,
				},
				Returns: []*luadoc.Param{
					{
						Name:        "element",
						Description: "First element that satisfies predicate",
						Type:        luadoc.Any,
					},
					{
						Name:        "ok",
						Description: "True if element was found",
						Type:        luadoc.Boolean,
					},
				},
			},
			{
				Name:        "find_index",
				Description: "Find index of first element in list that satisfies predicate.",
				Value:       findIndex,
				Params: []*luadoc.Param{
					{
						Name:        "list",
						Description: "List to search",
						Type:        luadoc.List(luadoc.Any),
					},
					predicate,
				},
				Returns: []*luadoc.Param{
					{
						Name:        "index",
						Description: "Index of first element that satisfies predicate",
						Type:        luadoc.Number,
					},
					{
						Name:        "ok",
						Description: "True if element was found",
						Type:        luadoc.Boolean,
					},
				},
			},
			{
				Name:        "find_last_index",
				Description: "Find index of last element in list that satisfies predicate.",
				Value:       findLastIndex,
				Params: []*luadoc.Param{
					{
						Name:        "list",
						Description: "List to search",
						Type:        luadoc.List(luadoc.Any),
					},
					predicate,
				},
				Returns: []*luadoc.Param{
					{
						Name:        "index",
						Description: "Index of last element that satisfies predicate",
						Type:        luadoc.Number,
					},
					{
						Name:        "ok",
						Description: "True if element was found",
						Type:        luadoc.Boolean,
					},
				},
			},
			{
				Name:        "head",
				Description: "Get first element of list",
				Value:       head,
				Params: []*luadoc.Param{
					{
						Name:        "list",
						Description: "List to get first element from",
						Type:        luadoc.List(luadoc.Any),
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "element",
						Description: "First element of list",
						Type:        luadoc.Any,
						Optional:    true,
					},
				},
			},
			{
				Name:        "tail",
				Description: "Get all elements of list except first",
				Value:       tail,
				Params: []*luadoc.Param{
					{
						Name:        "list",
						Description: "List to get tail from",
						Type:        luadoc.List(luadoc.Any),
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "tail",
						Description: "List of all elements except first",
						Type:        luadoc.List(luadoc.Any),
						Optional:    true,
					},
				},
			},
			{
				Name:        "last",
				Description: "Get last element of list",
				Value:       last,
				Params: []*luadoc.Param{
					{
						Name:        "list",
						Description: "List to get last element from",
						Type:        luadoc.List(luadoc.Any),
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "element",
						Description: "Last element of list",
						Type:        luadoc.Any,
						Optional:    true,
					},
				},
			},
			{
				Name:        "init",
				Description: "Get all elements of list except last",
				Value:       initial,
				Params: []*luadoc.Param{
					{
						Name:        "list",
						Description: "List to get initial from",
						Type:        luadoc.List(luadoc.Any),
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "initial",
						Description: "List of all elements except last",
						Type:        luadoc.List(luadoc.Any),
						Optional:    true,
					},
				},
			},
			{
				Name:        "take",
				Description: "Get first n elements of list",
				Value:       take,
				Params: []*luadoc.Param{
					{
						Name:        "list",
						Description: "List to get elements from",
						Type:        luadoc.List(luadoc.Any),
					},
					{
						Name:        "n",
						Description: "Number of elements to get",
						Type:        luadoc.Number,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "elements",
						Description: "First n elements of list",
						Type:        luadoc.List(luadoc.Any),
						Optional:    true,
					},
				},
			},
			{
				Name:        "drop",
				Description: "Get all elements of list except first n",
				Value:       drop,
				Params: []*luadoc.Param{
					{
						Name:        "list",
						Description: "List to get elements from",
						Type:        luadoc.List(luadoc.Any),
					},
					{
						Name:        "n",
						Description: "Number of elements to drop",
						Type:        luadoc.Number,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "elements",
						Description: "All elements of list except first n",
						Type:        luadoc.List(luadoc.Any),
						Optional:    true,
					},
				},
			},
			{
				Name:        "drop_while",
				Description: "Drop elements from list while predicate is true",
				Value:       dropWhile,
				Params: []*luadoc.Param{
					{
						Name:        "list",
						Description: "List to get elements from",
						Type:        luadoc.List(luadoc.Any),
					},
					predicate,
				},
				Returns: []*luadoc.Param{
					{
						Name:        "elements",
						Description: "Elements of list after predicate is false",
						Type:        luadoc.List(luadoc.Any),
						Optional:    true,
					},
				},
			},
			{
				Name:        "take_while",
				Description: "Take elements from list while predicate is true",
				Value:       takeWhile,
				Params: []*luadoc.Param{
					{
						Name:        "list",
						Description: "List to get elements from",
						Type:        luadoc.List(luadoc.Any),
					},
					predicate,
				},
				Returns: []*luadoc.Param{
					{
						Name:        "elements",
						Description: "Elements of list before predicate is false",
						Type:        luadoc.List(luadoc.Any),
						Optional:    true,
					},
				},
			},
			{
				Name:        "reverse",
				Description: "Reverse list",
				Value:       reverse,
				Params: []*luadoc.Param{
					{
						Name:        "list",
						Description: "List to reverse",
						Type:        luadoc.List(luadoc.Any),
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "reversed",
						Description: "Reversed list",
						Type:        luadoc.List(luadoc.Any),
					},
				},
			},
			{
				Name:        "slice",
				Description: "Get slice of list",
				Value:       slice,
				Params: []*luadoc.Param{
					{
						Name:        "list",
						Description: "List to get slice from",
						Type:        luadoc.List(luadoc.Any),
					},
					{
						Name:        "start",
						Description: "Start index of slice",
						Type:        luadoc.Number,
					},
					{
						Name:        "finish",
						Description: "End index of slice. Defaults to length of list",
						Type:        luadoc.Number,
						Optional:    true,
					},
				},
			},
			{
				Name:        "sorted",
				Description: "Returns a sorted list",
				Value:       sorted,
				Params: []*luadoc.Param{
					{
						Name:        "list",
						Description: "List to sort",
						Type:        luadoc.List(luadoc.Any),
					},
					{
						Name:        "cmp",
						Description: "Comparison function. Defaults to <",
						Type: luadoc.Func{
							Params: []*luadoc.Param{
								{
									Name: "a",
									Type: luadoc.Any,
								},
								{
									Name: "b",
									Type: luadoc.Any,
								},
							},
							Returns: []*luadoc.Param{
								{
									Name: "result",
									Type: luadoc.Boolean,
								},
							},
						}.AsType(),
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "sorted",
						Description: "Sorted list",
						Type:        luadoc.List(luadoc.Any),
					},
				},
			},
			{
				Name:        "zip2",
				Description: "Zips two lists together",
				Value:       zip2,
				Params: []*luadoc.Param{
					{
						Name:        "list1",
						Description: "First list to zip",
						Type:        luadoc.List(luadoc.Any),
					},
					{
						Name:        "list2",
						Description: "Second list to zip",
						Type:        luadoc.List(luadoc.Any),
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "zipped",
						Description: "Zipped list",
						Type:        luadoc.List(luadoc.List(luadoc.Any)),
					},
				},
			},
			{
				Name:        "zip3",
				Description: "Zips three lists together",
				Value:       zip3,
				Params: []*luadoc.Param{
					{
						Name:        "list1",
						Description: "First list to zip",
						Type:        luadoc.List(luadoc.Any),
					},
					{
						Name:        "list2",
						Description: "Second list to zip",
						Type:        luadoc.List(luadoc.Any),
					},
					{
						Name:        "list3",
						Description: "Third list to zip",
						Type:        luadoc.List(luadoc.Any),
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "zipped",
						Description: "Zipped list",
						Type:        luadoc.List(luadoc.List(luadoc.Any)),
					},
				},
			},
			{
				Name:        "map",
				Description: "Maps a function over a list",
				Value:       map_,
				Params: []*luadoc.Param{
					{
						Name:        "list",
						Description: "List to map over",
						Type:        luadoc.List(luadoc.Any),
					},
					{
						Name:        "func",
						Description: "Function to map",
						Type: luadoc.Func{
							Params: []*luadoc.Param{
								{
									Name: "value",
									Type: luadoc.Any,
								},
								{
									Name: "index",
									Type: luadoc.Number,
								},
							},
							Returns: []*luadoc.Param{
								{
									Name: "result",
									Type: luadoc.Any,
								},
							},
						}.AsType(),
					},
				},
			},
			{
				Name:        "for_each",
				Description: "Calls a function for each element in a list",
				Value:       forEach,
				Params: []*luadoc.Param{
					{
						Name:        "list",
						Description: "List to iterate over",
						Type:        luadoc.List(luadoc.Any),
					},
					{
						Name:        "func",
						Description: "Function to call",
						Type: luadoc.Func{
							Params: []*luadoc.Param{
								{
									Name:        "value",
									Description: "Value of the element",
									Type:        luadoc.Any,
								},
								{
									Name:        "index",
									Description: "Index of the element in the list",
								},
							},
						}.AsType(),
					},
				},
			},
			{
				Name:        "contains_by",
				Description: "Checks if a list contains a value using predicate",
				Value:       containsBy,
				Params: []*luadoc.Param{
					{
						Name:        "list",
						Description: "List to check",
						Type:        luadoc.List(luadoc.Any),
					},
					predicate,
				},
				Returns: []*luadoc.Param{
					{
						Name:        "contains",
						Description: "Whether the list contains the value",
						Type:        luadoc.Boolean,
					},
				},
			},
			{
				Name:        "reduce",
				Description: "Reduces a list to a single value",
				Value:       reduce,
				Params: []*luadoc.Param{
					{
						Name:        "list",
						Description: "List to reduce",
						Type:        luadoc.List(luadoc.Any),
					},
					{
						Name:        "func",
						Description: "Function to reduce with",
						Type: luadoc.Func{
							Params: []*luadoc.Param{
								{
									Name: "accumulator",
									Type: luadoc.Any,
								},
								{
									Name: "value",
									Type: luadoc.Any,
								},
								{
									Name: "initial",
									Type: luadoc.Any,
								},
							},
							Returns: []*luadoc.Param{
								{
									Name: "result",
									Type: luadoc.Any,
								},
							},
						}.AsType(),
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "result",
						Description: "Result of the reduction",
						Type:        luadoc.Any,
					},
				},
			},
			{
				Name:        "reduce_right",
				Description: "Reduces a list to a single value",
				Value:       reduceRight,
				Params: []*luadoc.Param{
					{
						Name:        "list",
						Description: "List to reduce",
						Type:        luadoc.List(luadoc.Any),
					},
					{
						Name:        "func",
						Description: "Function to reduce with",
						Type: luadoc.Func{
							Params: []*luadoc.Param{
								{
									Name: "accumulator",
									Type: luadoc.Any,
								},
								{
									Name: "value",
									Type: luadoc.Any,
								},
								{
									Name: "initial",
									Type: luadoc.Any,
								},
							},
							Returns: []*luadoc.Param{
								{
									Name: "result",
									Type: luadoc.Any,
								},
							},
						}.AsType(),
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "result",
						Description: "Result of the reduction",
						Type:        luadoc.Any,
					},
				},
			},
			{
				Name:        "reject",
				Description: "Rejects all elements that match a predicate",
				Value:       reject,
				Params: []*luadoc.Param{
					{
						Name:        "list",
						Description: "List to reject from",
						Type:        luadoc.List(luadoc.Any),
					},
					{
						Name:        "func",
						Description: "Predicate to reject with",
						Type: luadoc.Func{
							Params: []*luadoc.Param{
								{
									Name: "value",
									Type: luadoc.Any,
								},
								{
									Name: "index",
									Type: luadoc.Number,
								},
							},
							Returns: []*luadoc.Param{
								{
									Name: "result",
									Type: luadoc.Boolean,
								},
							},
						}.AsType(),
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "result",
						Description: "List of rejected elements",
						Type:        luadoc.List(luadoc.Any),
					},
				},
			},
			{
				Name:        "id",
				Description: "Returns the first argument",
				Value:       id,
				Params: []*luadoc.Param{
					{
						Name:        "value",
						Description: "Value to return",
						Type:        luadoc.Any,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "result",
						Description: "The first argument",
						Type:        luadoc.Any,
					},
				},
			},
			// TODO: do the rest
		},
	}
}

func checkList(L *lua.LState, n int) (list []lua.LValue) {
	table := L.CheckTable(n)
	table.ForEach(func(key lua.LValue, value lua.LValue) {
		list = append(list, value)
	})
	return
}

func chunk(L *lua.LState) int {
	list := checkList(L, 1)
	size := L.OptInt(2, 1)

	chunks, _ := luautil.ToLValue(L, lo.Chunk(list, size))
	L.Push(chunks)
	return 1
}

func drop(L *lua.LState) int {
	list := checkList(L, 1)
	n := L.CheckInt(2)

	dropped, _ := luautil.ToLValue(L, lo.Drop(list, n))
	L.Push(dropped)
	return 1
}

func dropRight(L *lua.LState) int {
	list := checkList(L, 1)
	n := L.CheckInt(2)

	dropped, _ := luautil.ToLValue(L, lo.DropRight(list, n))
	L.Push(dropped)
	return 1
}

func dropWhile(L *lua.LState) int {
	list := checkList(L, 1)
	fn := L.CheckFunction(2)

	dropped, _ := luautil.ToLValue(L, lo.DropWhile(list, func(v lua.LValue) bool {
		L.Push(fn)
		L.Push(v)
		L.Call(1, 1)
		return L.ToBool(-1)
	}))
	L.Push(dropped)
	return 1
}

func dropRightWhile(L *lua.LState) int {
	list := checkList(L, 1)
	fn := L.CheckFunction(2)

	dropped, _ := luautil.ToLValue(L, lo.DropRightWhile(list, func(v lua.LValue) bool {
		L.Push(fn)
		L.Push(v)
		L.Call(1, 1)
		return L.ToBool(-1)
	}))
	L.Push(dropped)
	return 1
}

func filter(L *lua.LState) int {
	list := checkList(L, 1)
	fn := L.CheckFunction(2)

	filtered, _ := luautil.ToLValue(L, lo.Filter(list, func(v lua.LValue, i int) bool {
		L.Push(fn)
		L.Push(v)
		L.Push(lua.LNumber(i))
		L.Call(2, 1)
		return L.ToBool(-1)
	}))
	L.Push(filtered)
	return 1
}

func find(L *lua.LState) int {
	list := checkList(L, 1)
	fn := L.CheckFunction(2)
	found, ok := lo.Find(list, func(v lua.LValue) bool {
		L.Push(fn)
		L.Push(v)
		L.Call(1, 1)
		return L.ToBool(-1)
	})

	lvalue, _ := luautil.ToLValue(L, found)
	L.Push(lvalue)
	L.Push(lua.LBool(ok))
	return 2
}

func findIndex(L *lua.LState) int {
	list := checkList(L, 1)
	fn := L.CheckFunction(2)
	_, index, ok := lo.FindIndexOf(list, func(v lua.LValue) bool {
		L.Push(fn)
		L.Push(v)
		L.Call(1, 1)
		return L.ToBool(-1)
	})

	L.Push(lua.LNumber(index))
	L.Push(lua.LBool(ok))
	return 2
}

func findLastIndex(L *lua.LState) int {
	list := checkList(L, 1)
	fn := L.CheckFunction(2)
	_, index, ok := lo.FindLastIndexOf(list, func(v lua.LValue) bool {
		L.Push(fn)
		L.Push(v)
		L.Call(1, 1)
		return L.ToBool(-1)
	})

	L.Push(lua.LNumber(index))
	L.Push(lua.LBool(ok))

	return 2
}

func head(L *lua.LState) int {
	list := checkList(L, 1)

	var head lua.LValue
	if len(list) == 0 {
		head = nil
	} else {
		head = list[0]
	}

	lvalue, _ := luautil.ToLValue(L, head)
	L.Push(lvalue)
	return 1
}

func tail(L *lua.LState) int {
	list := checkList(L, 1)

	var tail []lua.LValue
	if len(list) == 0 {
		tail = nil
	} else {
		tail = list[1:]
	}

	lvalue, _ := luautil.ToLValue(L, tail)
	L.Push(lvalue)
	return 1
}

func initial(L *lua.LState) int {
	list := checkList(L, 1)

	var initial []lua.LValue
	if len(list) == 0 {
		initial = nil
	} else {
		initial = list[:len(list)-1]
	}

	lvalue, _ := luautil.ToLValue(L, initial)
	L.Push(lvalue)
	return 1
}

func take(L *lua.LState) int {
	list := checkList(L, 1)
	n := L.CheckInt(2)

	var taken []lua.LValue
	if len(list) == 0 {
		taken = nil
	} else {
		taken = list[:n]
	}

	t, _ := luautil.ToLValue(L, taken)
	L.Push(t)
	return 1
}

func takeWhile(L *lua.LState) int {
	list := checkList(L, 1)
	fn := L.CheckFunction(2)

	for i, v := range list {
		L.Push(fn)
		L.Push(v)
		L.Call(1, 1)
		if !L.ToBool(-1) {
			list = list[:i]
			break
		}
	}

	l, _ := luautil.ToLValue(L, list)
	L.Push(l)
	return 1
}

func last(L *lua.LState) int {
	list := checkList(L, 1)

	var last lua.LValue
	if len(list) == 0 {
		last = nil
	} else {
		last = list[len(list)-1]
	}

	lvalue, _ := luautil.ToLValue(L, last)
	L.Push(lvalue)
	return 1
}

func reverse(L *lua.LState) int {
	list := checkList(L, 1)

	reversed, _ := luautil.ToLValue(L, lo.Reverse(list))
	L.Push(reversed)
	return 1
}

func slice(L *lua.LState) int {
	list := checkList(L, 1)
	start := L.CheckInt(2)
	end := L.OptInt(3, len(list))

	sliced, _ := luautil.ToLValue(L, lo.Slice(list, start, end))
	L.Push(sliced)
	return 1
}

func sorted(L *lua.LState) int {
	lessThan := L.NewFunction(func(L *lua.LState) int {
		a := L.CheckInt(1)
		b := L.CheckInt(2)

		L.Push(lua.LBool(a < b))
		return 1
	})

	list := checkList(L, 1)
	fn := L.OptFunction(2, lessThan)

	slices.SortFunc(list, func(i, j lua.LValue) bool {
		L.Push(fn)
		L.Push(i)
		L.Push(j)
		L.Call(2, 1)
		return L.ToBool(-1)
	})

	sorted, _ := luautil.ToLValue(L, list)
	L.Push(sorted)
	return 1
}

func zip2(L *lua.LState) int {
	list1 := checkList(L, 1)
	list2 := checkList(L, 2)

	zipped, _ := luautil.ToLValue(L, lo.Zip2(list1, list2))
	L.Push(zipped)
	return 1
}

func zip3(L *lua.LState) int {
	list1 := checkList(L, 1)
	list2 := checkList(L, 2)
	list3 := checkList(L, 3)

	zipped, _ := luautil.ToLValue(L, lo.Zip3(list1, list2, list3))
	L.Push(zipped)
	return 1
}

func map_(L *lua.LState) int {
	list := checkList(L, 1)
	fn := L.CheckFunction(2)

	mapped, _ := luautil.ToLValue(L, lo.Map(list, func(v lua.LValue, i int) lua.LValue {
		L.Push(fn)
		L.Push(v)
		L.Push(lua.LNumber(i))
		L.Call(2, 1)
		return L.Get(-1)
	}))
	L.Push(mapped)
	return 1
}

func forEach(L *lua.LState) int {
	list := checkList(L, 1)
	fn := L.CheckFunction(2)

	lo.ForEach(list, func(v lua.LValue, i int) {
		L.Push(fn)
		L.Push(v)
		L.Push(lua.LNumber(i))
		L.Call(2, 0)
	})

	return 0
}

func containsBy(L *lua.LState) int {
	list := checkList(L, 1)
	fn := L.CheckFunction(2)

	contains := lo.ContainsBy(list, func(v lua.LValue) bool {
		L.Push(fn)
		L.Push(v)
		L.Call(1, 1)
		return L.ToBool(-1)
	})

	L.Push(lua.LBool(contains))
	return 1
}

func reduce(L *lua.LState) int {
	list := checkList(L, 1)
	fn := L.CheckFunction(2)
	initial := L.CheckAny(3)

	r := lo.Reduce[lua.LValue, lua.LValue](list, func(acc lua.LValue, v lua.LValue, i int) lua.LValue {
		L.Push(fn)
		L.Push(acc)
		L.Push(v)
		L.Push(lua.LNumber(i))
		L.Call(3, 1)
		return L.Get(-1)
	}, initial)

	L.Push(r)
	return 1
}

func reduceRight(L *lua.LState) int {
	list := checkList(L, 1)
	fn := L.CheckFunction(2)
	initial := L.CheckAny(3)

	r := lo.ReduceRight[lua.LValue, lua.LValue](list, func(acc lua.LValue, v lua.LValue, i int) lua.LValue {
		L.Push(fn)
		L.Push(acc)
		L.Push(v)
		L.Push(lua.LNumber(i))
		L.Call(3, 1)
		return L.Get(-1)
	}, initial)

	L.Push(r)
	return 1
}

func reject(L *lua.LState) int {
	list := checkList(L, 1)
	fn := L.CheckFunction(2)

	rejected, _ := luautil.ToLValue(L, lo.Reject(list, func(v lua.LValue, i int) bool {
		L.Push(fn)
		L.Push(v)
		L.Push(lua.LNumber(i))
		L.Call(2, 1)
		return L.ToBool(-1)
	}))
	L.Push(rejected)
	return 1
}

func sample(L *lua.LState) int {
	list := checkList(L, 1)

	sampled, _ := luautil.ToLValue(L, lo.Sample(list))
	L.Push(sampled)
	return 1
}

func samples(L *lua.LState) int {
	list := checkList(L, 1)
	n := L.CheckInt(2)

	sampled, _ := luautil.ToLValue(L, lo.Samples(list, n))
	L.Push(sampled)
	return 1
}

func shuffle(L *lua.LState) int {
	list := checkList(L, 1)

	shuffled, _ := luautil.ToLValue(L, lo.Shuffle(list))
	L.Push(shuffled)
	return 1
}

func someBy(L *lua.LState) int {
	list := checkList(L, 1)
	fn := L.CheckFunction(2)

	L.Push(lua.LBool(lo.SomeBy(list, func(v lua.LValue) bool {
		L.Push(fn)
		L.Push(v)
		L.Call(2, 1)
		return L.ToBool(-1)
	})))
	return 1
}

func everyBy(L *lua.LState) int {
	list := checkList(L, 1)
	fn := L.CheckFunction(2)

	L.Push(lua.LBool(lo.EveryBy(list, func(v lua.LValue) bool {
		L.Push(fn)
		L.Push(v)
		L.Call(2, 1)
		return L.ToBool(-1)
	})))
	return 1
}

func id(L *lua.LState) int {
	value := L.CheckAny(1)
	L.Push(value)
	return 1
}
