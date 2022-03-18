package main

import (
	"encoding/gob"
	"fmt"
	"github.com/kreditplus/scorepro/config/credential"
	"github.com/kreditplus/scorepro/docs"
	"github.com/kreditplus/scorepro/utils"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"

	"io"
	"os"
	"strings"
	"time"

	"github.com/gorilla/context"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	"github.com/joho/godotenv"
	"github.com/kreditplus/scorepro/config"
	middlewareFunc "github.com/kreditplus/scorepro/middleware"
	"github.com/kreditplus/scorepro/models"
	"github.com/kreditplus/scorepro/routes"
	"github.com/kreditplus/scorepro/utils/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/newrelic/go-agent/v3/integrations/nrecho-v4"
	"github.com/newrelic/go-agent/v3/newrelic"
)

// @contact.name Kredit Plus
// @contact.url https://kreditplus.com
// @contact.email support@kreditplus.com

// @host localhost:9105
// @BasePath /api/v1
// @query.collection.format multi
func main() {
	//Load config.env file
	err := godotenv.Load("conf/config.env")
	if err != nil {
		log.Fatal("ERROR ", err)
	}

	if err := credential.CredentialsConfig(); err != nil {
		log.Fatal("ERROR ", err)
	}

	gob.Register(session.UserInfo{})
	gob.Register(session.FlashMessage{})
	gob.Register(models.User{})
	gob.Register(models.Menu{})
	gob.Register(map[string]interface{}{})
	gob.Register([]models.ValidationError{})

	//New instance echo
	e := echo.New()

	echo.NotFoundHandler = func(c echo.Context) error {
		return c.Render(http.StatusNotFound, "auth/error.html", nil)
	}

	env := strings.ToLower(os.Getenv("APP_ENV"))
	appVersion := strings.ToLower(os.Getenv("APP_VERSION"))

	// newrelic
	if env != "production" {
		e.GET("/swagger/*", echoSwagger.WrapHandler)
		appHost, _ := utils.DecryptCredential(os.Getenv("APP_HOST"))
		appPort, _ := utils.DecryptCredential(os.Getenv("APP_PORT"))
		docs.SwaggerInfo.Title = "SCOREPRO"
		docs.SwaggerInfo.Description = "This is a SCOREPRO API server."
		docs.SwaggerInfo.Version = "1.0"
		docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", appHost, appPort)
		docs.SwaggerInfo.BasePath = "/api/v1"
		docs.SwaggerInfo.Schemes = []string{"http", "https"}
		e.GET("/swagger/*", echoSwagger.WrapHandler)
	} else {
		app, err := newrelic.NewApplication(
			newrelic.ConfigAppName(os.Getenv("APP_NAME")),
			newrelic.ConfigLicense(credential.NewRelicConfigLicense),
			newrelic.ConfigDistributedTracerEnabled(true),
			func(config *newrelic.Config) {
				config.Labels = map[string]string{
					"Env": env,
					"Ver": appVersion,
				}
			},
		)
		if err == nil {
			e.Use(nrecho.Middleware(app))
		}
	}

	e.Static("/scorepro/assets", "assets")

	e.Pre(middleware.RemoveTrailingSlash())

	//Database
	db := config.NewDbMssql()
	dbScoreproV2 := config.NewDbMssqlScorepro()

	// Setup log folder
	if _, err := os.Stat(os.Getenv("LOG_FILE")); os.IsNotExist(err) {
		err = os.MkdirAll(os.Getenv("LOG_FILE"), 0755)
		if err != nil {
			panic(err)
		}
	}
	if _, err := os.Stat("./assets/upload/avatars/"); os.IsNotExist(err) {
		err = os.MkdirAll("./assets/upload/avatars/", 0755)
		if err != nil {
			panic(err)
		}
	}
	// Setup Log
	logPath := os.Getenv("LOG_FILE")
	logFileName := time.Now().Format("2006-01-02") + "-" + "los_cms.log"
	logFile, err := os.OpenFile(logPath+logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Error create or open log file")
	}

	//Validation
	e.Validator = config.NewValidator()

	//Set Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output: io.MultiWriter(logFile, os.Stdout),
	}))

	e.Use(echo.WrapMiddleware(context.ClearHandler))

	session.Manager = session.NewSessionManager(middlewareFunc.NewCookieStore())

	routes.BackendRoute(e, db,dbScoreproV2)
	routes.FrontendRoute(e, dbScoreproV2)
	routes.ApiRoute(e, db, dbScoreproV2)

	// Start server
	if err := e.Start(fmt.Sprintf("%s:%s", credential.AppHost, credential.AppPort)); err != nil {
		e.Logger.Info("shutting down the server")
	}
}
