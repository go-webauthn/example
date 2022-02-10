package model

type MessageResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type DataResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
}
