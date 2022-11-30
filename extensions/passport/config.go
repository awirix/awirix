package passport

import (
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"reflect"
)

type ConfigSection struct {
	Display string `toml:"display,omitempty"`
	About   string `toml:"about,omitempty"`
	Default any    `toml:"default"`
	Values  []any  `toml:"values,omitempty"`
	value   any
}

func (c *ConfigSection) UnmarshalText(text []byte) error {
	type raw struct {
		Display string `toml:"display"`
		About   string `toml:"about"`
		Default any    `toml:"default"`
		Values  []any  `toml:"values"`
		value   any
	}

	var r raw

	if err := toml.Unmarshal(text, &r); err != nil {
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

		// check if all values are the same type as default
		for _, v := range r.Values {
			if reflect.TypeOf(v) != reflect.TypeOf(r.Default) {
				return fmt.Errorf("values must be of the same type as default")
			}
		}

		//if !lo.Contains(r.Values, r.Default) {
		//	return fmt.Errorf("default value must one of the specified values")
		//}
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
