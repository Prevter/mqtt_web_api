package models

import (
	"database/sql"
	"encoding/json"
)

type JsonNullInt64 struct {
	sql.NullInt64
}

func (v JsonNullInt64) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Int64)
	} else {
		return json.Marshal(nil)
	}
}

type JsonNullInt32 struct {
	sql.NullInt32
}

func (v JsonNullInt32) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Int32)
	} else {
		return json.Marshal(nil)
	}
}

type Station struct {
	Id            int64         `json:"id_station"`
	City          string        `json:"city"`
	Name          string        `json:"station_name"`
	StationStatus ServerStatus  `json:"station_status"`
	IdServer      JsonNullInt64 `json:"id_server"`
	IdSaveEcoBot  string        `json:"id_saveecobot"`
}
