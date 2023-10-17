package models

type ServerStatus string

const (
	Active   ServerStatus = "ACTIVE"
	Inactive ServerStatus = "INACTIVE"
)

type MqttServer struct {
	Id           int64        `json:"id_server"`
	ServerUrl    string       `json:"server_url"`
	ServerStatus ServerStatus `json:"server_status"`
}
