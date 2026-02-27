package user

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/crimsonn/zm_server/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	Logger  *logrus.Entry
	Service *UserService
}

func NewUserHandler(logger *logrus.Entry, service *UserService) *UserHandler {
	return &UserHandler{
		Logger:  logger,
		Service: service,
	}
}

func (h *UserHandler) RegisterRoutes(router *gin.Engine) {
	router.GET("/users", h.GetUsers)
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

func (h *UserHandler) handleUserJoined(data json.RawMessage) (bool, string, any, error) {
	var userJoined dto.UserJoinedDTO
	if err := json.Unmarshal(data, &userJoined); err != nil {
		h.Logger.Errorf("failed to unmarshal user joined data: %v", err)
		return true, "", nil, fmt.Errorf("failed to unmarshal user joined data")
	}
	h.Logger.Infof("User %s joined the server", userJoined.Name)
	return true, "joined.response", userJoined, nil
}

func (h *UserHandler) HandleSocketEvent(action string, data json.RawMessage) (bool, string, any, error) {
	switch action {
	case "list":
		users, err := h.Service.GetUsers()
		if err != nil {
			h.Logger.Errorf("failed to get users over websocket: %v", err)
			return true, "", nil, fmt.Errorf("failed to fetch users")
		}
		return true, "list.response", users, nil
	case "upsert":
		var user dto.CreateOrUpdateUserDTO
		if err := json.Unmarshal(data, &user); err != nil {
			h.Logger.Errorf("failed to unmarshal user data: %v", err)
			return true, "", nil, fmt.Errorf("failed to unmarshal user data")
		}
		err := h.Service.CreateOrUpdateUser(&user)
		if err != nil {
			h.Logger.Errorf("failed to create or update user: %v", err)
			return true, "", nil, fmt.Errorf("failed to create or update user")
		}
		return true, "create.response", user, nil
	case "joined":
		return h.handleUserJoined(data)
	default:
		return false, "", nil, nil
	}
}
