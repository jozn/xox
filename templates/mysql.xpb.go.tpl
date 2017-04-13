{{- $short := (shortname .Name "err" "res" "sqlstr" "db" "XOLog") -}}
{{- $table := (schema .Schema .Table.TableName) -}}
{{- $typ := .Name}}
{{- if .Comment -}}
// {{ .Comment }}
{{- else -}}
// {{ .Name }} represents a row from '{{ $table }}'.
{{- end }}

// Manualy copy this to project
message {{ .Name }}_Row {
{{- range $i, $e := .Fields }}
	required {{ retype $e.Type }}  {{ $e.Col.ColumnName }} {{ $i }} ;// {{ .Col.ColumnName }} -
{{- end }}

}
