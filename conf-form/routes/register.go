package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/nrf24l01/ad-services/conf-form/handlers"
)

func RegisterRoutes(e *echo.Group, h *handlers.Handler) {
	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "templates/new-form.html", nil)
	})
}
