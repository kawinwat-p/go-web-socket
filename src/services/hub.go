package services

import (
	"websocketjingjing/domain/entities"
	"websocketjingjing/domain/repositories"

	"github.com/gofiber/contrib/websocket"
)

type hubService struct {
	hub repositories.IHub
}

type IHubService interface {
	CreateRoom(room *entities.Room) error
	// Boardcast(message entities.Message)
	JoinRoom(c *websocket.Conn, client entities.Client)
	LeaveRoom(c *websocket.Conn, client entities.Client)
	// WriteMessage(c *websocket.Conn,message entities.Message)
	// ReadMessage(c *websocket.Conn, client entities.Client)
	GetRooms() *[]entities.RoomResponse
	GetClients(RoomID string) *[]entities.ClientResponse
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
	if room := h.GetRoom(client.RoomID); room == nil {
		return
	}
	h.hub.JoinRoom(c, client)
}

func (h *hubService) LeaveRoom(c *websocket.Conn, client entities.Client) {
	h.hub.LeaveRoom(c, client)
}

// func (h *hubService) Boardcast(message entities.Message) {
// 	h.hub.Boardcast(message)
// }

// func (h *hubService) WriteMessage(c *websocket.Conn,message entities.Message) {
// 	h.hub.WriteMessage(c,message)
// }

// func (h *hubService) ReadMessage(c *websocket.Conn, client entities.Client) {
// 	h.hub.ReadMessage(c, client)
// }

func (h *hubService) GetRooms() *[]entities.RoomResponse {
	return h.hub.GetRooms()
}

func (h *hubService) GetClients(RoomID string) *[]entities.ClientResponse {
	return h.hub.GetClients(RoomID)
}

func (h *hubService) GetRoom(roomID string) *entities.Room {
	return h.hub.GetRoom(roomID)
}
