package handlers

import (
	"github.com/nrf24l01/ad-services/conf-form/config"
	gologger "github.com/nrf24l01/go-logger"
	"github.com/nrf24l01/go-web-utils/pgkit"
)

type Handler struct {
	DB     *pgkit.DB
	Config *config.Config
	Logger *gologger.Logger
}
