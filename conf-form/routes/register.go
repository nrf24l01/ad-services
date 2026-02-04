package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/nrf24l01/ad-services/conf-form/handlers"
	gologger "github.com/nrf24l01/go-logger"
)

func RegisterRoutes(e *echo.Group, h *handlers.Handler) {
	h.Logger.Log(gologger.LevelWarn, gologger.LogType("SETUP"), "NO ROUTES REGISTERED", "")
}
