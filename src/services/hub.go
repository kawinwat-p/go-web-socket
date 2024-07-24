package services

import (
	"fmt"
	"websocketjingjing/domain/entities"
	"websocketjingjing/domain/repositories"

	"github.com/gofiber/contrib/websocket"
)

type hubService struct {
	hub repositories.IHub
}

type IHubService interface {
	CreateRoom(room *entities.Room) error
	Boardcast(message []byte,roomId string)
	JoinRoom(c *websocket.Conn, client entities.Client)
	LeaveRoom(c *websocket.Conn, client entities.Client)
	// WriteMessage(c *websocket.Conn,message entities.Message)
	// ReadMessage(c *websocket.Conn, client entities.Client)
	GetRooms() *[]entities.RoomResponse
	GetClients(RoomID string) (*[]entities.ClientResponse,error)
	GetRoom(roomID string) *entities.Room
}

func NewHubService(h repositories.IHub) IHubService {
	return &hubService{
		hub: h,
	}
}

func (h *hubService) CreateRoom(room *entities.Room) error {
	err := h.hub.CreateRoom(room)
	if err != nil {
		return err
	}
	return nil
}

func (h *hubService) JoinRoom(c *websocket.Conn, client entities.Client) {
	var room *entities.Room
	room = h.GetRoom(client.RoomID)
	fmt.Println(room)
	if room == nil {
		fmt.Println("room not found")
	}else{
		fmt.Println(room.ID)
		room.Clients[client.ID] = &client
		h.hub.JoinRoom(c, room)
	}
}

func (h *hubService) LeaveRoom(c *websocket.Conn, client entities.Client) {
	var room *entities.Room
	if room = h.GetRoom(client.RoomID); room == nil {
		return
	}
	h.hub.LeaveRoom(c, client.RoomID)
}

func (h *hubService) Boardcast(message []byte,roomId string) {
	h.hub.Boardcast(message,roomId)
}

// func (h *hubService) WriteMessage(c *websocket.Conn,message entities.Message) {
// 	h.hub.WriteMessage(c,message)
// }

// func (h *hubService) ReadMessage(c *websocket.Conn, client entities.Client) {
// 	h.hub.ReadMessage(c, client)
// }

func (h *hubService) GetRooms() *[]entities.RoomResponse {
	return h.hub.GetRooms()
}

func (h *hubService) GetClients(RoomID string) (*[]entities.ClientResponse,error) {
	rooms, err := h.hub.GetClients(RoomID)
	if err != nil {
		return nil,err
	}
	return rooms,nil
}

func (h *hubService) GetRoom(roomID string) *entities.Room {
	return h.hub.GetRoom(roomID)
}
