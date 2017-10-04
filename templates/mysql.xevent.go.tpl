{{- $short := (shortname .Name "err" "res" "sqlstr" "db" "XOLog") -}}
{{- $table := (schema .Schema .Table.TableName) -}}
{{- $typ := .Name }}
{{- $_ := "" }}

{{/* - * (Manually copy this to other location) */}}
//{{ .Name }} Events
{{if (eq .PrimaryKey.Type "int") }}
func  On{{ .Name }}_AfterInsert{{$_}} (row *{{ .Name }}) {
	RowCache.Set("{{ .Name }}:"+strconv.Itoa(row.{{.PrimaryKey.Name}}), row,time.Hour* 0)
}

func  On{{ .Name }}_AfterUpdate{{$_}} (row *{{ .Name }}) {
	RowCache.Set("{{ .Name }}:"+strconv.Itoa(row.{{.PrimaryKey.Name}}), row,time.Hour* 0)
}

func  On{{ .Name }}_AfterDelete{{$_}} (row *{{ .Name }}) {
	RowCache.Delete("{{ .Name }}:"+strconv.Itoa(row.{{.PrimaryKey.Name}}))
}

func  On{{ .Name }}_LoadOne{{$_}} (row *{{ .Name }}) {
	RowCache.Set("{{ .Name }}:"+strconv.Itoa(row.{{.PrimaryKey.Name}}), row,time.Hour* 0)
}

func  On{{ .Name }}_LoadMany{{$_}} (rows []*{{ .Name }}) {
	for _, row:= range rows {
		RowCache.Set("{{ .Name }}:"+strconv.Itoa(row.{{.PrimaryKey.Name}}), row,time.Hour* 0)
	}
}
{{else if ( eq .PrimaryKey.Type "string" ) }}
func  On{{ .Name }}_AfterInsert{{$_}} (row *{{ .Name }}) {
	RowCache.Set("{{ .Name }}:"+row.{{.PrimaryKey.Name}}, row,time.Hour* 0)
}

func  On{{ .Name }}_AfterUpdate{{$_}} (row *{{ .Name }}) {
	RowCache.Set("{{ .Name }}:"+row.{{.PrimaryKey.Name}}, row,time.Hour* 0)
}

func  On{{ .Name }}_AfterDelete{{$_}} (row *{{ .Name }}) {
	RowCache.Delete("{{ .Name }}:"+row.{{.PrimaryKey.Name}})
}

func  On{{ .Name }}_LoadOne{{$_}} (row *{{ .Name }}) {
	RowCache.Set("{{ .Name }}:"+row.{{.PrimaryKey.Name}}, row,time.Hour* 0)
}

func  On{{ .Name }}_LoadMany{{$_}} (rows []*{{ .Name }}) {
	for _, row:= range rows {
		RowCache.Set("{{ .Name }}:"+row.{{.PrimaryKey.Name}}, row,time.Hour* 0)
	}
}
{{end}}


