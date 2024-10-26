package websocket

import (
	"log"
	"server/utils"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     utils.CheckOrigin,
	}
)

type Manager struct {
	clients ClientLists
	sync.Mutex
}

func NewManager() *Manager {
	return &Manager{
		clients: make(ClientLists),
	}
}

func (m *Manager) ServeWS(c *gin.Context) {
	log.Println("WEBSOCKET CONNECTION | new connection")

	conn, err := websocketUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		utils.PrintError("WEBSOCKET CONNECTION | error upgrading connection")
	}
	client := NewClient(conn, m)
	m.AddClients(client)
	go client.Read()
	go client.Write()
}

func (m *Manager) AddClients(client *Client) {
	m.Lock()
	defer m.Unlock()
	m.clients[client] = true
}
