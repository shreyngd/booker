package slotserver

type WsServer struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
}

func NewWebsocketServer() *WsServer {
	return &WsServer{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
	}
}

// Run our websocket server, accepting various requests
func (server *WsServer) Run() {
	for {
		select {

		case client := <-server.register:
			server.registerClient(client)

		case client := <-server.unregister:
			server.unregisterClient(client)
		case message := <-server.broadcast:
			server.broadcastToClient(message)
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
