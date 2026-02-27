package server

import (
	"encoding/json"
	"strings"

	"github.com/crimsonn/zm_server/internal/logger"
	"github.com/crimsonn/zm_server/internal/repository"
	"github.com/crimsonn/zm_server/internal/user"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type ClientSocketEventHandler interface {
	HandleSocketEvent(action string, data json.RawMessage) (handled bool, responseEvent string, responseData any, err error)
}

type BackendSocketEventHandler interface {
	HandleSocketEvent(action string, data json.RawMessage) (handled bool, responseEvent string, responseData any, err error)
}

type ServerHandlers struct {
	userHandler     *user.UserHandler
	wsLogger        *logrus.Entry
	clientHandlers  map[string]ClientSocketEventHandler
	backendHandlers map[string]BackendSocketEventHandler
}

func NewServerHandlers(db *sqlx.DB) *ServerHandlers {
	userRepository := repository.NewUserRepository(db)
	userService := user.NewUserService(userRepository)
	handlers := &ServerHandlers{
		userHandler: user.NewUserHandler(logger.NewComponentLogger("user"), userService),
		wsLogger:    logger.NewComponentLogger("websocket"),
	}
	handlers.clientHandlers = map[string]ClientSocketEventHandler{
		"user": handlers.userHandler,
	}
	handlers.backendHandlers = map[string]BackendSocketEventHandler{
		"user": handlers.userHandler,
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
