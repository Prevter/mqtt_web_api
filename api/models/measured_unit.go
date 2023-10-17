package models

type MeasuredUnit struct {
	Id    int64  `json:"id_measured_unit"`
	Title string `json:"title"`
	Unit  string `json:"unit"`
}
