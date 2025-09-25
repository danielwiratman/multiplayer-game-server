package auth

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
			Pattern:     "/auth/register",
			HandlerFunc: c.handleRegister,
		},
		{
			Method:      "POST",
			Pattern:     "/auth/login",
			HandlerFunc: c.handleLogin,
		},
		{
			Method:      "POST",
			Pattern:     "/auth/logout",
			HandlerFunc: c.handleLogout,
		},
	}
}

func (c *Controller) handleRegister(w http.ResponseWriter, r *http.Request) {
}

func (c *Controller) handleLogin(w http.ResponseWriter, r *http.Request) {
}

func (c *Controller) handleLogout(w http.ResponseWriter, r *http.Request) {
}
