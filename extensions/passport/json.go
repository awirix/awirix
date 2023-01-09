package passport

import (
	"encoding/json"
	"fmt"
	"github.com/enescakir/emoji"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"strings"
)

func errPassport(err error) error {
	return errors.Wrap(err, "passport")
}

func (p *Passport) UnmarshalJSON(data []byte) error {
	type Alias Passport
	aux := &struct {
		*Alias
		AuxIcon string `json:"icon"`
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

	if aux.AuxIcon == "" {
		return nil
	}

	// if AuxIcon is emoji already and not an alias
	if lo.Contains(lo.Values(emoji.Map()), aux.AuxIcon) {
		p.Icon = emoji.Emoji(aux.AuxIcon)
		return nil
	}

	if !emoji.Exist(aux.AuxIcon) {
		return errPassport(fmt.Errorf("invalid emoji alias: %q", aux.AuxIcon))
	}

	p.Icon = emoji.Emoji(emoji.Map()[aux.AuxIcon])

	return nil
}
