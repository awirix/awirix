package language

import (
	"encoding/json"
	"fmt"
)

// Language is an ISO 639-1 language with code, name and native name.
type Language struct {
	Code       string
	Name       string
	NativeName string
}

func (l *Language) UnmarshalJSON(data []byte) error {
	var s string

	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	lang, ok := FromCode(s)
	if !ok {
		return fmt.Errorf("language: invalid ISO 639-1 code: %s", s)
	}

	*l = *lang
	return nil
}

func (l *Language) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.Code)
}

func FromCode(code string) (*Language, bool) {
	lang, ok := Languages[code]
	return lang, ok
}

func FromName(name string) (*Language, bool) {
	for _, lang := range Languages {
		if lang.Name == name {
			return lang, true
		}
	}

	return nil, false
}

func FromNativeName(name string) (*Language, bool) {
	for _, lang := range Languages {
		if lang.NativeName == name {
			return lang, true
		}
	}

	return nil, false
}
