package models

//type Point struct {
//	Lat float64 `json:"lat"`
//	Lng float64 `json:"lng"`
//}

type Coordinates struct {
	IdStation   int64 `json:"id_station"`
	Coordinates Point `json:"coordinates"`
}
