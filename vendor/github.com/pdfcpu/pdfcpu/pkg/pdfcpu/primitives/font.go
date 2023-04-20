/*
	Copyright 2021 The pdfcpu Authors.

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package primitives

import (
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/font"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/color"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
	"github.com/pkg/errors"
)

type FormFont struct {
	pdf    *PDF
	Name   string
	Lang   string // ISO-639
	Script string // ISO-15924
	Size   int
	Color  string `json:"col"`
	col    *color.SimpleColor
}

// ISO-639 country codes
// See https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes
var ISO639Codes = []string{"ab", "aa", "af", "ak", "sq", "am", "ar", "an", "hy", "as", "av", "ae", "ay", "az", "bm", "ba", "eu", "be", "bn", "bi", "bs", "br", "bg",
	"my", "ca", "ch", "ce", "ny", "zh", "cu", "cv", "kw", "co", "cr", "hr", "cs", "da", "dv", "nl", "dz", "en", "eo", "et", "ee", "fo", "fj", "fi", "fr", "fy", "ff",
	"gd", "gl", "lg", "ka", "de", "el", "kl", "gn", "gu", "ht", "ha", "he", "hz", "hi", "ho", "hu", "is", "io", "ig", "id", "ia", "ie", "iu", "ik", "ga", "it", "ja",
	"jv", "kn", "kr", "ks", "kk", "km", "ki", "rw", "ky", "kv", "kg", "ko", "kj", "ku", "lo", "la", "lv", "li", "ln", "lt", "lu", "lb", "mk", "mg", "ms", "ml", "mt",
	"gv", "mi", "mr", "mh", "mn", "na", "nv", "nd", "nr", "ng", "ne", "no", "nb", "nn", "ii", "oc", "oj", "or", "om", "os", "pi", "ps", "fa", "pl", "pt", "pa", "qu",
	"ro", "rm", "rn", "ru", "se", "sm", "sg", "sa", "sc", "sr", "sn", "sd", "si", "sk", "sl", "so", "st", "es", "su", "sw", "ss", "sv", "tl", "ty", "tg", "ta", "tt",
	"te", "th", "bo", "ti", "to", "ts", "tn", "tr", "tk", "tw", "ug", "uk", "ur", "uz", "ve", "vi", "vo", "wa", "cy", "wo", "xh", "yi", "yo", "za", "zu"}

func (f *FormFont) validateISO639() error {
	if !types.MemberOf(f.Lang, ISO639Codes) {
		return errors.Errorf("pdfcpu: invalid ISO-639 code: %s", f.Lang)
	}
	return nil
}

func (f *FormFont) validateScriptSupport() error {
	fd, ok := font.UserFontMetrics[f.Name]
	if !ok {
		return errors.Errorf("pdfcpu: userfont %s not available", f.Name)
	}
	ok, err := fd.SupportsScript(f.Script)
	if err != nil {
		return err
	}
	if !ok {
		return errors.Errorf("pdfcpu: userfont (%s) does not support script: %s", f.Name, f.Script)
	}
	return nil
}

func (f *FormFont) validate() error {
	if f.Name == "$" {
		return errors.New("pdfcpu: invalid font reference $")
	}

	if f.Name != "" && f.Name[0] != '$' {
		if !font.SupportedFont(f.Name) {
			return errors.Errorf("pdfcpu: font %s is unsupported, please refer to \"pdfcpu fonts list\".\n", f.Name)
		}
		if font.IsUserFont(f.Name) {
			if f.Lang != "" {
				f.Lang = strings.ToLower(f.Lang)
				if err := f.validateISO639(); err != nil {
					return err
				}
			}
			if f.Script != "" {
				f.Script = strings.ToUpper(f.Script)
				if err := f.validateScriptSupport(); err != nil {
					return err
				}
			}
		}
		if f.Size <= 0 {
			return errors.Errorf("pdfcpu: invalid font size: %d", f.Size)
		}
	}

	if f.Color != "" {
		sc, err := f.pdf.parseColor(f.Color)
		if err != nil {
			return err
		}
		f.col = sc
	}

	return nil
}

func (f *FormFont) mergeIn(f0 *FormFont) {
	if f.Name == "" {
		f.Name = f0.Name
	}
	if f.Size == 0 {
		f.Size = f0.Size
	}
	if f.col == nil {
		f.col = f0.col
	}
	if f.Lang == "" {
		f.Lang = f0.Lang
	}
	if f.Script == "" {
		f.Script = f0.Script
	}
}

func (f *FormFont) SetCol(c color.SimpleColor) {
	f.col = &c
}

func (f FormFont) RTL() bool {
	return types.MemberOf(f.Script, []string{"Arab", "Hebr"}) || types.MemberOf(f.Lang, []string{"ar", "fa", "he"})
}
