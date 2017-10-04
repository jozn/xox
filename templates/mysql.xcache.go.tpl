{{- $short := (shortname .Name "err" "res" "sqlstr" "db" "XOLog") -}}
{{- $table := (schema .Schema .Table.TableName) -}}
{{- $typ := .Name }}
{{- $_ := "" }}
{{- $id := (.PrimaryKey.Name) }}

{{if (eq .PrimaryKey.Type "int") }}
{{/* //{{ .Name }} Events - * (Manually copy this to other location) */}}
func (c _StoreImpl) Get{{ .Name }}By{{$id}}{{$_}} ({{$id}} int) (*{{ .Name }},bool){
	o ,ok :=RowCache.Get("{{ .Name }}:"+strconv.Itoa({{$id}}))
	if ok {
		if obj, ok := o.(*{{ .Name }});ok{
			return obj, true
		}
	}
	obj2 ,err := {{ .Name }}By{{.PrimaryKey.Name}}(base.DB, {{$id}})
	if err == nil {
		return obj2, true
	}
	XOLogErr(err)
	return nil, false
}

func (c _StoreImpl) PreLoad{{ .Name }}By{{$id}}s{{$_}} (ids []int) {
	not_cached := make([]int,0,len(ids))

	for _,id := range ids {
		_ ,ok :=RowCache.Get("{{ .Name }}:"+strconv.Itoa(id))
		if !ok {
			not_cached = append(not_cached,id)
		}
	}

	if len(not_cached) > 0 {
		New{{ .Name }}_Selector().{{$id}}_In(not_cached).GetRows(base.DB)
	}
}
{{else if ( eq .PrimaryKey.Type "string" ) }}
func (c _StoreImpl) Get{{ .Name }}By{{$id}}{{$_}} ({{$id}} string) (*{{ .Name }},bool){
	o ,ok :=RowCache.Get("{{ .Name }}:"+{{$id}})
	if ok {
		if obj, ok := o.(*{{ .Name }});ok{
			return obj, true
		}
	}
	obj2 ,err := {{ .Name }}By{{.PrimaryKey.Name}}(base.DB, {{$id}})
	if err == nil {
		return obj2, true
	}
	XOLogErr(err)
	return nil, false
}

func (c _StoreImpl) PreLoad{{ .Name }}By{{$id}}s{{$_}} (ids []string) {
	not_cached := make([]string,0,len(ids))

	for _,id := range ids {
		_ ,ok :=RowCache.Get("{{ .Name }}:"+id)
		if !ok {
			not_cached = append(not_cached,id)
		}
	}

	if len(not_cached) > 0 {
		New{{ .Name }}_Selector().{{$id}}_In(not_cached).GetRows(base.DB)
	}
}
{{end}}
// yes 222 {{.PrimaryKey.Type }}




