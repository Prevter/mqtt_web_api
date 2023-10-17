package models

type MqttUnit struct {
	IdStation      int64  `json:"id_station"`
	IdMeasuredUnit int64  `json:"id_measured_unit"`
	UnitMessage    string `json:"unit_message"`
	UnitOrder      int    `json:"unit_order"`
}
