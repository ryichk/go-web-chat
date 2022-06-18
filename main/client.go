package main

import (
	"github.com/gorilla/websocket"
	"time"
)

// client is a chat user
type client struct {
	// socket is a WebSocket for this client
	socket *websocket.Conn
	// send is the channel to which the message is sent
	send chan *message
	// room is chat room in which this client participates
	room *room
	// ユーザに関する情報を保持する
	userData map[string]interface{}
}

func (c *client) read() {
	for {
		var msg *message
		if err := c.socket.ReadJSON(&msg); err == nil {
			msg.When = time.Now()
			msg.Name = c.userData["name"].(string)
			if avatarURL, ok := c.userData["avatar_url"]; ok {
				msg.AvatarURL = avatarURL.(string)
			}
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
