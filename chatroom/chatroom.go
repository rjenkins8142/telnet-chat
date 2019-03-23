package chatroom

import "net"

// Room describes a chat room that can be joined.
type Room struct {
	name     string
	users    []*User
	messages chan string
}

// User describes a user connected to the server.
type User struct {
	nick    string
	conn    net.Conn
	room    *Room
	message chan string
	prompt  string
	unread  string
}
