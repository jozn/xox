package internal

import (
	"bytes"
	"fmt"
	"text/template"
    "strings"
)

type ProtoFile struct {
	Messages []ProtoMessageDef
	FileName string
    OutPut string
}

type ProtoMessageDef struct {
	Fields      []ProtoMessageFieldDef
	MessageName string
    IsTableNotInline bool //for complex types like user we need to refrence fields with dot (.)
}

type ProtoMessageFieldDef struct {
	TagId    int
	Name     string
    TypeMix SqlToPBType
    PBType string
	Repeat   bool
}

func Gen_ProtosForTables(args *ArgType) {
	tbls := c.Loader.CacheTables
	filePB := ProtoFile{FileName: "pb_tabels"}

	for _, t := range tbls {
		tpb := ProtoMessageDef{
			MessageName: t.Name,
            IsTableNotInline: skipTableModel(t.Name),
		}

		for i, f := range t.Fields {
			fpb := ProtoMessageFieldDef{
				TagId:    (i*2 + 1),
				Name:     f.Name,
                TypeMix: MysqlParseTypeToProtoclBuffer(f.Col.DataType,true),
				Repeat:   false,
			}
			tpb.Fields = append(tpb.Fields, fpb)
		}

		filePB.Messages = append(filePB.Messages, tpb)
	}

    //gen proto def
	tmpl, err := template.New("t").Parse(TMP_PB)
	if err != nil {
		panic(err)
	}
	out := bytes.NewBufferString("")
	err = tmpl.Execute(out, filePB)
	if err != nil {
		panic(err)
	}
	c.GeneratedPb += out.String()

    //gen proto conv
    tmplCon, err := template.New("t").Parse(TMP_PB_CONVERTER)
    if err != nil {
        panic(err)
    }
    outConv := bytes.NewBufferString("")
    err = tmplCon.Execute(outConv, filePB)
    if err != nil {
        panic(err)
    }
    c.GeneratedPbConverter = outConv.String()

	fmt.Println("size of PB (tabels) : ", len(c.Loader.CacheTables))
}

var GRPC_TYOPES_MAP = map[string]string{// go type to => PB types
    "int": "int64",
    "string": "string",
    "float32": "float",
    "float64": "double",
}

const TMP_PB = `
syntax = "proto3";

{{range .Messages}}
message {{.MessageName }}_PB {
    {{- range .Fields }}
    {{.TypeMix.PB }} {{.Name}} = {{.TagId}};
    {{- end }}
}

{{- end}}
`

const TMP_PB_CONVERTER = `
package x

{{range .Messages}}
/*
func PBConv_{{.MessageName }}_PB_To_{{.MessageName }}( o *{{.MessageName }}_PB) *{{.MessageName }} {
  {{- if .IsTableNotInline -}}
   n := &{{.MessageName}}{}
    {{- range .Fields }}
   n.{{.Name}} = {{.TypeMix.Go}} ( o.{{.Name}} )
    {{- end -}}

  {{else }}
     n := &{{.MessageName}}{
    {{- range .Fields }}
      {{.Name}}: {{.TypeMix.Go}} ( o.{{.Name}} ),
      {{- end }}
    }
  {{- end }}
    return n
}

func PBConv_{{.MessageName }}_To_{{.MessageName }}_PB ( o *{{.MessageName }}) *{{.MessageName }}_PB {
  {{- if .IsTableNotInline -}}
   n := &{{.MessageName}}_PB{}
    {{- range .Fields }}
   n.{{.Name}} = {{.TypeMix.GoGen}} ( o.{{.Name}} )
    {{- end -}}

  {{else }}
     n := &{{.MessageName}}_PB{
    {{- range .Fields }}
      {{.Name}}: {{.TypeMix.GoGen}} ( o.{{.Name}} ),
      {{- end }}
    }
  {{- end }}
    return n
}
*/
{{- end}}
`

//============================================================

type SqlToPBType struct{
    Go string //simple go
    GoGen string //go type from pb genrator
    table string
    Java string
    PB string
}

//cp of MyParseType(...)
func MysqlParseTypeToProtoclBuffer(dt string, fromMysql bool) SqlToPBType {
    precision := 0
    unsigned := false

    res := SqlToPBType{}
    // extract unsigned

    if fromMysql {
        if strings.HasSuffix(dt, " unsigned") {
            unsigned = true
            dt = dt[:len(dt)-len(" unsigned")]
        }

        // extract precision
        dt, precision, _ = ParsePrecision(dt)
        _ = precision
        _ = unsigned
    }

    switch dt {
    case "bool", "boolean":
        res = SqlToPBType{
            Go: "bool",
            GoGen: "bool",
            table: "bool",
            Java: "Boolean",
            PB: "bool",
        }

    case "char", "varchar", "tinytext", "text", "mediumtext", "longtext":
        res = SqlToPBType{
            Go: "string",
            GoGen: "string",
            table: "text",
            Java: "String",
            PB: "string",
        }

    case "tinyint", "smallint", "mediumint", "int", "integer":
        res = SqlToPBType{
            Go: "int",
            GoGen: "int32",
            table: "int",
            Java: "Integer",
            PB: "int32",
        }

    case "bigint":
        //the main diffrence is for int64
        res = SqlToPBType{
            Go: "int",
            GoGen: "int64",
            table: "bigint",
            Java: "Long",
            PB: "int64",
        }

    case "float":
        res = SqlToPBType{
            Go: "float32",
            GoGen: "float32",
            table: "float",
            Java: "Float",
            PB: "float",
        }

    case "decimal", "double":
        res = SqlToPBType{
            Go: "float64",
            GoGen: "float64",
            table: "double",
            Java: "Double",
            PB: "double",
        }

    case "binary", "varbinary", "tinyblob", "blob", "mediumblob", "longblob":
        res = SqlToPBType{
            Go: "[]byte",
            GoGen: "[]byte",
            table: "binary",
            Java: "[]byte",
            PB: "????",
        }

    case "timestamp", "datetime", "date", "time":

    default:
    }
return res
}

