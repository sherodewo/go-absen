package controllers

import (
	"github.com/kreditplus/scorepro/utils/session"
	"github.com/labstack/echo/v4"
)

type HomeController struct {
	BaseBackendController
}

func NewHomeController() HomeController {
	return HomeController{
		BaseBackendController: BaseBackendController{
			Menu:        "Home",
			BreadCrumbs: []map[string]interface{}{},
		},
	}
}

func (c *HomeController) Index(ctx echo.Context) error {
	breadCrumbs := map[string]interface{}{
		"menu": "Home",
		"link": "/scorepro/admin/home",
	}
	userInfo, _ := session.Manager.Get(ctx, session.SessionId)
	return Render(ctx, "Home", "index", c.Menu, session.GetFlashMessage(ctx),
		append(c.BreadCrumbs, breadCrumbs), userInfo)
}
