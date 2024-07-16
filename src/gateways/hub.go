package gateways

import (
	"websocketjingjing/domain/entities"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"log"
)

func (s Server) CreateRoom(ctx *fiber.Ctx) error {
	var req entities.CreateRoomRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseMessage{Message: "invalid json body"})
	}

	room := entities.Room{
		ID:      req.ID,
		Name:    req.Name,
		Clients: make(map[string]*entities.Client),
	}
	err := s.HubService.CreateRoom(&room)

	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseModel{Message: "cannot create new room."})
	}
	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{Message: "success"})
}

func (s Server) JoinRoom(c *websocket.Conn) {
	log.Println(c.Locals("allowed"))  // true
	log.Println(c.Params("room_id"))  // 123
	log.Println(c.Query("v"))         // 1.0
	log.Println(c.Cookies("session")) // ""

	roomId := c.Query("room_id")
	userId := c.Query("userId")
	username := c.Query("username")
	role := c.Query("role")

	user := &entities.Client{
		Conn:     c,
		ID:       userId,
		RoomID:   roomId,
		Username: username,
		Role:     role,
	}

	log.Printf("roomId: %s", roomId)

	s.HubService.JoinRoom(c, *user)
	// websocket.Conn bindings https://pkg.go.dev/github.com/fasthttp/websocket?tab=doc#pkg-index
	// s.socketService.SocketTestService(c,roomId)
}

func (s Server) LeaveRoom(c *websocket.Conn) {
	log.Println(c.Locals("allowed"))  // true
	log.Println(c.Params("room_id"))  // 123
	log.Println(c.Query("v"))         // 1.0
	log.Println(c.Cookies("session")) // ""

	roomId := c.Query("room_id")
	userId := c.Query("userId")
	username := c.Query("username")

	user := &entities.Client{
		Conn:     c,
		ID:       userId,
		RoomID:   roomId,
		Username: username,
	}

	s.HubService.LeaveRoom(c, *user)
}

// func (s Server) Boardcast(c *websocket.Conn) {
// 	s.HubService.Boardcast(message)
// }

// func (s Server) WriteMessage(c *websocket.Conn) {
// 	s.HubService.WriteMessage(c,message)
// }

// func (s Server) ReadMessage(c *websocket.Conn) {
// 	s.HubService.ReadMessage(c, client)
// }

func (s Server) GetRooms(ctx *fiber.Ctx) error {
	rooms := s.HubService.GetRooms()

	return ctx.Status(fiber.StatusOK).JSON(rooms)
}

func (s Server) GetClients(ctx *fiber.Ctx) error {
	roomId := ctx.Query("roomId")
	clients := s.HubService.GetClients(roomId)

	return ctx.Status(fiber.StatusOK).JSON(clients)
}

func (s Server) GetRoom(ctx *fiber.Ctx) error {
	roomId := ctx.Query("roomId")
	room := s.HubService.GetRoom(roomId)
	if room == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(entities.ResponseMessage{Message: "room not found"})
	}
	return ctx.Status(fiber.StatusOK).JSON(room)
}
