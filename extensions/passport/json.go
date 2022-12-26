package passport

import (
	"encoding/json"
	"fmt"
	"github.com/vivi-app/vivi/language"
	"golang.org/x/exp/slices"
	"strings"
)

func (p *Passport) UnmarshalJSON(data []byte) error {
	type Alias Passport
	aux := &struct {
		*Alias
	}{}

	aux.Alias = (*Alias)(p)

	err := json.Unmarshal(data, aux)
	if err != nil {
		return err
	}

	aux.Name = strings.TrimSpace(aux.Name)

	if strings.Contains(aux.ID, " ") {
		return fmt.Errorf("passport: id cannot contain spaces")
	}

	if !slices.Contains(language.Codes, aux.LanguageRaw) {
		return fmt.Errorf("passport: invalid ISO 639-1 language code: %s", aux.LanguageRaw)
	}

	if len(aux.Tags) > 5 {
		return fmt.Errorf("passport: maximum of 5 tags")
	}

	if len(aux.About) > 100 {
		return fmt.Errorf("passport: maximum of 100 characters for about")
	}

	return nil
}
