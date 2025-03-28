package models

type MqttServer struct {
	Id           int64          `json:"id_server"`
	ServerUrl    string         `json:"server_url"`
	ServerStatus JsonNullString `json:"server_status"`
}
