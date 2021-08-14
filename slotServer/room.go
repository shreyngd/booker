package slotserver

import "fmt"

type Room struct {
	name       string
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan *Message
}

// NewRoom creates a new Room
func NewRoom(name string) *Room {
	return &Room{
		name:       name,
		clients:    make(map[*Client]bool),
		register:   make(chan *Client, 100),
		unregister: make(chan *Client, 100),
		broadcast:  make(chan *Message, 1000),
	}
}

// RunRoom runs our room, accepting various requests
func (room *Room) RunRoom() {
	for {
		select {

		case client := <-room.register:
			room.registerClientInRoom(client)

		case client := <-room.unregister:
			room.unregisterClientInRoom(client)

		case message := <-room.broadcast:
			room.broadcastToClientsInRoom(message.encode())
		}
	}
}

func (room *Room) GetName() string {
	return room.name
}

func (room *Room) registerClientInRoom(client *Client) {
	room.notifyClientJoined(client)
	room.clients[client] = true
}

func (room *Room) unregisterClientInRoom(client *Client) {
	if _, ok := room.clients[client]; ok {
		delete(room.clients, client)
	}
}

func (room *Room) broadcastToClientsInRoom(message []byte) {
	for client := range room.clients {
		client.send <- message
	}
}

const welcomeMessage = "%s joined the room"

func (room *Room) notifyClientJoined(client *Client) {
	message := &Message{
		Action:  SendMessageAction,
		Target:  room.name,
		Message: fmt.Sprintf(welcomeMessage, client.GetName()),
	}

	room.broadcastToClientsInRoom(message.encode())
}
