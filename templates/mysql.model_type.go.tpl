{{- $short := (shortname .Name "err" "res" "sqlstr" "db" "XOLog") -}}
{{- $table := (schema .Schema .Table.TableName) -}}
{{- $typ := .Name}}
{{- if .Comment -}}
// {{ .Comment }}
{{- else -}}
// {{ .Name }} represents a row from '{{ $table }}'.
{{- end }}

// Manualy copy this to project
type {{ .Name }} struct {
{{- range .Fields }}
	{{ .Col.ColumnName }} {{ retype .Type }}// `json:"{{ .Col.ColumnName }}"` // {{ .Col.ColumnName }} -
{{- end }}
{{- if .PrimaryKey }}

	// xo fields
	_exists, _deleted bool
{{ end }}
}
