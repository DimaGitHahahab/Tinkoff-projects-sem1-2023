package model

type Device struct {
	SerialNum string `json:"serial_number"`
	Model     string `json:"model"`
	IP        string `json:"ip"`
}
