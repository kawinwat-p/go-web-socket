package gateways

import "github.com/gofiber/fiber/v2"

func GatewayUsers(gateway HTTPGateway, app *fiber.App) {
	api := app.Group("/api/surveys")

	api.Get("/add_survey_credits", gateway.AddCredits)
	api.Get("/set_redis_alert_msg", gateway.SetRedisAlertMessage)
}
