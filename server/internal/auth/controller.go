package auth

import (
	"log/slog"
	"net/http"
	"server/internal/models"
	"server/pkg/betools"
)

type Controller struct {
	svc *Service
}

func NewController(svc *Service) *Controller {
	return &Controller{
		svc: svc,
	}
}

func (c *Controller) GetRoutes() []betools.Route {
	return []betools.Route{
		{
			Method:      "POST",
			Pattern:     "/auth/register",
			HandlerFunc: c.handleRegister,
			Middlewares: []betools.Middleware{
				betools.BodyParser[models.RegisterRequest](),
			},
		},
		{
			Method:      "POST",
			Pattern:     "/auth/login",
			HandlerFunc: c.handleLogin,
			Middlewares: []betools.Middleware{
				betools.BodyParser[models.LoginRequest](),
			},
		},
		{
			Method:      "POST",
			Pattern:     "/auth/logout",
			HandlerFunc: c.handleLogout,
		},
	}
}

func (c *Controller) handleRegister(w http.ResponseWriter, r *http.Request) {
	req := betools.GetBodyCtx[models.RegisterRequest](r)

	if err := c.svc.Register(req); err != nil {
	}

	betools.SendCreatedResponse(w)
}

func (c *Controller) handleLogin(w http.ResponseWriter, r *http.Request) {
	req := betools.GetBodyCtx[models.LoginRequest](r)

	slog.Debug("login", "email", req.Email, "password", req.Password)

	betools.SendOKResponse(w, "token")
}

func (c *Controller) handleLogout(w http.ResponseWriter, r *http.Request) {
}
