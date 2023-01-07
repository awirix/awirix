package js

import (
	"github.com/awirix/awirix/luadoc"
)

func Lib() *luadoc.Lib {
	classVMValue := &luadoc.Class{
		Name:        vmValueTypeName,
		Description: "A value returned from a JavaScript VM.",
		Methods: []*luadoc.Method{
			{
				Name:        "export",
				Description: "Exports the value to a Lua value.",
				Value:       vmValueExport,
				Returns: []*luadoc.Param{
					{
						Name:        "value",
						Description: "The exported value.",
						Type:        luadoc.Any,
					},
					{
						Name:        "error",
						Description: "The error message, if any.",
						Type:        luadoc.String,
						Optional:    true,
					},
				},
			},
			{
				Name:        "to_string",
				Description: "Converts the value to a string.",
				Value:       vmValueString,
				Returns: []*luadoc.Param{
					{
						Name:        "string",
						Description: "The string representation of the value.",
						Type:        luadoc.String,
					},
				},
			},
			{
				Name:        "is_null",
				Description: "Returns whether the value is null.",
				Value:       vmValueIsNull,
				Returns: []*luadoc.Param{
					{
						Name:        "ok",
						Description: "Whether the value is null.",
						Type:        luadoc.Boolean,
					},
				},
			},
			{
				Name:        "is_undefined",
				Description: "Returns whether the value is undefined.",
				Value:       vmValueIsUndefined,
				Returns: []*luadoc.Param{
					{
						Name:        "ok",
						Description: "Whether the value is undefined.",
						Type:        luadoc.Boolean,
					},
				},
			},
			{
				Name:        "is_boolean",
				Description: "Returns whether the value is a boolean.",
				Value:       vmValueIsBoolean,
				Returns: []*luadoc.Param{
					{
						Name:        "ok",
						Description: "Whether the value is a boolean.",
						Type:        luadoc.Boolean,
					},
				},
			},
			{
				Name:        "is_number",
				Description: "Returns whether the value is a number.",
				Value:       vmValueIsNumber,
				Returns: []*luadoc.Param{
					{
						Name:        "ok",
						Description: "Whether the value is a number.",
						Type:        luadoc.Boolean,
					},
				},
			},
			{
				Name:        "is_string",
				Description: "Returns whether the value is a string.",
				Value:       vmValueIsString,
				Returns: []*luadoc.Param{
					{
						Name:        "ok",
						Description: "Whether the value is a string.",
						Type:        luadoc.Boolean,
					},
				},
			},
			{
				Name:        "is_object",
				Description: "Returns whether the value is an object.",
				Value:       vmValueIsObject,
				Returns: []*luadoc.Param{
					{
						Name:        "ok",
						Description: "Whether the value is an object.",
						Type:        luadoc.Boolean,
					},
				},
			},
			{
				Name:        "is_nan",
				Description: "Returns whether the value is NaN.",
				Value:       vmValueIsNaN,
				Returns: []*luadoc.Param{
					{
						Name:        "ok",
						Description: "Whether the value is NaN.",
						Type:        luadoc.Boolean,
					},
				},
			},
			{
				Name:        "is_function",
				Description: "Returns whether the value is a function.",
				Value:       vmValueIsFunction,
				Returns: []*luadoc.Param{
					{
						Name:        "ok",
						Description: "Whether the value is a function.",
						Type:        luadoc.Boolean,
					},
				},
			},
		},
	}

	classVM := &luadoc.Class{
		Name:        "js_vm",
		Description: "A JavaScript virtual machine. This is used to execute JavaScript code.",
		Methods: []*luadoc.Method{
			{
				Name:        "run",
				Description: "Runs the given JavaScript code.",
				Value:       vmRun,
				Params: []*luadoc.Param{
					{
						Name:        "code",
						Description: "The JavaScript code to run.",
						Type:        luadoc.String,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "value",
						Description: "The value returned by the code.",
						Type:        classVMValue.Name,
					},
					{
						Name:        "error",
						Description: "The error that occurred, if any.",
						Type:        luadoc.String,
						Optional:    true,
					},
				},
			},
			{
				Name:        "get",
				Description: "Gets the value of the given property on the global object.",
				Value:       vmGet,
				Params: []*luadoc.Param{
					{
						Name:        "name",
						Description: "The name of the property.",
						Type:        luadoc.String,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "value",
						Description: "The value of the property.",
						Type:        classVMValue.Name,
					},
					{
						Name:        "error",
						Description: "The error that occurred, if any.",
						Type:        luadoc.String,
						Optional:    true,
					},
				},
			},
			{
				Name:        "set",
				Description: "Sets the value of the given property on the global object. It will convert the given Lua value to a JavaScript value.",
				Value:       vmSet,
				Params: []*luadoc.Param{
					{
						Name:        "name",
						Description: "The name of the property.",
						Type:        luadoc.String,
					},
					{
						Name:        "value",
						Description: "The value to set.",
						Type:        luadoc.Any,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "error",
						Description: "The error that occurred, if any.",
						Type:        luadoc.String,
						Optional:    true,
					},
				},
			},
		},
	}

	return &luadoc.Lib{
		Name:        "js",
		Description: "JavaScript execution.",
		Funcs: []*luadoc.Func{
			{
				Name:        "new_vm",
				Description: "Creates a new JavaScript virtual machine.",
				Value:       newVM,
				Returns: []*luadoc.Param{
					{
						Name:        "vm",
						Description: "The new JavaScript virtual machine.",
						Type:        classVM.Name,
					},
				},
			},
		},
		Classes: []*luadoc.Class{
			classVM,
			classVMValue,
		},
	}
}
