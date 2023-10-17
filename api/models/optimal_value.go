package models

type OptimalValue struct {
	IdCategory     int64         `json:"id_category"`
	IdMeasuredUnit int64         `json:"id_measured_unit"`
	BottomBorder   int           `json:"bottom_border"`
	UpperBorder    JsonNullInt32 `json:"upper_border"`
}
