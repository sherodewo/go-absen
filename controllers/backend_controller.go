package controllers

import (
	"github.com/foolin/goview/supports/echoview-v4"
	"github.com/kreditplus/scorepro/utils/session"
	"github.com/labstack/echo/v4"
	"net/http"
)

type BaseBackendController struct {
	Title       string
	Menu        string
	BreadCrumbs []map[string]interface{}
}

func Render(ctx echo.Context, title string, view string, activeMenu string, flashMessage session.FlashMessage,
	breadcrumbs []map[string]interface{}, data interface{}) error {
	return echoview.Render(ctx, http.StatusOK, view, echo.Map{
		"title":        title,
		"activeMenu":   activeMenu,
		"breadCrumbs":  breadcrumbs,
		"ctx":          ctx,
		"flashMessage": flashMessage,
		"data":         data,
	})
}
