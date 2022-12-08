package passport

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type ConfigSection struct {
	Display string `json:"display,omitempty" jsonschema:"description=The display name of the config section"`
	About   string `json:"about,omitempty" jsonschema:"description=The description of the config section"`
	Default any    `json:"default" jsonschema:"required,description=The default value of the config section"`
	Values  []any  `json:"values,omitempty" jsonschema:"description=The possible values of the config section. If specified the config section will be an enum"`
	value   any
}

func (c *ConfigSection) UnmarshalJSON(text []byte) error {
	type raw struct {
		Display string `json:"display,omitempty"`
		About   string `json:"about,omitempty"`
		Default any    `json:"default"`
		Values  []any  `json:"values,omitempty"`
		value   any
	}

	var r raw

	if err := json.Unmarshal(text, &r); err != nil {
		return err
	}

	if r.Default == nil {
		return fmt.Errorf("default is required")
	}

	if r.Values != nil {
		// handle it as enum
		if len(r.Values) == 0 {
			return fmt.Errorf("values must not be empty")
		}

		var contains bool

		// check if all values are the same type as default
		for _, v := range r.Values {
			if reflect.TypeOf(v) != reflect.TypeOf(r.Default) {
				return fmt.Errorf("values must be of the same type as default")
			}

			if !contains && reflect.DeepEqual(v, r.Default) {
				contains = true
			}
		}

		if !contains {
			return fmt.Errorf("default value must one of the specified values")
		}
	}

	*c = ConfigSection(r)

	return nil
}

func (c *ConfigSection) SetValue(value any) {
	c.value = value
}

func (c *ConfigSection) Value() any {
	if c.value != nil {
		return c.value
	}

	return c.Default
}
