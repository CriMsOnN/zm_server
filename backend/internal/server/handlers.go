package server

import (
	"encoding/json"
	"strings"

	"github.com/crimsonn/zm_server/internal/events"
	"github.com/crimsonn/zm_server/internal/logger"
	"github.com/crimsonn/zm_server/internal/repository"
	"github.com/crimsonn/zm_server/internal/session"
	"github.com/crimsonn/zm_server/internal/stores"
	"github.com/crimsonn/zm_server/internal/user"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type ClientSocketEventHandler interface {
	HandleSocketEvent(action string, data json.RawMessage) (handled bool, responseData any, err error)
}

type BackendSocketEventHandler interface {
	HandleSocketEvent(action string, data json.RawMessage) (handled bool, responseData any, err error)
}

type ServerHandlers struct {
	userHandler     *user.UserHandler
	sessionHandler  *session.SessionHandler
	wsLogger        *logrus.Entry
	clientHandlers  map[string]ClientSocketEventHandler
	backendHandlers map[string]BackendSocketEventHandler
}

func NewServerHandlers(db *sqlx.DB, events events.Bus) *ServerHandlers {
	onlineUsers := stores.NewInMemoryUserStore()
	userRepository := repository.NewUserRepository(db)
	userService := user.NewUserService(userRepository)
	userHandler := user.NewUserHandler(logger.NewComponentLogger("user"), userService, events, onlineUsers)
	sessionHandler := session.NewSessionHandler(logger.NewComponentLogger("session"), events, onlineUsers)
	handlers := &ServerHandlers{
		userHandler:    userHandler,
		wsLogger:       logger.NewComponentLogger("websocket"),
		sessionHandler: sessionHandler,
	}
	handlers.clientHandlers = map[string]ClientSocketEventHandler{
		"user": handlers.userHandler,
	}
	handlers.backendHandlers = map[string]BackendSocketEventHandler{
		"user":    handlers.userHandler,
		"session": handlers.sessionHandler,
	}
	return handlers
}

func (h *ServerHandlers) splitEventDomain(event string) (domain string, action string, ok bool) {
	parts := strings.Split(event, ".")
	if len(parts) != 2 {
		return "", "", false
	}
	return parts[0], parts[1], true
}

func (h *ServerHandlers) RegisterRoutes(router *gin.Engine) {
	h.userHandler.RegisterRoutes(router)
	router.GET("/ws/client", h.HandleClientSocket)
	backend := router.Group("/ws/backend")
	backend.Use(h.BackendSocketJWTMiddleware())
	backend.GET("", h.HandleBackendSocket)
}
