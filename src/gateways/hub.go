package gateways

import (
	"encoding/json"
	"fmt"
	"websocketjingjing/domain/entities"

	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func (s Gateway) CreateRoom(ctx *fiber.Ctx) error {
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

func (s Gateway) JoinRoom(c *websocket.Conn) {
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
}

func (s Gateway) LeaveRoom(c *websocket.Conn) {
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

// func (s Gateway) Boardcast(c *websocket.Conn) {
// 	s.HubService.Boardcast(message)
// }

// func (s Gateway) WriteMessage(c *websocket.Conn) {
// 	s.HubService.WriteMessage(c,message)
// }

// func (s Gateway) ReadMessage(c *websocket.Conn) {
// 	s.HubService.ReadMessage(c, client)
// }

func (s Gateway) GetRooms(ctx *fiber.Ctx) error {
	rooms := s.HubService.GetRooms()

	return ctx.Status(fiber.StatusOK).JSON(rooms)
}

func (s Gateway) GetClients(ctx *fiber.Ctx) error {
	roomId := ctx.Query("roomId")
	clients,err := s.HubService.GetClients(roomId)

	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(entities.ResponseMessage{Message: "clients not found"})
	}

	return ctx.Status(fiber.StatusOK).JSON(clients)
}

func (s Gateway) GetRoom(ctx *fiber.Ctx) error {
	roomId := ctx.Query("roomId")
	room := s.HubService.GetRoom(roomId)
	if room == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(entities.ResponseMessage{Message: "room not found"})
	}
	return ctx.Status(fiber.StatusOK).JSON(room)
}

func (s Gateway) handleWebSocket(conn *websocket.Conn) {
	// Get room name and username from the URL parameters
	roomId := conn.Query("room_id")
	userId := conn.Query("userId")
	username := conn.Query("username")
	role := conn.Query("role")
	if role == "" {
		role = "client"
	}
	
	user := &entities.Client{
		Conn:     conn,
		ID:       userId,
		RoomID:   roomId,
		Username: username,
		Role:     role,
	}
	fmt.Println(user.RoomID)
	s.HubService.JoinRoom(conn, *user)
	defer s.LeaveRoom(conn)
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(fmt.Sprintf("Error: %s", err))
			break
		}
		fmt.Println("Received:", string(msg))

		// Create a JSON message including the username
		// message := entities.Message{
		// 	Content:   string(msg),
		// 	Username:  username,
		// 	Role:      role,
		// }
		// messageBytes, _ := json.Marshal(message)

		message := map[string]string{
			"username": username,
			"message":  string(msg),
			"role":     role,
		}
		messageBytes, _ := json.Marshal(message)

		s.HubService.Boardcast(messageBytes,roomId)
	}
}

// func(s Gateway) Roomhandle(c *websocket.Conn) {
// 	var prevRooms []string // Store previous room names

// 	for {
// 		// Create a new message when there is a change in chatServer
// 		rooms := s.HubService.GetRooms()

// 		// Check for changes
// 		if len(*rooms) != len(prevRooms) {
// 			message := fiber.Map{
// 				"rooms": rooms,
// 			}
// 			// Send the updated room list to the client
// 			err := c.WriteJSON(message)
// 			if err != nil {
// 				log.Println(err)
// 				break
// 			}

// 			prevRooms = rooms
// 		}

// 		// Wait for 1 second before checking again
// 		time.Sleep(1 * time.Second)
// 	}
// }