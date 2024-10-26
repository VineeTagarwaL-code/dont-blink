package websocket

import (
	"server/utils"

	"github.com/gorilla/websocket"
)

type ClientLists map[*Client]bool

type Client struct {
	Conn    *websocket.Conn
	manager *Manager
	egress  chan []byte
}

func NewClient(conn *websocket.Conn, m *Manager) *Client {
	return &Client{
		Conn:    conn,
		manager: m,
		egress:  make(chan []byte),
	}
}

func (c *Client) Read() {
	// setting the readLimit to prevent the client from sending too much data

	c.Conn.SetReadLimit(512)

	for {
		_, payload, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				utils.PrintError("WEBSOCKET CONNECTION | error reading message")
			}
			break
		}
		for wsClient := range c.manager.clients {
			wsClient.egress <- payload
		}

	}

}

func (c *Client) Write() {
	for {
		select {
		case message, ok := <-c.egress:
			if !ok {
				return
			}
			c.Conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}
