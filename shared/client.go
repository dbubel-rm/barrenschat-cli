package bchatcommon

import (
	"github.com/gorilla/websocket"
	"sync"
	//b "github.com/engineerbeard/barrenschat/shared"
)

type BChatClient struct {
	Name   string
	Room   string
	WsConn *websocket.Conn
	Uid    string
	mu     sync.Mutex
}

func (c *BChatClient) ChangeName(s string) {
	c.Name = s
}

func (c *BChatClient) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.WsConn.Close()
}

func (c *BChatClient) ChangeRoom(s string) {
	c.Room = s
}

func (c *BChatClient) SendMessage(s BMessage) {
	c.mu.Lock()
	c.WsConn.WriteJSON(s)
	c.mu.Unlock()
}

func (c *BChatClient) ReadMessage() (BMessage, error) {
	c.mu.Lock()
	var bMessage BMessage
	err := c.WsConn.ReadJSON(&bMessage)
	c.mu.Unlock()
	return bMessage, err
}
