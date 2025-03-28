package models

type MeasurementReport struct {
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
	Avg   float64 `json:"avg"`
	Count int     `json:"count"`
	Unit  string  `json:"unit"`
}

type StationReport struct {
	StationName string         `json:"station_name"`
	FirstDate   JsonNullString `json:"first_date"`
	Units       JsonNullString `json:"units"`
}

type MaxPMByCity struct {
	City string  `json:"city"`
	Unit string  `json:"unit"`
	Max  float64 `json:"max"`
}

type Particles struct {
	Designation string `json:"designation"`
	Count       int    `json:"count"`
}
