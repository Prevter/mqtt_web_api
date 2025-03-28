package models

type Measurement struct {
	Id             int64          `json:"id_measurment"`
	Time           string         `json:"measurment_time"`
	Value          float64        `json:"measurment_value"`
	IdStation      int64          `json:"id_station"`
	StationName    string         `json:"station_name"`
	IdMeasuredUnit int64          `json:"id_measured_unit"`
	UnitTitle      string         `json:"unit_title"`
	UnitUnit       string         `json:"unit_unit"`
	UnitMessage    JsonNullString `json:"unit_message"`
	UnitOrder      JsonNullInt32  `json:"unit_order"`
	Designation    JsonNullString `json:"designation"`
}
