package user

import (
	"encoding/json"
	"fmt"
	"net/http"

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

func (h *UserHandler) HandleSocketEvent(action string, _ json.RawMessage) (bool, string, any, error) {
	switch action {
	case "list":
		users, err := h.Service.GetUsers()
		if err != nil {
			h.Logger.Errorf("failed to get users over websocket: %v", err)
			return true, "", nil, fmt.Errorf("failed to fetch users")
		}
		return true, "list.response", users, nil
	default:
		return false, "", nil, nil
	}
}
