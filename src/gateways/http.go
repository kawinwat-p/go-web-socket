package gateways

import (
	service "websocketjingjing/src/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/contrib/websocket"
)

type Gateway struct {
	HubService service.IHubService
}

func NewHTTPGateway(app *fiber.App, h service.IHubService, config websocket.Config) {
	gateway := &Gateway{
		HubService: h,
	}
	GatewayUsers(*gateway, app,config)
}
