package passport

import (
	"encoding/json"
	"github.com/pkg/errors"
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
		return errPassport(errors.New("id cannot contain spaces"))
	}

	if len(aux.Tags) > 5 {
		return errPassport(errors.Errorf("maximum of 5 tags, got %d", len(aux.Tags)))
	}

	if len(aux.About) > 100 {
		return errPassport(errors.Errorf("maximum of 100 characters for about, got %d", len(aux.About)))
	}

	return nil
}
