package formsgen

import (
	"html/template"
)

func init() {
	fieldsTemplate = template.Must(template.New("fields").Parse(templates))
}

var fieldsTemplate *template.Template

var templates = `
{{define "text"}}
    <label for="{{ .Field }}">{{ .Name}}</label>
    <input type="text" name="{{ .Field }}"{{ if .Default }} value="{{ .Value }}"{{ end }}{{ if .Required }} required{{ end }}>
{{end}}

{{define "password"}}
    <label for="{{ .Field }}">{{ .Name}}</label>
    <input type="password" name="{{ .Field }}"{{ if .Default }} value="{{ .Value }}"{{ end }}{{ if .Required }} required{{ end }}>
{{end}}

{{define "hidden"}}
    <input type="hidden" name="{{ .Field }}"{{ if .Default }} value="{{ .Value }}"{{ end }}>
{{end}}

{{define "button"}}
    <input type="button" name="{{ .Field }}" value="{{ .Name }}">
{{end}}

{{define "checkbox"}}
    <label>
        <input type="checkbox" name="{{ .Field }}"{{ if and .Default .Value.Bool }} checked{{ end }}> {{ .Name }}
    </label>
{{end}}

{{define "radio"}}
    <label>
        <input type="radio" name="{{ .Field }}" value="{{ .RadioValue.Value }}"{{ if .RadioValue.Checked }} checked{{ end }}> {{ .Name }}
    </label>
{{end}}

{{define "select"}}
    <label for="{{ .Field }}">{{ .Name}}</label>
    <select name="{{ .Field }}">{{ range .SelectValues }}
        <option value="{{ .Value }}" {{if .Selected}}selected{{end}}>{{ .Label }}</option>{{ end }}
    </select>
{{end}}
`
