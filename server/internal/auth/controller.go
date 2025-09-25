package auth

import "server/pkg/betools"

type Controller struct{}

func NewController() *Controller {
	return &Controller{}
}

func (c *Controller) GetRoutes() []betools.Route {
	return []betools.Route{}
}
