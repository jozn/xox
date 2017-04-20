package com.mardomsara.social.json;

public class J {

{{range $key,$model := . }}
{{- with $model }}
	public static class {{.Name}} {
		{{- range .Fields }}
		public {{ retype .Type }} {{ .Col.ColumnName }};
		{{- end }}
	}
{{end -}}
{{end}}
}