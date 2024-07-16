package repositories

import (
	"log"
	"websocketjingjing/domain/entities"

	"github.com/gofiber/contrib/websocket"
)

type Hub struct {
	Rooms map[string]*entities.Room
	conns map[*websocket.Conn]bool
}

type IHub interface {
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

func NewHub() IHub {
	return &Hub{
		Rooms: make(map[string]*entities.Room),
		conns: make(map[*websocket.Conn]bool),
	}
}

func (h *Hub) CreateRoom(room *entities.Room) error {
	h.Rooms[room.ID] = room
	return nil
}

func (h *Hub) JoinRoom(c *websocket.Conn, client entities.Client) {
	if _, ok := h.Rooms[client.RoomID]; ok {
		room := h.Rooms[client.RoomID]
		if _, ok := room.Clients[client.ID]; !ok {
			room.Clients[client.ID] = &client

		}
	}
	message := &entities.Message{
		Content:  client.Username + " has joined the room",
		Username: client.Username,
		RoomID:   client.RoomID,
		Role: client.Role,
	}

	var (
		// mt  int
		msg entities.Message
		err error
	)

	if err = h.WriteMessage(c, *message); err != nil {
		log.Println("write:", err)
	}

	for {
		// log.Println("read")
		if msg, err = h.ReadMessage(c, client); err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", msg)

		// ss := []byte("test")
		// c.WriteMessage(mt,ss)

		// if string(msg) == "test" {
		// 	c.WriteMessage(mt,[]byte("this is test message"))
		// }

		if err = h.WriteMessage(c, msg); err != nil {
			log.Println("write:", err)
			break
		}

	}
}

func (h *Hub) LeaveRoom(c *websocket.Conn, client entities.Client) {
	if _, ok := h.Rooms[client.RoomID]; ok {
		room := h.Rooms[client.RoomID]
		if _, ok := room.Clients[client.ID]; ok {
			if len(h.Rooms[client.RoomID].Clients) != 0 {
				message := entities.Message{
					Content:  "user has left the room",
					RoomID:   client.RoomID,
					Username: client.Username,
				}
				h.WriteMessage(c, message)
			}
			delete(room.Clients, client.ID)
			c.Conn.Close()
		}
	}
}

// func (h *Hub) Boardcast(message entities.Message) {
// 	if _,ok := h.Rooms[message.RoomID]; !ok {
// 		for _, client := range h.Rooms[message.RoomID].Clients {
// 			client.Message <- &message
// 			h.WriteMessage(client.Conn,message)
// 		}
// 	}

// }

func (h *Hub) WriteMessage(c *websocket.Conn, message entities.Message) error {
	if err := c.WriteJSON(message); err != nil {
		log.Println("write:", err)
		return err
	}
	return nil
}

func (h *Hub) ReadMessage(c *websocket.Conn, client entities.Client) (entities.Message, error) {
	var (
		// mt  int
		msg []byte
		err error
	)
	if _, msg, err = c.ReadMessage(); err != nil {
		log.Println("read:", err)
		return entities.Message{}, err
	}

	message := entities.Message{
		Content:  string(msg),
		Username: client.Username,
		RoomID:   client.RoomID,
		Role: client.Role,
	}

	return message, nil
}

func (h *Hub) GetRooms() *[]entities.RoomResponse {
	rooms := make([]entities.RoomResponse, 0)

	for _, r := range h.Rooms {
		room := &entities.RoomResponse{
			ID:   r.ID,
			Name: r.Name,
		}
		rooms = append(rooms, *room)
	}
	return &rooms
}

func (h *Hub) GetClients(roomID string) *[]entities.ClientResponse {
	clients := make([]entities.ClientResponse, 0)

	if _, ok := h.Rooms[roomID]; !ok {
		return &clients
	}

	for _, r := range h.Rooms[roomID].Clients {
		clients = append(clients, entities.ClientResponse{
			ID:       r.ID,
			Username: r.Username,
			RoomID:   r.RoomID,
		})

	}
	return &clients
}

func (h *Hub) GetRoom(roomID string) *entities.Room {

	if _, ok := h.Rooms[roomID]; !ok {
		return nil
	}

	return h.Rooms[roomID]
}
