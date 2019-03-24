package chatroom

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

const (
	// Cleanup is how we tell the different goroutines that we're done and they should exit their respective loops.
	Cleanup = "/cleanup/"
)

// User describes a user connected to the server.
type User struct {
	nick    string
	conn    net.Conn
	room    *Room
	message chan string
	prompt  string
	unread  string
}

var nickNames = make(map[string]bool)

// CreateUser a new user with the given nick and connect them to a room.
func CreateUser(conn net.Conn, room *Room) *User {
	newUser := &User{
		message: make(chan string),
		conn:    conn,
		prompt:  "%s > ",
	}
	nick := newUser.askForNick()
	newUser.nick = nick
	newUser.JoinRoom(room)
	go newUser.commandHandler()
	go newUser.messageHandler()
	newUser.ServerMessage("Welcome %s, type /help for a list of commands.", newUser.nick)
	return newUser
}

func (u *User) askForNick() string {
	var nick string
	var err error
	for {
		u.conn.Write([]byte("Please enter your name: "))
		nick, err = u.readString()
		if err != nil {
			log.Fatal("unable to read user name")
		}
		nick = strings.Trim(nick, "\r\n")
		if !existingNick(nick) {
			addNick(nick)
			break
		}
		u.conn.Write([]byte("That name is already taken, please choose another.\n"))
	}
	return nick
}

func addNick(nick string) {
	allCaps := strings.ToUpper(nick)
	nickNames[allCaps] = true
}

func existingNick(nick string) bool {
	allCaps := strings.ToUpper(nick)
	if _, ok := nickNames[allCaps]; ok {
		return true
	}
	return false
}

// JoinRoom is how a user joins a room. Automatically leaves any previous room.
func (u *User) JoinRoom(room *Room) {
	if u.room != nil {
		u.room.Leave(u)
	}
	room.Join(u)
	u.room = room
}

// JoinRoomName is how a user joins a room by name.
func (u *User) JoinRoomName(roomName string) error {
	room := FindRoom(roomName)
	if room == nil {
		return fmt.Errorf("Unable to find room %s", roomName)
	}
	u.JoinRoom(room)
	return nil
}

// ServerMessage is how we communicate server messages to this user. Takes printf style parameters.
func (u *User) ServerMessage(format string, msg ...interface{}) {
	newMsg := fmt.Sprintf(format, msg...)
	u.SimpleMessage(">> " + newMsg)
}

// Message is how we communicate to this user. Takes printf style parameters.
func (u *User) Message(format string, msg ...interface{}) {
	newMsg := fmt.Sprintf(format, msg...)
	u.SimpleMessage(newMsg)
}

// SimpleMessage is how we communicate to this user.
func (u *User) SimpleMessage(msg string) {
	u.message <- msg
}

// DirectMessage just writes to user.
func (u *User) DirectMessage(msg string) {
	msg = strings.Trim(msg, "\r\n")
	u.conn.Write([]byte(msg + "\n"))
}

func (u *User) commandHandler() {
	for {
		msg, err := u.readString()
		if err != nil {
			log.Printf("Warning: Unable to read input: %s", err)
		}
		msg = strings.Trim(msg, "\r\n")
		msg, err = HandleCommand(msg, u)
		if err != nil {
			log.Printf("Warning: Error during HandleCommand: %s", err)
		}
		if msg == Cleanup {
			u.cleanup()
			break
		}
		if msg != "" {
			msg = fmt.Sprintf("%s | %s | %s", time.Now().Format("2006-01-02 15:04:05"), u.nick, msg)
			u.room.Send(msg)
		} else {
			u.writePrompt()
		}
	}
}

func (u *User) readString() (string, error) {
	reader := bufio.NewReader(u.conn)
	final, err := reader.ReadString('\n')
	return final, err
}

func (u *User) messageHandler() {
	for {
		msg := <-u.message
		if msg == Cleanup {
			break
		}
		_, err := u.conn.Write([]byte("\n" + msg + "\n"))
		u.writePrompt()
		if err != nil {
			log.Printf("Warning: Unable to write message: %s", err)
		}
	}
}

func (u *User) writePrompt() {
	prmt := fmt.Sprintf(u.prompt, u.nick)
	u.conn.Write([]byte(prmt))
}

func (u *User) cleanup() {
	// Tell the current room that we've left.
	u.room.Leave(u)
	// Tell the messageHandler that we're done.
	u.message <- Cleanup
	// Close the connection
	u.conn.Close()
}
