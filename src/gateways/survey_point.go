package gateways

import (
	"websocketjingjing/domain/entities"

	"github.com/gofiber/fiber/v2"
)

func (h HTTPGateway) AddCredits(ctx *fiber.Ctx) error {
	response := h.surveyPointService.AddCredits()
	return ctx.Status(response.Status).JSON(response)
}

func (h HTTPGateway) SetRedisAlertMessage(ctx *fiber.Ctx) error {
	err := h.surveyPointService.SetAlertMessageRedis()

	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseModel{Message: "cannot set redis"})
	}

	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{Message: "successfully set redis"})
}
