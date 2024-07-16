package gateways

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func GatewayUsers(gateway HTTPGateway, app *fiber.App, server Server) {
	api := app.Group("/api/surveys")

	api.Get("/add_survey_credits", gateway.AddCredits)
	api.Get("/set_redis_alert_msg", gateway.SetRedisAlertMessage)

	webSocket := app.Group("/api/ws")
	// webSocket.Get("/test", websocket.New(server.SocketTest))
	webSocket.Post("/create_room", server.CreateRoom)
	webSocket.Get("/join_room", websocket.New(server.JoinRoom))
	webSocket.Get("/leave_room", websocket.New(server.LeaveRoom))
	webSocket.Get("/get_rooms", server.GetRooms)
	webSocket.Get("/get_clients", server.GetClients)
	webSocket.Get("/get_room", server.GetRoom)
}
