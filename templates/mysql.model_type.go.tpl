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
	{{ .Col.ColumnName }} {{ retype .Type }} {{ ms_col_comment_json .Comment }} {{ ms_col_comment_raw .Comment }}      {{/* `json:"{{ .Col.ColumnName }}"` // {{ .Col.ColumnName }} */}}
{{- end }}
{{- if .PrimaryKey }}
	{{/* // xox fields */}}
	_exists, _deleted bool
{{ end }}
}
/*
:= &{{ .Name }} {
{{- range .Fields }}
	{{ .Col.ColumnName }}: {{datatype_to_defualt_go_type .Type }},
{{- end }}
*/


