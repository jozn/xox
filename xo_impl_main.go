package main

import (
    "os"
    "fmt"
    "github.com/alexflint/go-arg"
    "ms/xox/internal"
)

func main_old() {
    var err error

    // get defaults
    internal.Args = internal.NewDefaultArgs()
    args := internal.Args

    fmt.Println(os.Args) //me
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
    if args.Schema == "" {
        args.Schema, err = args.Loader.SchemaName(args)
        if err != nil {
            fmt.Fprintf(os.Stderr, "error: %v\n", err)
            os.Exit(1)
        }
    }

    // load defs into type map
    if args.QueryMode {
        err = args.Loader.ParseQuery(args)
    } else {
        err = args.Loader.LoadSchema(args)
    }
    if err != nil {
        fmt.Fprintf(os.Stderr, "error: %v\n", err)
        os.Exit(1)
    }

    // add xo
    err = args.ExecuteTemplate(internal.XOTemplate, "xo_db", "", args)
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
}

