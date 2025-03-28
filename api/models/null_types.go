package models

import (
	"database/sql"
	"encoding/json"
)

type JsonNullInt64 struct {
	sql.NullInt64
}

type JsonNullInt32 struct {
	sql.NullInt32
}

type JsonNullString struct {
	sql.NullString
}

func (v JsonNullInt64) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Int64)
	} else {
		return json.Marshal(nil)
	}
}

func (v JsonNullInt32) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Int32)
	} else {
		return json.Marshal(nil)
	}
}

func (v JsonNullString) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.String)
	} else {
		return json.Marshal(nil)
	}
}
