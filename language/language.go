package language

import "fmt"

// Language is an ISO 639-1 language with code, name and native name.
type Language struct {
	Code       string
	Name       string
	NativeName string
}

func (l Language) MarshalText() ([]byte, error) {
	return []byte(l.Code), nil
}

func (l *Language) UnmarshalText(text []byte) error {
	lang, ok := Languages[string(text)]
	if !ok {
		return fmt.Errorf("invalid language code: %s", text)
	}

	*l = *lang
	return nil
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
