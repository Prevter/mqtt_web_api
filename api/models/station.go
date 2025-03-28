package models

type Point struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lon"`
}

type Station struct {
	Id            int64          `json:"id_station"`
	City          string         `json:"city"`
	Name          string         `json:"station_name"`
	StationStatus JsonNullString `json:"station_status"`
	ServerUrl     JsonNullString `json:"server_url"`
	ServerStatus  JsonNullString `json:"server_status"`
	IdSaveEcoBot  string         `json:"id_saveecobot"`
	Coordinates   Point          `json:"coordinates"`
}
