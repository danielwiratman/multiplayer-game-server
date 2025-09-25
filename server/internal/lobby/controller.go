package lobby

import (
	"net/http"

	"server/pkg/betools"
)

type Controller struct{}

func NewController() *Controller {
	return &Controller{}
}

func (c *Controller) GetRoutes() []betools.Route {
	return []betools.Route{
		{
			Method:      "POST",
			Pattern:     "/lobby/create",
			HandlerFunc: c.handleCreateLobby,
		},
		{
			Method:      "POST",
			Pattern:     "/lobby/join/{lobbyId}",
			HandlerFunc: c.handleJoinLobby,
		},
		{
			Method:      "POST",
			Pattern:     "/lobby/leave/{lobbyId}",
			HandlerFunc: c.handleLeaveLobby,
		},
	}
}

func (c *Controller) handleCreateLobby(w http.ResponseWriter, r *http.Request) {
}

func (c *Controller) handleJoinLobby(w http.ResponseWriter, r *http.Request) {
}

func (c *Controller) handleLeaveLobby(w http.ResponseWriter, r *http.Request) {
}
