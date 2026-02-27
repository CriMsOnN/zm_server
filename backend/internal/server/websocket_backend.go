package server

import (
	"encoding/json"
	"net/http"

	"github.com/crimsonn/zm_server/internal/env"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type backendSocketMessage struct {
	Event string          `json:"event"`
	Data  json.RawMessage `json:"data,omitempty"`
}

var backendWebsocketUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *ServerHandlers) BackendSocketJWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		secret := env.GetEnvString("BACKEND_SECRET", "")
		if secret == "" {
			h.wsLogger.Error("BACKEND_SECRET is not configured")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "backend socket secret is not configured"})
			return
		}

		token := c.Query("secret")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing secret"})
			return
		}

		if token != secret {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid secret"})
			return
		}

		c.Next()
	}
}

func (h *ServerHandlers) HandleBackendSocket(c *gin.Context) {
	conn, err := backendWebsocketUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.wsLogger.Errorf("failed to upgrade backend websocket connection: %v", err)
		return
	}
	defer conn.Close()

	h.wsLogger.Infof("Backend websocket connected from %s", c.ClientIP())
	_ = conn.WriteJSON(gin.H{
		"event": "backend.connected",
		"data":  gin.H{"message": "Connected to backend websocket namespace from " + c.ClientIP()},
	})

	for {
		_, messageBytes, err := conn.ReadMessage()
		if err != nil {
			h.wsLogger.Infof("backend websocket disconnected: %v", err)
			return
		}

		var msg backendSocketMessage
		if err := json.Unmarshal(messageBytes, &msg); err != nil {
			h.writeError(conn, "invalid_json", ErrInvalidJSON.Error())
			continue
		}

		domain, action, ok := h.splitEventDomain(msg.Event)
		if !ok {
			h.writeError(conn, "invalid_event", ErrInvalidEvent.Error())
			continue
		}

		handler, exists := h.backendHandlers[domain]
		if !exists {
			h.writeError(conn, "unknown_domain", ErrUnknownDomain.Error())
			continue
		}

		handled, responseData, handlerErr := handler.HandleSocketEvent(action, msg.Data)
		if !handled {
			h.writeError(conn, "unknown_event", ErrUnknownEvent.Error())
			continue
		}
		if handlerErr != nil {
			h.writeError(conn, msg.Event+":error", handlerErr.Error())
			continue
		}

		if responseData != nil {
			if err := conn.WriteJSON(gin.H{
				"event": msg.Event,
				"data":  responseData,
			}); err != nil {
				h.wsLogger.Errorf("Failed to write backend websocket response: %v", err)
				return
			}
		}
	}
}
