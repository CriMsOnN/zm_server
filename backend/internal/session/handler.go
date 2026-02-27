package session

import (
	"encoding/json"
	"fmt"

	"github.com/crimsonn/zm_server/internal/events"
	"github.com/crimsonn/zm_server/internal/stores"
	"github.com/crimsonn/zm_server/internal/types"
	"github.com/sirupsen/logrus"
)

type SessionHandler struct {
	Logger      *logrus.Entry
	events      events.Bus
	actionMap   map[string]types.SocketActionHandler
	onlineUsers stores.OnlineUserStore
}

func NewSessionHandler(logger *logrus.Entry, events events.Bus, onlineUsers stores.OnlineUserStore) *SessionHandler {
	h := &SessionHandler{
		Logger:      logger,
		actionMap:   make(map[string]types.SocketActionHandler),
		events:      events,
		onlineUsers: onlineUsers,
	}
	h.initActions()
	return h
}

func (h *SessionHandler) registerAction(action string, handler types.SocketActionHandler) {
	h.actionMap[action] = handler
}

func (h *SessionHandler) initActions() {
	h.registerAction("dropped", h.handleSessionDropped)
	h.registerAction("joined", h.handleSessionJoined)
	_ = h.events.Subscribe(events.PlayerAddToSession, h.handlePlayerAddToSession)
}

func (h *SessionHandler) handlePlayerAddToSession(event events.Event) {
	data, ok := event.Data.(events.PlayerAddToSessionPayload)
	if !ok {
		h.Logger.Errorf("failed to get player add to session data: %v", event.Data)
		return
	}
	if _, ok := h.onlineUsers.Get(data.NetID); ok {
		h.Logger.Errorf("player %s already in online users", data.NetID)
		return
	}

	h.onlineUsers.Set(data.NetID, data.User)
	h.Logger.Infof("Player %s added to session", data.NetID)

}

func (h *SessionHandler) handleSessionJoined(data json.RawMessage) (bool, any, error) {
	var playerJoinedPayload events.PlayerJoinedPayload
	if err := json.Unmarshal(data, &playerJoinedPayload); err != nil {
		h.Logger.Errorf("failed to unmarshal player joined payload: %v", err)
		return true, nil, fmt.Errorf("failed to unmarshal player joined payload")
	}
	h.events.Publish(events.Event{
		Name: events.PlayerJoined,
		Data: playerJoinedPayload,
	})
	return true, nil, nil
}

func (h *SessionHandler) handleSessionDropped(data json.RawMessage) (bool, any, error) {
	var playerDisconnectedPayload events.PlayerDisconnectedPayload
	if err := json.Unmarshal(data, &playerDisconnectedPayload); err != nil {
		h.Logger.Errorf("failed to unmarshal player disconnected payload: %v", err)
		return true, nil, fmt.Errorf("failed to unmarshal player disconnected payload")
	}
	h.events.Publish(events.Event{
		Name: events.PlayerDisconnected,
		Data: playerDisconnectedPayload,
	})

	return true, playerDisconnectedPayload, nil
}

func (h *SessionHandler) HandleSocketEvent(action string, data json.RawMessage) (bool, any, error) {
	handler, ok := h.actionMap[action]
	if !ok {
		return false, nil, nil
	}
	return handler(data)
}
