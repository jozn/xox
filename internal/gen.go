package internal

//go:generate ./tpl.sh
//go:generate ./gen.sh models

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/alexflint/go-arg"

	"github.com/knq/dburl"

	//_ "ms/xox/loaders"

	//_ "github.com/jozn/xoutil"
    "sort"
)

func Gen() {
	var err error

	// get defaults
	Args = NewDefaultArgs_MS()
	args := Args

	//fmt.Println(os.Args) //me
	// parse args
	arg.MustParse(args)

	// process args
	err = processArgs(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// open database
	err = openDB(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	defer args.DB.Close()

	// load schema name
	//me: schema is equal to name of database in mysql
	if args.Schema == "" {
		args.Schema, err = args.Loader.SchemaName(args) //me: MySchema(..) returns current db
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}

	// load defs into type map
	if args.QueryMode { //me: not for us
		err = args.Loader.ParseQuery(args)
	} else {
		err = GenProccess()
		//err = args.Loader.LoadSchema(args)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// add xo
	err = ExecuteTemplate(XOTemplate, "xo_db", "", args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

    // output
    err = writeTypesJavaJson(args)
    if err != nil {
        fmt.Fprintf(os.Stderr, "error: %v\n", err)
        os.Exit(1)
    }

	// output
	err = writeTypes(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	//fmt.Printf("%+v", args)
}

//me: just configs - skip it
// processArgs processs cli args.
func processArgs(args *ArgType) error {
	var err error

	// get working directory
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	// determine out path
	if args.Out == "" {
		args.Path = cwd
	} else {
		// determine what to do with Out
		fi, err := os.Stat(args.Out)
		if err == nil && fi.IsDir() {
			// out is directory
			args.Path = args.Out
		} else if err == nil && !fi.IsDir() {
			// file exists (will truncate later)
			args.Path = path.Dir(args.Out)
			args.Filename = path.Base(args.Out)

			// error if not split was set, but destination is not a directory
			if !args.SingleFile {
				return errors.New("output path is not directory")
			}
		} else if _, ok := err.(*os.PathError); ok {
			// path error (ie, file doesn't exist yet)
			args.Path = path.Dir(args.Out)
			args.Filename = path.Base(args.Out)

			// error if split was set, but dest doesn't exist
			if !args.SingleFile {
				return errors.New("output path must be a directory and already exist when not writing to a single file")
			}
		} else {
			return err
		}
	}

	// check user template path
	if args.TemplatePath != "" {
		fi, err := os.Stat(args.TemplatePath)
		if err == nil && !fi.IsDir() {
			return errors.New("template path is not directory")
		} else if err != nil {
			return errors.New("template path must exist")
		}
	}

	// fix path
	if args.Path == "." {
		args.Path = cwd
	}

	// determine package name
	if args.Package == "" {
		args.Package = path.Base(args.Path)
	}

	// determine filename if not previously set
	if args.Filename == "" {
		args.Filename = args.Package + args.Suffix
	}

	// if query mode toggled, but no query, read Stdin.
	if args.QueryMode && args.Query == "" {
		buf, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		args.Query = string(buf)
	}

	// query mode parsing
	if args.Query != "" {
		args.QueryMode = true
	}

	// check that query type was specified
	if args.QueryMode && args.QueryType == "" {
		return errors.New("query type must be supplied for query parsing mode")
	}

	// query trim
	if args.QueryMode && args.QueryTrim {
		args.Query = strings.TrimSpace(args.Query)
	}

	// escape all
	if args.EscapeAll {
		args.EscapeSchemaName = true
		args.EscapeTableNames = true
		args.EscapeColumnNames = true
	}

	// if verbose
	if args.Verbose {
		XOLog = func(s string, p ...interface{}) {
			fmt.Printf("SQL:\n%s\nPARAMS:\n%v\n\n", s, p)
		}
	}

	return nil
}

// openDB attempts to open a database connection.
func openDB(args *ArgType) error {
	var err error

	// parse dsn
	u, err := dburl.Parse(args.DSN)
	if err != nil {
		return err
	}

	// save driver type
	args.LoaderType = u.Driver

	// grab loader
	var ok bool
	args.Loader, ok = SchemaLoaders[u.Driver]
	if !ok {
		return errors.New("unsupported database type")
	}

	// open database connection
	args.DB, err = sql.Open(u.Driver, u.DSN)
	if err != nil {
		return err
	}

	return nil
}

// files is a map of filenames to open file handles.
var files = map[string]*os.File{}

// getFile builds the filepath from the TBuf information, and retrieves the
// file from files. If the built filename is not already defined, then it calls
// the os.OpenFile with the correct parameters depending on the state of args.
func getFile(args *ArgType, t *TBuf_OutputToFileHolder) (*os.File, error) {
	var f *os.File
	var err error

	// determine filename
	var filename = strings.ToLower(t.Name) + args.Suffix
	if args.SingleFile {
		filename = args.Filename
	}
	filename = path.Join(args.Path, filename)

	// lookup file
	f, ok := files[filename]
	if ok {
		return f, nil
	}

	// default open mode
	mode := os.O_RDWR | os.O_CREATE | os.O_TRUNC

	// stat file to determine if file already exists
	fi, err := os.Stat(filename)
	if err == nil && fi.IsDir() {
		return nil, errors.New("filename cannot be directory")
	} else if _, ok = err.(*os.PathError); !ok && args.Append && t.TemplateType != XOTemplate {
		// file exists so append if append is set and not XO type
		mode = os.O_APPEND | os.O_WRONLY
	}

	// skip
	if t.TemplateType == XOTemplate && fi != nil {
		return nil, nil
	}

	// open file
	f, err = os.OpenFile(filename, mode, 0666)
	if err != nil {
		return nil, err
	}

	// file didn't originally exist, so add package header
	//me:writes header to each file
	if fi == nil || !args.Append {
		// add build tags
		if args.Tags != "" {
			f.WriteString(`// +build ` + args.Tags + "\n\n")
		}

		// execute
		err = GetTemplateSet().Execute(f, "xo_package.go.tpl", args)
		if err != nil {
			return nil, err
		}
	}

	// store file
	files[filename] = f

	return f, nil
}

// writeTypes writes the generated definitions.
func writeTypes(args *ArgType) error {
	var err error

	//fmt.Println("args.Generated: ", args.Generated)
	out := TBufSlice(args.Generated)

	// sort segments
	sort.Stable(out)
	//sort.Sort(out)

	// loop, writing in order
	for _, t := range out {
		var f *os.File

		// skip when in append and type is XO
		if args.Append && t.TemplateType == XOTemplate {
			continue
		}

		// check if generated template is only whitespace/empty
		bufStr := strings.TrimSpace(t.Buf.String())
		if len(bufStr) == 0 {
			continue
		}

		// get file and filename
		f, err = getFile(args, &t)
		if err != nil {
			return err
		}

		// should only be nil when type == xo
		if f == nil {
			continue
		}

		// write segment
		if !args.Append || (t.TemplateType != TypeTemplate && t.TemplateType != QueryTypeTemplate) {
			_, err = t.Buf.WriteTo(f)
			if err != nil {
				return err
			}
		}
	}

	// build goimports parameters, closing files
	params := []string{"-w"}
	for k, f := range files {
		params = append(params, k)

		// close
		err = f.Close()
		if err != nil {
			return err
		}
	}

	// process written files with goimports
	return exec.Command("goimports", params...).Run()
	//return nil
}

// writeTypes writes the generated definitions.
func writeTypesJavaJson(args *ArgType) error {
    var err error
    var f *os.File

    // check if generated template is only whitespace/empty
    bufStr := strings.TrimSpace(args.GeneratedJavaJson.Buf.String())
    if len(bufStr) == 0 {
        ErrLog(err)
        return nil
    }

    filename := path.Join(args.Path, "j.java")

    mode := os.O_RDWR | os.O_CREATE | os.O_TRUNC
    f, err = os.OpenFile(filename, mode, 0666)
    if err != nil {
        ErrLog(err)
        return  err
    }

    _, err = args.GeneratedJavaJson.Buf.WriteTo(f)
    if err != nil {
        ErrLog(err)
        return err
    }

    return nil
}

