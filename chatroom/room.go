package chatroom

import (
	"errors"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Room describes a chat room that can be joined.
type Room struct {
	name     string
	users    []*User
	messages chan string
}

// ListRooms contains the current list of existing rooms.
var ListRooms []*Room

// CreateRoom creates a new room.
func CreateRoom(roomName string) (*Room, error) {
	// Check to see if this room already exists
	existingRoom := FindRoom(roomName)
	if existingRoom != nil {
		return existingRoom, errors.New("Room already exists")
	}
	room := &Room{
		name:     roomName,
		users:    make([]*User, 0),
		messages: make(chan string),
	}
	ListRooms = append(ListRooms, room)
	// Kick off go-routine to listen for new messages in this room.
	go room.listen()
	log.Printf("Created room [%s]", roomName)
	return room, nil
}

// Listen is the main routine that listens for new messages to display to people in this room.
func (r *Room) listen() {
	for {
		msg := <-r.messages
		log.Printf("Received %s room message [%s]", r.name, msg)
		if msg == Cleanup {
			if r.name != "lobby" {
				err := RemoveRoom(r.name)
				if err != nil {
					log.Printf("Warning: Unable to remove room: %s, error: %s", r.name, err)
				}
			}
			break
		}
		r.notifyRoom(msg, false)
	}
}

func (r *Room) notifyRoom(msg string, system bool) {
	for _, person := range r.users {
		// Send a message to each user
		log.Printf("Sending [%s] to user [%s]", msg, person.nick)
		if system {
			person.ServerMessage(msg)
		} else {
			person.SimpleMessage(msg)
		}
	}
}

// RemoveRoom removes an existing room.
func RemoveRoom(roomName string) error {
	if strings.EqualFold(roomName, "lobby") {
		return errors.New("Cannot remove the lobby")
	}
	i := FindRoomIdx(roomName)
	if i >= 0 {
		ourRoom := ListRooms[i]
		// First dump everyone in the room back to the lobby.
		for _, person := range ourRoom.users {
			person.ServerMessage("Deleting room %s, returning you to the lobby", ourRoom.name)
			person.JoinRoomName("lobby")
		}
		// Delete element at i, in a GC safe way.
		copy(ListRooms[i:], ListRooms[i+1:])
		ListRooms[len(ListRooms)-1] = nil
		ListRooms = ListRooms[:len(ListRooms)-1]
		return nil
	}
	return errors.New("Unknown room")
}

// FindRoomIdx finds a room, given its name (case insensitive search). Returns the index to the found room or -1, if not found.
func FindRoomIdx(roomName string) int {
	for i, r := range ListRooms {
		if strings.EqualFold(r.name, roomName) {
			return i
		}
	}
	return -1
}

// FindRoom finds a room, given its name (case insensitive search). Returns pointer to found room or nil, if not found.
func FindRoom(roomName string) *Room {
	i := FindRoomIdx(roomName)
	if i >= 0 {
		return ListRooms[i]
	}
	return nil
}

// Join is how a user gets connected to a room.
func (r *Room) Join(user *User) error {
	r.notifyRoom(fmt.Sprintf("User [%s] has joined the room", user.nick), true)
	r.users = append(r.users, user)
	log.Printf("User %s joined room %s", user.nick, r.name)
	return nil
}

// Leave is how a user leaves a room.
func (r *Room) Leave(user *User) error {
	found := false
	// Loop through and find this user.
	for i, p := range r.users {
		if p == user {
			// Delete element at i, in a GC safe way.
			copy(r.users[i:], r.users[i+1:])
			r.users[len(r.users)-1] = nil
			r.users = r.users[:len(r.users)-1]
			found = true
		}
	}
	if found == false {
		return errors.New("could not find the user to remove from room users list")
	}
	r.notifyRoom(fmt.Sprintf("User [%s] has left the room", user.nick), true)
	return nil
}

// Users returns the list of users in the room.
func (r *Room) Users() []*User {
	return r.users
}

// Name returns the name of the room.
func (r *Room) Name() string {
	return r.name
}

// Send is how we communicate to the room.
func (r *Room) Send(msg string) {
	log.Printf("Sending [%s] to room [%s]", msg, r.name)
	r.messages <- msg
}
