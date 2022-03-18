package routes

import (
	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/echoview-v4"
	"github.com/kreditplus/scorepro/config"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
)

func FrontendRoute(e *echo.Echo, db *gorm.DB) {
	e.Renderer = echoview.New(goview.Config{
		Root:         "views",
		Extension:    ".html",
		DisableCache: false,

	})

	e.GET("/", func(ctx echo.Context) error {
		return ctx.Redirect(http.StatusTemporaryRedirect,"/scorepro/auth/login")
	})

	frontendGroup:=e.Group("/scorepro")
	authController := config.InjectAuthController(db)
	authGroup:=frontendGroup.Group("/auth")
	authGroup.GET("/login", authController.Index)
	authGroup.POST("/login", authController.LoginLos)
}
