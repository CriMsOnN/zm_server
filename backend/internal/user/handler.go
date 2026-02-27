package user

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/crimsonn/zm_server/internal/dto"
	"github.com/crimsonn/zm_server/internal/events"
	"github.com/crimsonn/zm_server/internal/stores"
	"github.com/crimsonn/zm_server/internal/types"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	Logger      *logrus.Entry
	Service     *UserService
	actionMap   map[string]types.SocketActionHandler
	events      events.Bus
	onlineUsers stores.OnlineUserStore
}

func NewUserHandler(logger *logrus.Entry, service *UserService, events events.Bus, onlineUsers stores.OnlineUserStore) *UserHandler {
	h := &UserHandler{
		Logger:      logger,
		actionMap:   make(map[string]types.SocketActionHandler),
		Service:     service,
		events:      events,
		onlineUsers: onlineUsers,
	}
	h.initActions()
	return h
}

func (h *UserHandler) RegisterRoutes(router *gin.Engine) {
	router.GET("/users", h.GetUsers)

}

func (h *UserHandler) registerAction(action string, handler types.SocketActionHandler) {
	h.actionMap[action] = handler
}

func (h *UserHandler) initActions() {
	// h.registerAction("joined", h.handleUserJoined)
	h.registerAction("upsert", h.handleUserUpsert)
	_ = h.events.Subscribe(events.PlayerDisconnected, h.handlePlayerDisconnected)
	_ = h.events.Subscribe(events.PlayerJoined, h.handlePlayerJoined)
}
func (h *UserHandler) handlePlayerDisconnected(event events.Event) {
	data, ok := event.Data.(events.PlayerDisconnectedPayload)
	if !ok {
		h.Logger.Errorf("failed to get player disconnected data: %v", event.Data)
		return
	}
	h.Logger.Infof("Player %s disconnected", data.NetID)
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	h.Logger.Info("Getting users")
	users, err := h.Service.GetUsers()
	if err != nil {
		h.Logger.Errorf("failed to get users: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) handlePlayerJoined(event events.Event) {
	data, ok := event.Data.(events.PlayerJoinedPayload)
	if !ok {
		h.Logger.Errorf("failed to get player joined data: %v", event.Data)
		return
	}
	user, err := h.Service.GetUserByFivemIdentifier(data.Identifier)
	if err != nil {
		h.Logger.Errorf("failed to get user by fivem identifier: %v", err)
		return
	}
	h.events.Publish(events.Event{
		Name: events.PlayerAddToSession,
		Data: events.PlayerAddToSessionPayload{
			NetID: data.NetID,
			User:  *user,
		},
	})
	h.Logger.Infof("Player %s joined", data.NetID)
}

func (h *UserHandler) handleUserUpsert(data json.RawMessage) (bool, any, error) {
	var userUpsert dto.CreateOrUpdateUserDTO
	if err := json.Unmarshal(data, &userUpsert); err != nil {
		h.Logger.Errorf("failed to unmarshal user upsert data: %v", err)
		return true, nil, fmt.Errorf("failed to unmarshal user upsert data")
	}
	err := h.Service.CreateOrUpdateUser(&userUpsert)
	if err != nil {
		h.Logger.Errorf("failed to create or update user: %v", err)
		return true, nil, fmt.Errorf("failed to create or update user")
	}
	h.Logger.Infof("User %s upserted", userUpsert.Name)
	return true, nil, nil
}

func (h *UserHandler) HandleSocketEvent(action string, data json.RawMessage) (bool, any, error) {
	handler, ok := h.actionMap[action]
	if !ok {
		return false, nil, nil
	}
	return handler(data)
}
