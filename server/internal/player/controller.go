package player

import (
	"log/slog"
	"net/http"
	"server/internal/middlewares"
	"server/internal/models"
	"server/pkg/betools"
)

type Controller struct {
	mws *middlewares.Middlewares
}

func NewController(mws *middlewares.Middlewares) *Controller {
	return &Controller{
		mws: mws,
	}
}

func (c *Controller) GetRoutes() []betools.Route {
	return betools.WithMiddlewares(
		[]betools.Middleware{
			c.mws.AuthMiddleware(),
		},
		[]betools.Route{
			{
				Method:      "GET",
				Pattern:     "/player/me",
				HandlerFunc: c.handleGetMe,
			},
			{
				Method:      "PUT",
				Pattern:     "/player/me",
				HandlerFunc: c.handleUpdateMe,
				Middlewares: []betools.Middleware{
					betools.BodyParser[models.UpdatePlayerRequest](),
				},
			},
		},
	)
}

func (c *Controller) handleGetMe(w http.ResponseWriter, r *http.Request) {
	account := betools.GetAuthCtx(r)

	slog.Debug("get me", "account", account)

	betools.SendOKResponse(w, account)
}

func (c *Controller) handleUpdateMe(w http.ResponseWriter, r *http.Request) {
}
