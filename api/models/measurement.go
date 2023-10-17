package models

type Measurement struct {
	Id             int64   `json:"id_measurment"`
	Time           string  `json:"measurment_time"`
	Value          float64 `json:"measurment_value"`
	IdStation      int64   `json:"id_station"`
	IdMeasuredUnit int64   `json:"id_measured_unit"`
}
