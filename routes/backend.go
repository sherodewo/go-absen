package routes

import (
	"errors"
	"html/template"
	"reflect"
	"strconv"
	"time"

	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/echoview-v4"
	"github.com/kreditplus/scorepro/config"
	"github.com/kreditplus/scorepro/controllers"
	"github.com/kreditplus/scorepro/middleware"
	"github.com/kreditplus/scorepro/models"
	"github.com/kreditplus/scorepro/utils"
	"github.com/kreditplus/scorepro/utils/session"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

func BackendRoute(e *echo.Echo, db *gorm.DB, dbScoreproV2 *gorm.DB) {

	//=========== Backend ===========//
	var userInfo session.UserInfo
	//new middleware
	mv := echoview.NewMiddleware(goview.Config{
		Root:      "views",
		Extension: ".html",
		Master:    "layouts/master",
		Partials: []string{
			"partials/asside",
			"partials/brand-menu",
			"partials/chart",
			"partials/chat",
			"partials/demo-panel",
			"partials/header-mobile",
			"partials/language",
			"partials/notification",
			"partials/quick-action",
			"partials/quick-panel",
			"partials/quick-panel-toogle",
			"partials/scrolltop",
			"partials/search",
			"partials/sticky-toolbar",
			"partials/sub-header",
			"partials/user-bar",
		},
		Funcs: template.FuncMap{
			"user": func(ctx echo.Context) models.User {
				userModel, err := utils.GetUserInfoFromContext(ctx, db)
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return models.User{}
					}
					return models.User{}
				}
				return userModel
			},
			"floatToString": func(value *float64) string {
				if value == nil {
					return ""
				}
				return strconv.FormatFloat(*value, 'f', -1, 64)
			},
			"floatNotPointerToString": func(value float64) string {
				return strconv.FormatFloat(value, 'f', -1, 64)
			},
			"formatDate": func(date time.Time, layout string) string {
				return date.Format(layout)
			},
			"maritalStatus": func(data string) string {
				var result string
				if data == "S" {
					result = "Belum Pernah Menikah"
				}
				return result
			},
			"getCsrfToken": func(ctx echo.Context) string {
				return ctx.Get("csrf_token").(string)
			},
			"formatRupiah": func(prefix bool, amount interface{}) string {

				if reflect.ValueOf(amount).Kind() == reflect.Ptr {
					if n, ok := amount.(*float64); ok {
						return utils.FormatRupiah(prefix, *n)
					}
					//error set to 0
					return utils.FormatRupiah(false, 0)
				} else {
					if n, ok := amount.(float64); ok {
						return utils.FormatRupiah(prefix, n)
					}
					//error set to 0
					return utils.FormatRupiah(false, 0)
				}
			},
			"MenuParent": func(ctx echo.Context) []map[string]interface{} {
				var dataMenu map[string]interface{}
				var listOfMenu []map[string]interface{}
				result, err := session.Manager.Get(ctx, session.SessionId)
				if err != nil {
					panic(err)
				}
				userInfo = result.(session.UserInfo)
				//var menu models.Menu
				//var menuRole []models.MenuRole
				//if err := db.Raw("select * from web_menu_role where user_role_id = ? and active = ?", userInfo.UserRoleID, 1 ).
				//Scan(&menuRole); err != nil {
				//}

				var listParentMenu []models.Menu
				if err := dbScoreproV2.Raw("select * from web_menu where is_active = ? and menu_type = ? ",
					1, "parent_menu").Scan(&listParentMenu).Error; err != nil {
				}

				for _, v := range listParentMenu {
					var menus []models.Menu
					if err := dbScoreproV2.Raw("select * from web_menu where parent_id = ? and is_active = ? ",
						v.ID, 1).Scan(&menus).Error; err != nil {
						log.Info("ERROR GET MENU BY ROLE ", err.Error())
					}
					dataMenu = map[string]interface{}{
						"Name":  v.Name,
						"Icon":  v.IconClass,
						"Menus": menus,
					}
					listOfMenu = append(listOfMenu, dataMenu)
				}

				return listOfMenu
			},
			"Menu" : func(ctx echo.Context) []map[string]interface{}{
				var dataMenu map[string]interface{}
				var listOfMenu []map[string]interface{}
				result, err := session.Manager.Get(ctx, session.SessionId)
				if err != nil {
					panic(err)
				}
				userInfo = result.(session.UserInfo)

				var menus []models.Menu
				if err := dbScoreproV2.Raw("select * from web_menu where parent_id = ? and is_active = ? and menu_type = ? ",
					0,1,"menu").Scan(&menus).Error; err != nil {
					log.Info("ERROR GET MENU BY ROLE ", err.Error())
				}
				dataMenu = map[string]interface{}{
					"Menus": menus,
				}
				listOfMenu = append(listOfMenu, dataMenu)
				return listOfMenu
			},
		},
		DisableCache: true,
	})

	// You should use helper func `Middleware()` to set the supplied
	// TemplateEngine and make `HTML()` work validly.
	bGroup := e.Group("/scorepro")
	backendGroup := bGroup.Group("/admin", mv, echoMiddleware.CSRFWithConfig(echoMiddleware.CSRFConfig{
		TokenLookup: "form:csrf",
		ContextKey:  "csrf_token",
		Skipper: func(i echo.Context) bool {
			return false
		},
	}), middleware.SessionMiddleware(session.Manager))
	authorizationMiddleware := middleware.NewAuthorizationMiddleware(db)

	var menus []models.Menu
	if err := dbScoreproV2.Raw("select * from web_menu where is_active = ? ",
		 1).Scan(&menus).Error; err != nil {
		log.Info("ERROR GET MENU BY ROLE ", err.Error())
	}

	homeController := controllers.NewHomeController()
	backendGroup.GET("/home", homeController.Index)

	authController := config.InjectAuthController(db)
	backendGroup.POST("/logout", authController.Logout)

	//telcoScoreController
	telcoScoreController := config.InjectTelcoScoreController(db)
	telcoGroup := backendGroup.Group("/score", authorizationMiddleware.AuthorizationMiddleware(menus, "incoming_online"))
	telcoGroup.GET("", telcoScoreController.Index)
	telcoGroup.GET("/datatable", telcoScoreController.List)

	//experianScoringController
	experianScoringController := config.InjectTExperianScoringController(dbScoreproV2)
	experianGroup := backendGroup.Group("/experian", authorizationMiddleware.AuthorizationMiddleware(menus, "experian_scoring"))
	experianGroup.GET("", experianScoringController.Index)
	experianGroup.GET("/datatable", experianScoringController.List)
	experianGroup.GET("/:id", experianScoringController.Detail)
	experianGroup.POST("/store", experianScoringController.Store)
	experianGroup.POST("/update/:id", experianScoringController.Update)

	//userController
	userController := config.InjectUserController(dbScoreproV2)
	userGroup := backendGroup.Group("/register", authorizationMiddleware.AuthorizationMiddleware(menus, "user"))
	userGroup.GET("", userController.Index)
	userGroup.POST("/store", userController.Store)
	userGroup.GET("/add", userController.Add)
	userGroup.GET("/profile", userController.Profile)
	userGroup.GET("/datatable", userController.List)
	userGroup.DELETE("/delete/:id", userController.Delete)
	userGroup.GET("/detail/:id", userController.View)
	userGroup.GET("/edit/:id", userController.Edit)
	userGroup.POST("/update/:id", userController.Update)

	//ConfigController
	configController := config.InjectConfigController(dbScoreproV2)
	configGroup := backendGroup.Group("/config", authorizationMiddleware.AuthorizationMiddleware(menus, "config"))
	configGroup.GET("", configController.Index)
	configGroup.POST("/store", configController.Store)
	configGroup.GET("/datatable", configController.Datatable)
	configGroup.POST("/update/:id", configController.Update)
	configGroup.GET("/:id", configController.Detail)
	configGroup.POST("/delete/:id", configController.Delete)
	bGroup.POST("/admin/config/set-active/:id", configController.SetActive)
	bGroup.POST("/admin/config/set-inactive/:id", configController.SetInactive)

	// INCOMING OFFLINE
	scoreproController := config.InjectScoreproController(dbScoreproV2)
	{
		offlineGroup := backendGroup.Group("/incoming-offline", authorizationMiddleware.AuthorizationMiddleware(menus, "incoming_offline"))
		offlineGroup.GET("", scoreproController.Index)
		offlineGroup.GET("/datatable", scoreproController.List)
	}

	// SCORE GENERATOR
	generatorController := config.InjectGeneratorController(dbScoreproV2)
	{
		generatorGroup := backendGroup.Group("/score-generator", authorizationMiddleware.AuthorizationMiddleware(menus, "score_generator"))
		generatorGroup.GET("", generatorController.Index)
		generatorGroup.GET("/create", generatorController.Create)
		generatorGroup.GET("/datatable", generatorController.List)
		generatorGroup.GET("/score-model-rules", generatorController.GetScoreModels)
		generatorGroup.GET("/all", generatorController.GetAllScoreGenerator)
	}
	bGroup.POST("/admin/score-generator/store", generatorController.Store)
	bGroup.POST("/admin/score-generator/upload/:replace", generatorController.Upload)

	//Supplier ID
	supplierController := config.InjectSupplierController(dbScoreproV2)
	supplierGroup := backendGroup.Group("/supplier", authorizationMiddleware.AuthorizationMiddleware(menus, "supplier"))
	supplierGroup.GET("", supplierController.Index)
	supplierGroup.POST("/store", supplierController.Store)
	supplierGroup.GET("/datatable", supplierController.List)
	supplierGroup.POST("/update/:id", supplierController.Update)
	supplierGroup.GET("/:id", supplierController.GetByID)
	supplierGroup.POST("/delete/:id", supplierController.Delete)

}
