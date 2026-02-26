package server

import "github.com/gorilla/websocket"

type socketError struct {
	Event string `json:"event"`
	Error string `json:"error"`
}

func (h *ServerHandlers) writeError(conn *websocket.Conn, event string, message string) {
	if err := conn.WriteJSON(socketError{
		Event: event,
		Error: message,
	}); err != nil {
		h.wsLogger.Errorf("Failed to write websocket error response: %v", err)
	}
}
