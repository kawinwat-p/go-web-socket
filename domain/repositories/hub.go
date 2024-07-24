package repositories

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"websocketjingjing/domain/entities"

	"github.com/gofiber/contrib/websocket"
)

type Hub struct {
	Rooms map[string]*entities.Room
	conns map[string]map[*websocket.Conn]bool
	mu    sync.Mutex
	// conns map[*websocket.Conn]bool
	// Rooms map[string]map[*websocket.Conn]bool
}

type IHub interface {
	CreateRoom(room *entities.Room) error
	Boardcast(message []byte, roomId string)
	JoinRoom(c *websocket.Conn, room *entities.Room)
	LeaveRoom(c *websocket.Conn, roomId string)
	// WriteMessage(c *websocket.Conn,message entities.Message)
	// ReadMessage(c *websocket.Conn, client entities.Client)
	GetRooms() *[]entities.RoomResponse
	GetClients(roomID string) (*[]entities.ClientResponse, error)
	GetRoom(roomID string) *entities.Room
}

func NewHub() IHub {
	return &Hub{
		// Rooms: make(map[string]*entities.Room),
		// conns: make(map[*websocket.Conn]bool),
		conns: make(map[string]map[*websocket.Conn]bool),
		Rooms: make(map[string]*entities.Room),
	}
}

func (h *Hub) CreateRoom(room *entities.Room) error {
	// h.Rooms[room.ID] = room
	h.mu.Lock()
	defer h.mu.Unlock()
	h.Rooms[room.ID] = room
	return nil
}

func (h *Hub) JoinRoom(c *websocket.Conn, room *entities.Room) {
	// message := &entities.Message{
	// 	Content:  client.Username + " has joined the room",
	// 	Username: client.Username,
	// 	RoomID:   client.RoomID,
	// 	Role:     client.Role,
	// }

	// var (
	// 	// mt  int
	// 	msg entities.Message
	// 	err error
	// )

	// if err = h.WriteMessage(c, *message); err != nil {
	// 	log.Println("write:", err)
	// }

	// for {
	// 	// log.Println("read")
	// 	if msg, err = h.ReadMessage(c, client); err != nil {
	// 		log.Println("read:", err)
	// 		break
	// 	}
	// 	log.Printf("recv: %s", msg)

	// 	// ss := []byte("test")
	// 	// c.WriteMessage(mt,ss)

	// 	// if string(msg) == "test" {
	// 	// 	c.WriteMessage(mt,[]byte("this is test message"))
	// 	// }

	// 	if err = h.WriteMessage(c, msg); err != nil {
	// 		log.Println("write:", err)
	// 		break
	// 	}

	// }
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.Rooms[room.ID]; ok {
		if _, ok := h.conns[room.ID]; !ok {
			h.conns[room.ID] = make(map[*websocket.Conn]bool)
		}
		// h.conns[room.ID] = make(map[*websocket.Conn]struct{})
		h.conns[room.ID][c] = true
	} else {
		fmt.Println("room not found")
	}
}

func (h *Hub) LeaveRoom(c *websocket.Conn, roomId string) {
	// if _, ok := h.Rooms[client.RoomID]; ok {
	// 	room := h.Rooms[client.RoomID]
	// 	if _, ok := room.Clients[client.ID]; ok {
	// 		if len(h.Rooms[client.RoomID].Clients) != 0 {
	// 			message := entities.Message{
	// 				Content:  "user has left the room",
	// 				RoomID:   client.RoomID,
	// 				Username: client.Username,
	// 			}
	// 			h.WriteMessage(c, message)
	// 		}
	// 		delete(room.Clients, client.ID)
	// 		c.Conn.Close()
	// 	}
	// }
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.Rooms[roomId]; ok {
		delete(h.conns[roomId], c)
		// if len(h.Rooms[roomId]) == 0 {
		// 	delete(h.Rooms, roomId)
		// }
	}
}

func (h *Hub) Boardcast(message []byte, roomId string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	fmt.Println(roomId)
	for client := range h.conns[roomId] {
		err := client.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println(err)
			continue
		}
	}
}

// func (h *Hub) WriteMessage(c *websocket.Conn, message entities.Message) error {
// 	if err := c.WriteJSON(message); err != nil {
// 		log.Println("write:", err)
// 		return err
// 	}
// 	return nil
// }

// func (h *Hub) ReadMessage(c *websocket.Conn, client entities.Client) (entities.Message, error) {
// 	var (
// 		// mt  int
// 		msg []byte
// 		err error
// 	)
// 	if _, msg, err = c.ReadMessage(); err != nil {
// 		log.Println("read:", err)
// 		return entities.Message{}, err
// 	}

// 	message := entities.Message{
// 		Content:  string(msg),
// 		Username: client.Username,
// 		RoomID:   client.RoomID,
// 		Role:     client.Role,
// 	}

// 	return message, nil
// }

func (h *Hub) GetRooms() *[]entities.RoomResponse {
	h.mu.Lock()
	defer h.mu.Unlock()
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

func (h *Hub) GetClients(roomID string) (*[]entities.ClientResponse, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	clients := make([]entities.ClientResponse, 0)

	if _, ok := h.Rooms[roomID]; !ok {
		return nil, errors.New("room not exist")
	}

	room := h.Rooms[roomID]
	for _, c := range room.Clients {
		clients = append(clients, entities.ClientResponse{
			ID:       c.ID,
			Username: c.Username,
			RoomID:   c.RoomID,
		})
	}

	return &clients, nil
}

func (h *Hub) GetRoom(roomID string) *entities.Room {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.Rooms[roomID]; !ok {
		return nil
	}
	room := h.Rooms[roomID]

	return room
}
