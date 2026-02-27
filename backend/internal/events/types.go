package events

import "github.com/crimsonn/zm_server/internal/models"

var (
	PlayerDisconnected string = "player.disconnected"
	PlayerJoined       string = "player.joined"
	PlayerAddToSession string = "player.add_to_session"
)

type PlayerDisconnectedPayload struct {
	NetID  string `json:"netID"`
	Reason string `json:"reason"`
}

type PlayerJoinedPayload struct {
	NetID      string `json:"netID"`
	Name       string `json:"name"`
	Identifier string `json:"identifier"`
}

type PlayerAddToSessionPayload struct {
	NetID string
	User  models.User
}
