package luadoc

import (
	"fmt"
	"strings"
	"text/template"
)

var templateLuaDocLib = template.Must(template.
	New("luadoc").
	Funcs(map[string]any{
		"doc": func(s string) string {
			var b strings.Builder

			for _, line := range strings.Split(s, "\n") {
				b.WriteString("--- ")
				b.WriteString(line)
				b.WriteString("\n")
			}

			return strings.TrimSpace(b.String())
		},
		"join": func(s []string, delim string) string {
			return strings.Join(s, delim)
		},
		"params": func(s []*Param) []string {
			var strings = make([]string, len(s))
			for i, v := range s {
				strings[i] = v.String()
			}

			return strings
		},
		"vartype": func(v *Var) string {
			if v.Type == "" {
				return v.Value.Type().String()
			}

			return v.Type
		},
		"defineTypes": func(vars []*Var) string {
			var custom = make(map[string]*Var, 0)
			for _, v := range vars {
				if v.Type != "" {
					custom[v.Type] = v
				}
			}

			if len(custom) == 0 {
				return ""
			}

			var b strings.Builder
			for t, v := range custom {
				b.WriteString(fmt.Sprintf("--- @alias %s %s\n", t, v.Value.Type().String()))
			}

			return strings.TrimSpace(b.String())
		},
	}).
	Parse(`--- @meta
{{ define "lib" }}

{{ defineTypes .Vars }}

{{ doc .Description }}
{{ doc "@class" }} {{ .Name }}
{{- range $var := .Vars }}
{{ doc "@field" }} public {{ $var.Name }} {{ vartype $var }} {{ $var.Description }}
{{- end }}
{{- range $lib := .Libs }}
{{ doc "@field" }} public {{ $lib.Name }} {{ $lib.Name }} {{ $lib.Description }}
{{- end }}
local {{ .Name }} = {}

{{ range $class := .Classes }}
{{ doc $class.Description }}
{{ doc "@class" }} {{ $class.Name }}
local {{ $class.Name }} = {}

{{ range $method := $class.Methods }}
{{ doc $method.Description }}
{{- range $param := $method.Params }}
{{ doc "@param" }} {{ $param.Name }} {{ $param.Type }}{{ if $param.Optional }}?{{ end }} {{ $param.Description }}
{{- end }}
{{- range $return := $method.Returns }}
{{ doc "@return" }} {{ $return.Type }}{{ if $return.Optional }}?{{ end }} {{ $return.Name }} {{ $return.Description }}
{{- end }}
function {{ $class.Name }}:{{ $method.Name }}({{ join (params $method.Params) ", " }}) end
{{ end }}
{{ end }}

{{ range $func := .Funcs }}
{{ doc $func.Description }}
{{- range $param := $func.Params }}
{{ doc "@param" }} {{ $param.Name }} {{ $param.Type }}{{ if $param.Optional }}?{{ end }} {{ $param.Description }}
{{- end }}
{{- range $return := $func.Returns }}
{{ doc "@return" }} {{ $return.Type }}{{ if $return.Optional }}?{{ end }} {{ $return.Name }} {{ $return.Description }}
{{- end }}
function {{ $.Name }}.{{ $func.Name }}({{ join (params $func.Params) ", " }}) end
{{ end }}


{{ range $lib := .Libs }}
{{ template "lib" $lib }}
{{ end }}

{{- end }}

{{ template "lib" . }}

return {{ .Name }}
`))
