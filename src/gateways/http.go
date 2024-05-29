package gateways

import (
	service "bn-survey-point/src/services"

	"github.com/gofiber/fiber/v2"
)

type HTTPGateway struct {
	surveyPointService service.ISurveyPointService
}

func NewHTTPGateway(app *fiber.App, service service.ISurveyPointService) {
	gateway := &HTTPGateway{
		surveyPointService: service,
	}

	GatewayUsers(*gateway, app)
}
