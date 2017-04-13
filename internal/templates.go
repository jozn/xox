package internal

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"path"
	"text/template"
)

// TemplateLoader loads templates from the specified name.
func TemplateLoader(name string) ([]byte, error) {
	// no template path specified
	if c.TemplatePath == "" {
		return ioutil.ReadFile(path.Join("./templates/", name))
		//return templates.Asset(name)
	}

	return ioutil.ReadFile(path.Join(c.TemplatePath, name))
}

// TemplateSet retrieves the created template set.
func GetTemplateSet() *TemplateSet {
	if c.templateSet == nil {
		c.templateSet = &TemplateSet{
			funcs: NewTemplateFuncs(),
			l:     TemplateLoader,
			tpls:  map[string]*template.Template{},
		}
	}

	return c.templateSet
}

// ExecuteTemplate loads and parses the supplied template with name and
// executes it with obj as the context.
//me: tableName is table or views
func ExecuteTemplate(tt TemplateType, nameOfOutPutFile string, sub string, obj interface{}) error {
	var err error

	//fmt.Println("****** ", tableName)

	// setup generated
	if c.Generated == nil {
		c.Generated = []TBuf_OutputToFileHolder{}
	}

    //me
    if nameOfOutPutFile[0] != 'z'{
        nameOfOutPutFile ="z_" + nameOfOutPutFile
    }
	// create store
	v := TBuf_OutputToFileHolder{
		TemplateType: tt,
		Name:         nameOfOutPutFile, // table name: Post, User
		Subname:      sub,                     // ex: index name
		Buf:          new(bytes.Buffer),
	}

	// build template name
	loaderType := ""
	if tt != XOTemplate {
		loaderType = c.LoaderType + "."
	}
	templateName := fmt.Sprintf("%s%s.go.tpl", loaderType, tt)

	// execute template
	err = GetTemplateSet().Execute(v.Buf, templateName, obj)
	if err != nil {
		return err
	}

	c.Generated = append(c.Generated, v)
	return nil
}

// TemplateSet is a set of templates.
type TemplateSet struct {
	funcs template.FuncMap
	l     func(string) ([]byte, error)
	tpls  map[string]*template.Template
}

// Execute executes a specified template in the template set using the supplied
// obj as its parameters and writing the output to w.
func (ts *TemplateSet) Execute(w io.Writer, name string, obj interface{}) error {
	tpl, ok := ts.tpls[name]
	if !ok {
		// attempt to load and parse the template
		buf, err := ts.l(name)
		if err != nil {
			return err
		}

		// parse template
		tpl, err = template.New(name).Funcs(ts.funcs).Parse(string(buf))
		if err != nil {
			return err
		}
	}

	return tpl.Execute(w, obj)
}
