package luadoc

import (
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
	}).
	Parse(`--- @meta
{{ define "lib" }}
{{ doc .Description }}
{{ doc "@class" }} {{ .Name }}
{{- range $var := .Vars }}
{{ doc "@field" }} public {{ $var.Name }} {{ $var.Value.Type.String }} {{ $var.Description }}
{{- end }}
{{- range $lib := .Libs }}
{{ doc "@field" }} public {{ $lib.Name }} {{ $lib.Name }} {{ $lib.Description }}
{{- end }}
local {{ .Name }} = {}

{{ range $func := .Funcs }}
{{ doc $func.Description }}
{{- range $param := $func.Params }}
{{ doc "@param" }} {{ $param.Name }} {{ $param.Type }}{{ if $param.Opt }}?{{ end }} {{ $param.Description }}
{{- end }}
{{- range $return := $func.Returns }}
{{ doc "@return" }} {{ $return.Type }}{{ if $return.Opt }}?{{ end }} {{ $return.Name }} {{ $return.Description }}
{{- end }}
function {{ $.Name }}.{{ $func.Name }}({{ join (params $func.Params) ", " }}) end
{{ end }}

{{ range $class := .Classes }}
{{ doc $class.Description }}
{{ doc "@class" }} {{ $class.Name }}
local {{ $class.Name }} = {}

{{ range $method := $class.Methods }}
{{ doc $method.Description }}
{{- range $param := $method.Params }}
{{ doc "@param" }} {{ $param.Name }} {{ $param.Type }}{{ if $param.Opt }}?{{ end }} {{ $param.Description }}
{{- end }}
{{- range $return := $method.Returns }}
{{ doc "@return" }} {{ $return.Type }}{{ if $return.Opt }}?{{ end }} {{ $return.Name }} {{ $return.Description }}
{{- end }}
function {{ $class.Name }}:{{ $method.Name }}({{ join (params $method.Params) ", " }}) end
{{ end }}

{{ end }}

{{ range $lib := .Libs }}
{{ template "lib" $lib }}
{{ end }}

{{- end }}

{{ template "lib" . }}

return {{ .Name }}
`))
