package auth

import (
	"log/slog"
	"net/http"
	"server/internal/env"
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

	slog.Debug("register", "email", req.Email)

	if err := c.svc.Register(req); err != nil {
		slog.Error("register", "email", req.Email, "error", err)
		betools.SendErrorResponse(w, http.StatusBadRequest, "failed to register")
		return
	}

	betools.SendCreatedResponse(w)
}

func (c *Controller) handleLogin(w http.ResponseWriter, r *http.Request) {
	req := betools.GetBodyCtx[models.LoginRequest](r)

	slog.Debug("login", "email", req.Email)

	res, err := c.svc.Login(req)
	if err != nil {
		slog.Error("login", "email", req.Email, "error", err)
		betools.SendErrorResponse(w, http.StatusBadRequest, "failed to login")
		return
	}

	cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    res.RefreshToken,
		Expires:  res.ExpiresIn,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	if env.C.IsProd {
		cookie.Secure = true
	}

	http.SetCookie(w, cookie)

	betools.SendOKResponse(w, models.LoginResponse{
		AccessToken: res.AccessToken,
		TokenType:   "Bearer",
		ExpiresIn:   res.ExpiresIn,
	})
}

func (c *Controller) handleLogout(w http.ResponseWriter, r *http.Request) {
}
