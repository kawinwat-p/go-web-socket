package gateways

import (
	service "websocketjingjing/src/services"

	"github.com/gofiber/fiber/v2"
)

type HTTPGateway struct {
	surveyPointService service.ISurveyPointService
}

type Server struct {
	// socketService service.ISocketService
	HubService    service.IHubService
}

func NewHTTPGateway(app *fiber.App, service service.ISurveyPointService, h service.IHubService) {
	gateway := &HTTPGateway{
		surveyPointService: service,
	}
	server := &Server{
		// socketService: socketService,
		HubService:    h,
	}
	GatewayUsers(*gateway, app, *server)
}
