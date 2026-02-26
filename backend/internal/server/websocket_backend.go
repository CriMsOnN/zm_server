package server

import (
	"encoding/json"
	"net/http"

	"github.com/crimsonn/zm_server/internal/env"
	"github.com/crimsonn/zm_server/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
		secret := env.GetEnvString("FIVEM_BACKEND_SECRET", "")
		if secret == "" {
			h.wsLogger.Error("FIVEM_BACKEND_SECRET is not configured")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "backend socket secret is not configured"})
			return
		}

		tokenString := utils.ExtractTokenFromHeader(c, "Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing jwt token"})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid jwt token"})
			return
		}

		c.Set("backend_jwt", token)
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

		var message backendSocketMessage
		if err := json.Unmarshal(messageBytes, &message); err != nil {
			h.writeError(conn, "invalid_json", ErrInvalidJSON.Error())
			continue
		}

		switch message.Event {
		case "ping":
			h.writeError(conn, "pong", "")
		default:
			h.writeError(conn, "unknown_event", ErrUnknownEvent.Error())
		}
	}
}
