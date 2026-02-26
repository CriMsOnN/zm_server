package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	ErrInvalidJSON   = errors.New("invalid json")
	ErrInvalidEvent  = errors.New("invalid event")
	ErrUnknownDomain = errors.New("unknown domain")
	ErrUnknownEvent  = errors.New("unknown event")
	ErrUnknownError  = errors.New("unknown error")
)

type clientSocketMessage struct {
	Event string          `json:"event"`
	Data  json.RawMessage `json:"data,omitempty"`
}

var clientWebsocketUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *ServerHandlers) HandleClientSocket(c *gin.Context) {
	conn, err := clientWebsocketUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.wsLogger.Errorf("Failed to upgrade client websocket connection: %v", err)
		return
	}
	defer conn.Close()

	h.wsLogger.Infof("Client websocket connected from %s", c.ClientIP())
	_ = conn.WriteJSON(gin.H{
		"event": "client.connected",
		"data":  gin.H{"message": "Connected to client namespace from " + c.ClientIP()},
	})

	for {
		_, messageBytes, err := conn.ReadMessage()
		if err != nil {
			h.wsLogger.Infof("Client websocket disconnected: %v", err)
			return
		}

		var msg clientSocketMessage
		if err := json.Unmarshal(messageBytes, &msg); err != nil {
			h.writeError(conn, "invalid_json", ErrInvalidJSON.Error())
			continue
		}

		if msg.Event == "" {
			h.writeError(conn, "invalid_event", ErrInvalidEvent.Error())
			continue
		}

		domain, action, ok := h.splitEventDomain(msg.Event)
		if !ok {
			h.writeError(conn, "invalid_event", ErrInvalidEvent.Error())
			continue
		}

		handler, exists := h.clientHandlers[domain]
		if !exists {
			h.writeError(conn, "unknown_domain", ErrUnknownDomain.Error())
			continue
		}

		handled, responseEvent, responseData, handlerErr := handler.HandleSocketEvent(action, msg.Data)
		if !handled {
			h.writeError(conn, "unknown_event", ErrUnknownEvent.Error())
			continue
		}
		if handlerErr != nil {
			h.writeError(conn, msg.Event+":error", handlerErr.Error())
			continue
		}

		if responseEvent != "" {
			if err := conn.WriteJSON(gin.H{
				"event": domain + "." + responseEvent,
				"data":  responseData,
			}); err != nil {
				h.wsLogger.Errorf("Failed to write client websocket response: %v", err)
				return
			}
		}
	}
}
