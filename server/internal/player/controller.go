package player

import (
	"log/slog"
	"net/http"
	"server/internal/middlewares"
	"server/internal/models"
	"server/pkg/betools"
	"strconv"

	"github.com/go-chi/chi/v5"
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
	return betools.WithMiddlewares(
		[]betools.Middleware{
			middlewares.AuthMiddleware,
		},
		[]betools.Route{
			{
				Method:      "GET",
				Pattern:     "/player/me",
				HandlerFunc: c.handleGetMe,
			},
			{
				Method:      "GET",
				Pattern:     "/player/{uid}",
				HandlerFunc: c.handleGetInfo,
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
	uid := betools.GetAuthCtx(r)

	res, err := c.svc.GetInfo(uid)
	if err != nil {
		slog.Error("get player info", "uid", uid, "error", err)
		betools.SendErrorResponse(w, http.StatusBadRequest, "failed to get player info")
		return
	}

	betools.SendOKResponse(w, res)
}

func (c *Controller) handleGetInfo(w http.ResponseWriter, r *http.Request) {
	uid, err := strconv.Atoi(chi.URLParam(r, "uid"))
	if err != nil {
		slog.Error("get player info", "uid", chi.URLParam(r, "uid"), "error", err)
		betools.SendErrorResponse(w, http.StatusBadRequest, "failed to get player info")
		return
	}

	res, err := c.svc.GetInfo(uid)
	if err != nil {
		slog.Error("get player info", "uid", uid, "error", err)
		betools.SendErrorResponse(w, http.StatusBadRequest, "failed to get player info")
		return
	}

	betools.SendOKResponse(w, res)
}

func (c *Controller) handleUpdateMe(w http.ResponseWriter, r *http.Request) {
	req := betools.GetBodyCtx[models.UpdatePlayerRequest](r)
	uid := betools.GetAuthCtx(r)

	if err := c.svc.UpdateInfo(uid, req); err != nil {
		slog.Error("update player info", "uid", uid, "error", err)
		betools.SendErrorResponse(w, http.StatusBadRequest, "failed to update player info")
		return
	}

	betools.SendOKResponse(w)
}
