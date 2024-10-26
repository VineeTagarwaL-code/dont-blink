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

// implemeting waiting , paired and disconnected users queues
type Manager struct {
	clients    ClientLists
	waiting    []*Client
	paired     map[*Client]*Client
	register   chan *Client
	unregister chan *Client
	mu         sync.Mutex
}

func NewManager() *Manager {
	return &Manager{
		clients:    make(ClientLists),
		waiting:    make([]*Client, 0),
		paired:     make(map[*Client]*Client, 0),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (m *Manager) ServeWS(c *gin.Context) {
	log.Println("WEBSOCKET CONNECTION | new connection")

	conn, err := websocketUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		utils.PrintError("WEBSOCKET CONNECTION | error upgrading connection")
	}
	client := NewClient(conn, m)
	m.register <- client
	go client.Read()
	go client.Write()
}

func (m *Manager) Start() {

	for {
		select {
		case client := <-m.register:
			m.handleRegistration(client)
		case client := <-m.unregister:
			m.handleUnregistration(client)
		}
	}
}

func (m *Manager) handleRegistration(client *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// add the client in the map of clients first
	m.clients[client] = true

	if len(m.waiting) > 0 {
		partner := m.waiting[0]
		m.waiting = m.waiting[1:]

		m.pairUsers(client, partner)
	} else {
		// Add the client to the waiting queue if no partner is available
		m.waiting = append(m.waiting, client)
	}

}

func (m *Manager) handleUnregistration(client *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.clients[client]; ok {

		delete(m.clients, client)
		close(client.egress)

		// Unpair the client if they are currently paired

		if partner, paired := m.paired[client]; paired {

			m.unpairUsers(client, partner)
		}

		for i, waitingClient := range m.waiting {
			if waitingClient == client {
				m.waiting = append(m.waiting[:i], m.waiting[i+1:]...)
				break
			}
		}

	}

}

func (m *Manager) pairUsers(client1, client2 *Client) {
	m.paired[client1] = client2
	m.paired[client2] = client1

	client1.egress <- []byte(`{"type":"paired","msg":"You are connected to a user!"}`)
	client2.egress <- []byte(`{"type":"paired","msg":"You are connected to a user!"}`)
}

func (m *Manager) unpairUsers(client1, client2 *Client) {
	delete(m.paired, client1)
	delete(m.paired, client2)

	client2.egress <- []byte(`{"type":"unpaired","msg":"Your partner has disconnected."}`)

	m.waiting = append(m.waiting, client2)
}
