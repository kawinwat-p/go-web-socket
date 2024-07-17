package gateways

import (
	service "websocketjingjing/src/services"

	"github.com/gofiber/fiber/v2"
)

type Gateway struct {
	HubService service.IHubService
}

func NewHTTPGateway(app *fiber.App, h service.IHubService) {
	gateway := &Gateway{
		HubService: h,
	}
	GatewayUsers(*gateway, app)
}
