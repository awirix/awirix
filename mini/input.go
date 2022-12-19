package mini

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/samber/lo"
	"github.com/vivi-app/vivi/extensions/passport"
	"strconv"
)

func selectOne[T any](message string, items []T, render func(T) string) (T, error) {
	theSpinner.Stop()
	defer theSpinner.Start()

	var stringMap = make(map[string]T)

	var name string
	err := survey.AskOne(&survey.Select{
		Message: message,
		Options: lo.Map(items, func(item T, _ int) string {
			s := render(item)
			stringMap[s] = item
			return s
		}),
		PageSize: 10,
		VimMode:  true,
	}, &name)

	if err != nil {
		var t T
		return t, err
	}

	return stringMap[name], nil
}

func getString(message string, dflt string) (input string, err error) {
	theSpinner.Stop()
	defer theSpinner.Start()

	err = survey.AskOne(&survey.Input{
		Message: message,
		Default: dflt,
	}, &input)

	return
}

func getBool(message string, dflt bool) (input bool, err error) {
	theSpinner.Stop()
	defer theSpinner.Start()

	err = survey.AskOne(&survey.Confirm{
		Message: message,
		Default: dflt,
	}, &input)

	return
}

func getFloat(message string, dflt float64) (input float64, err error) {
	theSpinner.Stop()
	defer theSpinner.Start()

	var numberString string
	err = survey.AskOne(&survey.Input{
		Message: message,
		Default: fmt.Sprint(dflt),
	}, &numberString)

	if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(numberString, 64)
}

func getConfigValue(section *passport.ConfigSection) (value any, err error) {
	if len(section.Values) != 0 {
		return selectOne[any](section.Display, section.Values, func(a any) string { return fmt.Sprint(a) })
	}

	switch section.Value().(type) {
	case string:
		return getString(section.Display, section.Default.(string))
	case bool:
		return getBool(section.Display, section.Default.(bool))
	case float64:
		return getFloat(section.Display, section.Default.(float64))
	default:
		return nil, fmt.Errorf("unsupported type %T", section.Value())
	}
}
