package mini

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/samber/lo"
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

func getString(message string) (input string, err error) {
	theSpinner.Stop()
	defer theSpinner.Start()

	err = survey.AskOne(&survey.Input{
		Message: message,
	}, &input)

	return
}
