package internal

import "database/sql"


// Column represents column info.
type Column_Impl struct {
    FieldOrdinal int            // field_ordinal
    ColumnName   string         // column_name
    DataType     string         // data_type
    NotNull      bool           // not_null
    DefaultValue sql.NullString // default_value
    IsPrimaryKey bool           // is_primary_key
}


// Enum represents a enum.
type Enum_Impl struct {
    EnumName string // enum_name
}

//me: DEL THIS NOT USED
// EnumValue_Impl represents a enum value.
type EnumValue_Impl struct {
    EnumValue  string // enum_value
    ConstValue int    // const_value
}

// ForeignKey represents a foreign key.
type ForeignKey_Impl struct {
    ForeignKeyName string // foreign_key_name
    ColumnName     string // column_name
    RefIndexName   string // ref_index_name
    RefTableName   string // ref_table_name
    RefColumnName  string // ref_column_name
    KeyID          int    // key_id
    SeqNo          int    // seq_no
    OnUpdate       string // on_update
    OnDelete       string // on_delete
    Match          string // match
}

// Index represents an index.
type Index_Impl struct {
    IndexName string // index_name
    IsUnique  bool   // is_unique
    IsPrimary bool   // is_primary
    SeqNo     int    // seq_no
    Origin    string // origin
    IsPartial bool   // is_partial
}


// IndexColumn represents index column info.
type IndexColumn_Impl struct {
    SeqNo      int    // seq_no
    Cid        int    // cid
    ColumnName string // column_name
}


// MyAutoIncrement represents a row from '[custom my_auto_increment]'.
type MyAutoIncrement_Impl struct {
    TableName string // table_name
}

// MyEnumValue represents a row from '[custom my_enum_value]'.
type MyEnumValue_Impl struct {
    EnumValues string // enum_values
}

// Proc represents a stored procedure.
type Proc_Impl struct {
    ProcName   string // proc_name
    ReturnType string // return_type
}

// ProcParam represents a stored procedure param.
type ProcParam_Impl struct {
    ParamType string // param_type
}

// Table represents table info.
type Table_Impl struct {
    Type      string // type
    TableName string // table_name
    ManualPk  bool   // manual_pk
}


