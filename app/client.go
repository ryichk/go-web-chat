package main

import (
	"github.com/gorilla/websocket"
)

// client is a chat user
type client struct {
	// socket is a WebSocket for this client
	socket *websocket.Conn
	// send is the channel to which the message is sent
	send chan []byte
	// room is chat room in which this client participates
	room *room
}

func (c *client) read() {
	for {
		if _, msg, err := c.socket.ReadMessage(); err == nil {
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
