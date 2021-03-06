package slotserver

import (
	"github.com/shreyngd/booker/models"
)

type WsServer struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	rooms      map[*Room]bool
}

func NewWebsocketServer() *WsServer {
	return &WsServer{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client, 100),
		unregister: make(chan *Client, 100),
		broadcast:  make(chan []byte, 1000),
		rooms:      make(map[*Room]bool),
	}
}

// Run our websocket server, accepting various requests
func (server *WsServer) Run() {
	gb := *models.GetInstanceGlobal()
	for {
		select {

		case client := <-server.register:
			server.registerClient(client)
		case client := <-server.unregister:
			server.unregisterClient(client)
		case message := <-server.broadcast:
			server.broadcastToClient(message)
		case message := <-gb.Channel:
			server.createRoomAndBroadCast(message)
		}

	}
}

func (server *WsServer) registerClient(client *Client) {
	server.clients[client] = true
}

func (server *WsServer) unregisterClient(client *Client) {
	delete(server.clients, client)

}

func (server *WsServer) broadcastToClient(message []byte) {
	for client := range server.clients {
		client.send <- message
	}
}

func (server *WsServer) findRoomByName(name string) *Room {
	var foundRoom *Room
	for room := range server.rooms {
		if room.GetName() == name {
			foundRoom = room
			break
		}
	}

	return foundRoom
}

func (server *WsServer) createRoom(name string) *Room {
	room := NewRoom(name)
	go room.RunRoom()
	server.rooms[room] = true
	return room
}

func (server *WsServer) createRoomAndBroadCast(name string) {

	room := server.createRoom(name)
	var roomList []string
	for r := range server.rooms {
		roomList = append(roomList, r.GetName())
	}

	addRoomReply := models.AddRoomReply{
		RoomList:  roomList,
		AddedRoom: room.GetName(),
	}

	addReply := MessageAddRoom{
		Data:   addRoomReply,
		Action: "room-add-success",
	}

	server.broadcast <- addReply.encode()
}
