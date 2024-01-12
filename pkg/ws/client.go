package ws

// ReadSocketMessage this is socket msg struct
type ReadSocketMessage struct {
	Message string            `json:"message"`
	Data    map[string]string `json:"data"`
}
