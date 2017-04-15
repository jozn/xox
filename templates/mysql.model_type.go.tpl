{{- $short := (shortname .Name "err" "res" "sqlstr" "db" "XOLog") -}}
{{- $table := (schema .Schema .Table.TableName) -}}
{{- $typ := .Name}}
{{- if .Comment -}}
// {{ .Comment }}
{{- else -}}
// {{ .Name }} '{{ $table }}'.
{{- end }}
type {{ .Name }} struct {
{{- range .Fields }}
	{{ .Col.ColumnName }} {{ retype .Type }} {{/* `json:"{{ .Col.ColumnName }}"` // {{ .Col.ColumnName }} */}}
{{- end }}
{{- if .PrimaryKey }}
	{{/* // xox fields */}}
	_exists, _deleted bool
{{ end }}
}
