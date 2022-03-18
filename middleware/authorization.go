package middleware

import (
	"github.com/kreditplus/scorepro/models"
	"github.com/kreditplus/scorepro/utils"
	"github.com/kreditplus/scorepro/utils/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

type authorizationConfig struct {
	db      *gorm.DB
	Skipper middleware.Skipper
}

func NewAuthorizationMiddleware(db *gorm.DB) *authorizationConfig {
	return &authorizationConfig{
		db: db,
		Skipper: func(context echo.Context) bool {
			return false
		},
	}
}

func (m *authorizationConfig) AuthorizationMiddleware(menu []models.Menu, prefix string) echo.MiddlewareFunc {
	return func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
		return func(context echo.Context) error {
			if m.Skipper(context) {
				return handlerFunc(context)
			}

			//check role
			result, err := session.Manager.Get(context, session.SessionId)
			if err != nil {
				return echo.NewHTTPError(404, "you must login before access this resource")
			}
			userInfo := result.(session.UserInfo)

			for _, v := range menu{
				if v.Prefix == prefix {
					if !utils.ItemExists(strings.Split(v.Role, ","), userInfo.UserRoleID) {
						return echo.NewHTTPError(http.StatusForbidden, "this user role don't have permission to access this resource ")
					}else {
						return handlerFunc(context)
					}
				}
			}
			return echo.NewHTTPError(500, "error get user from context ")
		}
	}
}
