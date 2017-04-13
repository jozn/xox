package internal

//me: this is the copy of :(tl TypeLoader) LoadSchema(args *ArgType) error
func  GenProccess(tl TypeLoader) error {
    var err error

    // load enums
    _, err = tl.LoadEnums(c)
    if err != nil {
        return err
    }

    // load procs
    _, err = tl.LoadProcs(c)
    if err != nil {
        return err
    }

    // load tables
    tableMap, err := tl.LoadRelkind(c, Table)
    if err != nil {
        return err
    }

    // load views
    viewMap, err := tl.LoadRelkind(c, View)
    if err != nil {
        return err
    }

    // merge views with the tableMap
    for k, v := range viewMap {
        tableMap[k] = v
    }

    // load foreign keys
    _, err = tl.LoadForeignKeys(c, tableMap)
    if err != nil {
        return err
    }

    // load indexes
    _, err = tl.LoadIndexes(c, tableMap)
    if err != nil {
        return err
    }

    //Me:
    err = tl.XLoadEvents(c, tableMap)
    if err != nil {
        return err
    }

    return nil
}
