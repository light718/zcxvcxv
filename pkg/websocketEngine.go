package pkg

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// 创建一个WebSocket升级器
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebSocket连接结构
type WebSocketConnection struct {
	conn *websocket.Conn
	send chan Message
}

// 定义消息结构
type Message struct {
	MessageType int
	Content     []byte
}

// WebSocket服务器结构
type WebSocketServer struct {
	clients    map[*WebSocketConnection]bool
	message    chan Message
	register   chan *WebSocketConnection
	unregister chan *WebSocketConnection
	pool       *sync.Pool
}

// 创建新的WebSocket服务器
func NewWebSocketServer() *WebSocketServer {
	return &WebSocketServer{
		clients:    make(map[*WebSocketConnection]bool),
		message:    make(chan Message, 16384),
		register:   make(chan *WebSocketConnection),
		unregister: make(chan *WebSocketConnection),
		pool: &sync.Pool{
			New: func() interface{} {
				return &WebSocketConnection{}
			},
		},
	}
}

// 启动服务器
func (server *WebSocketServer) Run() {
	for {
		select {
		case connection := <-server.register:
			server.clients[connection] = true
		case connection := <-server.unregister:
			if _, ok := server.clients[connection]; ok {
				delete(server.clients, connection)
				close(connection.send)
			}
		case message := <-server.message:
			for connection := range server.clients {
				select {
				case connection.send <- message:
				default:
					close(connection.send)
					delete(server.clients, connection)
				}
			}
		}
	}
}

// 处理WebSocket连接
func (server *WebSocketServer) HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	connection := server.pool.Get().(*WebSocketConnection)
	connection.conn = conn
	connection.send = make(chan Message, 16)
	server.register <- connection

	go server.handleReads(connection)
	go server.handleWrites(connection)
}

// 处理读取消息
func (server *WebSocketServer) handleReads(connection *WebSocketConnection) {
	defer func() {
		server.unregister <- connection
		connection.conn.Close()
		server.pool.Put(connection)
	}()

	for {
		messageType, message, err := connection.conn.ReadMessage()
		if err != nil {
			break
		}

		server.message <- Message{MessageType: messageType, Content: message}
	}
}

// 处理发送消息
func (server *WebSocketServer) handleWrites(connection *WebSocketConnection) {
	for message := range connection.send {
		err := connection.conn.WriteMessage(message.MessageType, message.Content)
		if err != nil {
			break
		}
	}
}
