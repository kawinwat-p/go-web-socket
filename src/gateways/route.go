package gateways

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func GatewayUsers(gateway Gateway, app *fiber.App, config websocket.Config) {

	webSocket := app.Group("/api/ws")
	// webSocket.Get("/test", websocket.New(server.SocketTest))
	webSocket.Post("/create_room", gateway.CreateRoom)
	// webSocket.Get("/join_room", websocket.New(gateway.JoinRoom))
	// webSocket.Get("/leave_room", websocket.New(gateway.LeaveRoom))
	webSocket.Get("/get_rooms", gateway.GetRooms)
	webSocket.Get("/get_clients", gateway.GetClients)
	webSocket.Get("/get_room", gateway.GetRoom)

	webSocket.Get("/connect", websocket.New(gateway.handleWebSocket, config))
	// app.Get("/rooms", websocket.New(Roomhandle, config))
}
