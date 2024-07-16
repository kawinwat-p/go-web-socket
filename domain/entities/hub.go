package entities

import "github.com/gofiber/contrib/websocket"

type Room struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Clients map[string]*Client `json:"clients"`
}

type RoomResponse struct {
	ID string `json:"id"`
	Name string `json:"name"`
}

type ClientResponse struct {
	ID string `json:"id"`
	RoomID string `json:"roomId"`
	Username string `json:"username"`
}

type Client struct {
	Conn *websocket.Conn
	// Message chan *Message 
	ID string `json:"id"`
	RoomID string `json:"roomId"`
	Username string `json:"username"`
	Role string `json:"role"`
}

type Message struct {
	Content string `json:"content"`
	Username string `json:"username"`
	RoomID string `json:"roomId"`
	Role string `json:"role"`
}

type CreateRoomRequest struct {
	Name string `json:"name"`
	ID string `json:"id"`
}