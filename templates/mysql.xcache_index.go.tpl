{{- $short := (shortname .Type.Name "err" "sqlstr" "db" "q" "res" "XOLog" .Fields) -}}
{{- $table := (schema .Schema .Type.Table.TableName) -}}
{{- $model := .Type.Name -}}

{{if (and (eq (len .Fields) 1) (not .Index.IsPrimary) ) }}

{{- $col := (index .Fields 0) -}}//field
{{- $colType := (retype $col.Type ) -}}//field
{{- $indexName := (printf "%s_%s" $model  .Index.IndexName) -}}//field
{{$param := (printf "%s" $col.Name) }}
///// Generated from index '{{ .Index.IndexName }}'.
func (c _StoreImpl) {{ $model }}_By{{$col.Name}} ({{$param}} {{$colType}}) (*{{ $model }},bool){
	o ,ok :=RowCacheIndex.Get("{{ $indexName }}:"+fmt.Sprintf("%v",{{$param}}))
	if ok {
		if obj, ok := o.(*{{ $model }});ok{
			return obj, true
		}
	}

	row, err := New{{ $model }}_Selector().{{$col.Name}}_EQ({{$param}}).GetRow(base.DB)
	if err == nil{
        RowCacheIndex.Set("{{ $indexName }}:"+fmt.Sprintf("%v",row.{{$param}}), row,0)
        return row, true
    }

	XOLogErr(err)
	return nil, false
}
{{$param := (printf "%ss" $col.Name) }}
func (c _StoreImpl) PreLoad{{ $model }}_By{{$col.Name}}s ({{$param}} []{{$colType}}) {
	not_cached := make([]{{$colType}},0,len({{$param}}))

	for _,id := range {{$param}} {
		_ ,ok :=RowCacheIndex.Get("{{ $indexName }}:"+fmt.Sprintf("%v",id))
		if !ok {
			not_cached = append(not_cached,id)
		}
	}

	if len(not_cached) > 0 {
		rows, err := New{{ $model }}_Selector().{{$col.Name}}_In(not_cached).GetRows(base.DB)
		if err == nil{
            for _, row := range rows {
                RowCacheIndex.Set("{{ $indexName }}:"+fmt.Sprintf("%v",row.{{$col.Name}}), row,0)
            }
        }
	}
}
{{else}}
// {{$model}} - {{.Index.IndexName}}

{{end}}