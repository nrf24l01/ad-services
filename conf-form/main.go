package main

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

	echoMw "github.com/labstack/echo/v4/middleware"
	"github.com/nrf24l01/ad-services/conf-form/config"
	"github.com/nrf24l01/ad-services/conf-form/handlers"
	"github.com/nrf24l01/ad-services/conf-form/routes"

	echokitMw "github.com/nrf24l01/go-web-utils/echokit/middleware"
	echokitSchemas "github.com/nrf24l01/go-web-utils/echokit/schemas"
	"github.com/nrf24l01/go-web-utils/pgkit"

	gologger "github.com/nrf24l01/go-logger"
)

func main() {
	ctx := context.Background()

	// Logger create
	logger := gologger.NewLogger(os.Stdout, "conf-form",
		gologger.WithTypeColors(map[gologger.LogType]string{
			gologger.LogType("HTTP"):  gologger.BgCyan,
			gologger.LogType("DB"):    gologger.BgGreen,
			gologger.LogType("SETUP"): gologger.BgRed,
			gologger.LogType("AUTH"):  gologger.BgMagenta,
		}),
	)
	log.Printf("Logger initialized")

	err := godotenv.Load(".env")
	if err != nil {
		logger.Log(gologger.LevelWarn, gologger.LogType("SETUP"), fmt.Sprintf("Failed to load .env file: %v", err), "")
	} else {
		logger.Log(gologger.LevelSuccess, gologger.LogType("SETUP"), ".env file loaded", "")
	}

	// Configuration initialization
	config, err := config.BuildConfigFromEnv()
	if err != nil {
		logger.Log(gologger.LevelFatal, gologger.LogType("SETUP"), fmt.Sprintf("Failed to build config: %v", err), "")
		return
	} else {
		logger.Log(gologger.LevelSuccess, gologger.LogType("SETUP"), "Configuration loaded", "")
	}

	// Data sources initialization
	db, err := pgkit.NewDB(ctx, config.PGConfig)
	if err != nil {
		logger.Log(gologger.LevelFatal, gologger.LogType("SETUP"), fmt.Sprintf("Failed to connect to postgres: %v", err), "")
		return
	} else {
		logger.Log(gologger.LevelSuccess, gologger.LogType("SETUP"), "Connected to Postgres database", "")
	}
	// err = pgkit.RunMigrations(db.SQL, config.PGConfig)
	// if err != nil {
	// 	logger.Log(gologger.LevelFatal, gologger.LogType("SETUP"), fmt.Sprintf("Failed to run migrations: %v", err), "")
	// 	return
	// } else {
	// 	logger.Log(gologger.LevelSuccess, gologger.LogType("SETUP"), "Migrations ran successfully", "")
	// }

	// Create echo object
	e := echo.New()

	// Set up HTML template renderer
	tpl := template.Must(template.ParseGlob("templates/*.html"))
	e.Renderer = &Template{templates: tpl}

	// Register custom validator
	v := validator.New()
	e.Validator = &echokitMw.CustomValidator{Validator: v}

	// Echo Configs
	e.Use(echoMw.Recover())
	e.Use(echoMw.RemoveTrailingSlash())
	e.Use(echokitMw.TraceMiddleware())

	e.Use(echokitMw.RequestLogger(logger))

	// Cors
	log.Printf("Setting allowed origin to: %s", config.WebAppConfig.AllowOrigin)
	e.Use(echoMw.CORSWithConfig(echoMw.CORSConfig{
		AllowOrigins:     []string{config.WebAppConfig.AllowOrigin},
		AllowMethods:     []string{echo.GET, echo.POST, echo.OPTIONS, echo.DELETE},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	// Api group
	api := e.Group("")

	// Health check endpoint
	api.GET("/ping", func(c echo.Context) error {
		return c.JSON(200, echokitSchemas.Message{Status: "Sl-eco/bank backend is OK"})
	})
	api.GET("/", func(c echo.Context) error {
		return c.Render(200, "templates/new-form.html", nil)
	})

	// Register routes
	handler := &handlers.Handler{DB: db, Config: config, Logger: logger}
	routes.RegisterRoutes(api, handler)

	// Start server
	e.Logger.Fatal(e.Start(config.WebAppConfig.AppHost))
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmplName := filepath.Base(name)
	return t.templates.ExecuteTemplate(w, tmplName, data)
}
