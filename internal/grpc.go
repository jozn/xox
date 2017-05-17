package internal

import (
	"bytes"
	"fmt"
	"text/template"
)

type ProtoFile struct {
	Messages []ProtoMessageDef
	FileName string
    OutPut string
}

type ProtoMessageDef struct {
	Fields      []ProtoMessageFieldDef
	MessageName string
}

type ProtoMessageFieldDef struct {
	TagId    int
	Name     string
	DataType string
    OutType string
	Repeat   bool
}

func Gen_ProtosForTables(args *ArgType) {
	tbls := c.Loader.CacheTables
	filePB := ProtoFile{FileName: "pb_tabels"}

	for _, t := range tbls {
		tpb := ProtoMessageDef{
			MessageName: t.Name,
		}

		for i, f := range t.Fields {
            ot := GRPC_TYOPES_MAP[f.Type]
			fpb := ProtoMessageFieldDef{
				TagId:    (i*2 + 1),
				Name:     f.Name,
				DataType: f.Type,
				OutType: ot,
				Repeat:   false,
			}
			tpb.Fields = append(tpb.Fields, fpb)
		}

		filePB.Messages = append(filePB.Messages, tpb)
	}

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
    {{.OutType}} {{.Name}} = {{  ( (.TagId ) )  }} ;
    {{- end }}
}

{{- end}}
`
