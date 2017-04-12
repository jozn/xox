package internal

import (
	"strings"

	_ "github.com/go-sql-driver/mysql"

	"ms/xox/snaker"
)

func init() {
	SchemaLoaders["mysql"] = TypeLoader{
		ParamN:          func(int) string { return "?" },
		MaskFunc:        func() string { return "?" },
		ProcessRelkind:  MyRelkind,
		Schema:          MySchema,
		ParseType:       MyParseType,
		EnumList:        MyEnums,
		EnumValueList:   MyEnumValues,
		ProcList:        MyProcs,
		ProcParamList:   MyProcParams,
		TableList:       MyTables,
		ColumnList:      MyTableColumns,
		ForeignKeyList:  MyTableForeignKeys,
		IndexList:       MyTableIndexes,
		IndexColumnList: MyIndexColumns,
		QueryColumnList: MyQueryColumns,
	}
}

// MySchema retrieves the name of the current schema.
func MySchema(args *ArgType) (string, error) {
	var err error

	// sql query
	const sqlstr = `SELECT SCHEMA()`

	var schema string

	// run query
	XOLog(sqlstr)
	err = args.DB.QueryRow(sqlstr).Scan(&schema)
	if err != nil {
		return "", err
	}

	return schema, nil
}

// MyRelkind returns the mysql string representation for RelType.
func MyRelkind(relType RelType) string {
	var s string
	switch relType {
	case Table:
		s = "BASE TABLE"
	case View:
		s = "VIEW"
	default:
		panic("unsupported RelType")
	}
	return s
}

// MyParseType parse a mysql type into a Go type based on the column
// definition.
func MyParseType(args *ArgType, dt string, nullable bool) (int, string, string) {
	precision := 0
	nilVal := "nil"
	unsigned := false

	// extract unsigned
	if strings.HasSuffix(dt, " unsigned") {
		unsigned = true
		dt = dt[:len(dt)-len(" unsigned")]
	}

	// extract precision
	dt, precision, _ = args.ParsePrecision(dt)

	var typ string

switchDT:
	switch dt {
	case "bit":
		nilVal = "0"
		if precision == 1 {
			nilVal = "false"
			typ = "bool"
			if nullable {
				nilVal = "sql.NullBool{}"
				typ = "sql.NullBool"
			}
			break switchDT
		} else if precision <= 8 {
			typ = "uint8"
		} else if precision <= 16 {
			typ = "uint16"
		} else if precision <= 32 {
			typ = "uint32"
		} else {
			typ = "uint64"
		}
		if nullable {
			nilVal = "sql.NullInt64{}"
			typ = "sql.NullInt64"
		}

	case "bool", "boolean":
		nilVal = "false"
		typ = "bool"
		if nullable {
			nilVal = "sql.NullBool{}"
			typ = "sql.NullBool"
		}

	case "char", "varchar", "tinytext", "text", "mediumtext", "longtext":
		nilVal = `""`
		typ = "string"
		if nullable {
			nilVal = "sql.NullString{}"
			typ = "sql.NullString"
		}

	case "tinyint", "smallint":
		nilVal = "0"
		typ = "int16"
		if nullable {
			nilVal = "sql.NullInt64{}"
			typ = "sql.NullInt64"
		}

	case "mediumint", "int", "integer":
		nilVal = "0"
		typ = args.Int32Type
		if nullable {
			nilVal = "sql.NullInt64{}"
			typ = "sql.NullInt64"
		}

	case "bigint":
		nilVal = "0"
		typ = "int64"
		if nullable {
			nilVal = "sql.NullInt64{}"
			typ = "sql.NullInt64"
		}

	case "float":
		nilVal = "0.0"
		typ = "float32"
		if nullable {
			nilVal = "sql.NullFloat64{}"
			typ = "sql.NullFloat64"
		}

	case "decimal", "double":
		nilVal = "0.0"
		typ = "float64"
		if nullable {
			nilVal = "sql.NullFloat64{}"
			typ = "sql.NullFloat64"
		}

	case "binary", "varbinary", "tinyblob", "blob", "mediumblob", "longblob":
		typ = "[]byte"

	case "timestamp", "datetime", "date", "time":
		nilVal = "time.Time{}"
		typ = "time.Time"
		if nullable {
			nilVal = "mysql.NullTime{}"
			typ = "mysql.NullTime"
		}

	default:
		if strings.HasPrefix(dt, args.Schema+".") {
			// in the same schema, so chop off
			typ = snaker.SnakeToCamelIdentifier(dt[len(args.Schema)+1:])
			nilVal = typ + "(0)"
		} else {
			typ = snaker.SnakeToCamelIdentifier(dt)
			nilVal = typ + "{}"
		}
	}

	// add 'u' as prefix to type if its unsigned
	// FIXME: this needs to be tested properly...
	if unsigned && IntRE.MatchString(typ) {
		typ = "u" + typ
	}

	return precision, nilVal, typ
}

// MyEnumValues loads the enum values.
func MyEnumValues(db XODB, schema string, enum string) ([]*EnumValue_Impl, error) {
	var err error

	// load enum vals
	res, err := MyEnumValues_Impl(db, schema, enum)
	if err != nil {
		return nil, err
	}

	// process enum vals
	enumVals := []*EnumValue_Impl{}
	for i, ev := range strings.Split(res.EnumValues[1:len(res.EnumValues)-1], "','") {
		enumVals = append(enumVals, &EnumValue_Impl{
			EnumValue:  ev,
			ConstValue: i + 1,
		})
	}

	return enumVals, nil
}

// MyTables returns the MySql tables with the manual PK information added.
// ManualPk is true when the table's primary key is not autoincrement.
func MyTables(db XODB, schema string, relkind string) ([]*Table_Impl, error) {
	var err error

	// get the tables
	rows, err := MyTables_Impl(db, schema, relkind)
	if err != nil {
		return nil, err
	}

	// get the tables that have Autoincrementing included
	autoIncrements, err := MyAutoIncrements(db, schema)
	if err != nil {
		// Set it to an empty set on error.
		autoIncrements = []*MyAutoIncrement{}
	}

	// Add information about manual FK.
	var tables []*Table_Impl
	for _, row := range rows {
		manualPk := true
		// Look for a match in the table name where it contains the autoincrement
		for _, autoInc := range autoIncrements {
			if autoInc.TableName == row.TableName {
				manualPk = false
			}
		}
		tables = append(tables, &Table_Impl{
			TableName: row.TableName,
			Type:      row.Type,
			ManualPk:  manualPk,
		})
	}

	return tables, nil
}

// MyQueryColumns parses the query and generates a type for it.
func MyQueryColumns(args *ArgType, inspect []string) ([]*Column, error) {
	var err error

	// create temporary view xoid
	xoid := "_xo_" + GenRandomID()
	viewq := `CREATE VIEW ` + xoid + ` AS (` + strings.Join(inspect, "\n") + `)`
	XOLog(viewq)
	_, err = args.DB.Exec(viewq)
	if err != nil {
		return nil, err
	}

	// load columns
	cols, err := MyTableColumns(args.DB, args.Schema, xoid)

	// drop inspect view
	dropq := `DROP VIEW ` + xoid
	XOLog(dropq)
	_, _ = args.DB.Exec(dropq)

	// load column information
	return cols, err
}
