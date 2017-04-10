// Package ischema contains the types for schema 'information_schema'.
package ischema

import "ms/xox/examples/pgcatalog/pgtypes"

// GENERATED BY XO. DO NOT EDIT.

// SQLSizingProfile represents a row from 'information_schema.sql_sizing_profiles'.
type SQLSizingProfile struct {
	Tableoid      pgtypes.Oid            `json:"tableoid"`       // tableoid
	Cmax          pgtypes.Cid            `json:"cmax"`           // cmax
	Xmax          pgtypes.Xid            `json:"xmax"`           // xmax
	Cmin          pgtypes.Cid            `json:"cmin"`           // cmin
	Xmin          pgtypes.Xid            `json:"xmin"`           // xmin
	Ctid          pgtypes.Tid            `json:"ctid"`           // ctid
	SizingID      pgtypes.CardinalNumber `json:"sizing_id"`      // sizing_id
	SizingName    pgtypes.CharacterData  `json:"sizing_name"`    // sizing_name
	ProfileID     pgtypes.CharacterData  `json:"profile_id"`     // profile_id
	RequiredValue pgtypes.CardinalNumber `json:"required_value"` // required_value
	Comments      pgtypes.CharacterData  `json:"comments"`       // comments
}
