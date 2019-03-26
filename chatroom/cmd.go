package chatroom

import (
	"fmt"
	"regexp"
	"strings"
)

// Command is a struct that describes a chat command.
type Command struct {
	name        string      // The /slash command to look for
	aliases     []string    // List of other /slash commands that are aliases to this one
	description string      // The description of the command
	handler     commandFunc // The main handler function for the command
}

// ParsedCommand describes the parsed message
type ParsedCommand struct {
	message string   // Raw message
	command string   // The parsed name of the command
	args    []string // The arguments for the command
	user    *User    // The user that sent the command
}

type commandFunc func(cmd *ParsedCommand) (string, error)

var commands = make(map[string]*Command)

// HandleCommand parses any commands and handles any commands for the user.
func HandleCommand(message string, u *User) (string, error) {
	// See if the message possibly contains a command.
	if !strings.Contains(message, "/") {
		return message, nil
	}
	// Split the message into pieces.
	spaces := regexp.MustCompile(`\s+`)
	args := spaces.Split(message, -1)
	commandName := args[0]
	args = args[1:]

	if cmd, ok := commands[commandName]; ok {
		newmsg, err := cmd.handler(&ParsedCommand{
			message: message,
			command: commandName,
			args:    args,
			user:    u,
		})
		return newmsg, err
	}

	return message, nil
}

// RegisterCommand registers a command so we know about it.
func RegisterCommand(command *Command) {
	commands[command.name] = command
	for _, alias := range command.aliases {
		commands[alias] = command
	}
}

// InitCommands initalizes the list of known commands.
func InitCommands() {
	RegisterCommand(&Command{
		name:        "/help",
		aliases:     []string{"/h", "/?"},
		description: "Show help (list of possible commands)",
		handler:     helpHandler,
	})
	RegisterCommand(&Command{
		name:        "/exit",
		aliases:     []string{"/logout", "/lo", "/quit"},
		description: "Exit the chat program",
		handler:     exitHandler,
	})
	RegisterCommand(&Command{
		name:        "/create",
		aliases:     []string{"/createroom"},
		description: "Create a new chat room",
		handler:     createRoomHandler,
	})
	RegisterCommand(&Command{
		name:        "/delete",
		aliases:     []string{"/deleteroom", "/del"},
		description: "Delete a chat room",
		handler:     deleteRoomHandler,
	})
	RegisterCommand(&Command{
		name:        "/join",
		aliases:     []string{"/joinroom"},
		description: "Join a new chat room",
		handler:     joinRoomHandler,
	})
	RegisterCommand(&Command{
		name:        "/list",
		aliases:     []string{"/listrooms"},
		description: "List the available chat rooms",
		handler:     listRoomsHandler,
	})
	RegisterCommand(&Command{
		name:        "/nick",
		aliases:     []string{"/changenick"},
		description: "Change your nickname",
		handler:     changeNickHandler,
	})
	RegisterCommand(&Command{
		name:        "/who",
		aliases:     []string{"/w"},
		description: "Show who is logged in",
		handler:     whoHandler,
	})
	RegisterCommand(&Command{
		name:        "/tell",
		aliases:     []string{"/whisper"},
		description: "Send a direct message to a user",
		handler:     tellHandler,
	})
}

func checkArgs(cmd *ParsedCommand, usage commandFunc, minArgs, maxArgs int) (bool, string, error) {
	for _, arg := range cmd.args {
		switch arg {
		case "-h", "-?":
			msg, err := usage(cmd)
			return false, msg, err
		}
	}
	if maxArgs >= 0 && len(cmd.args) > maxArgs {
		msg := fmt.Sprintf("Too many arguments to %s. Got %d, expected %d.", cmd.command, len(cmd.args), maxArgs)
		cmd.user.DirectMessage(msg)
		msg, err := usage(cmd)
		return false, msg, err
	}
	if len(cmd.args) < minArgs {
		msg := fmt.Sprintf("Not enough arguments to %s. Got %d, expected %d.", cmd.command, len(cmd.args), minArgs)
		cmd.user.DirectMessage(msg)
		msg, err := usage(cmd)
		return false, msg, err
	}
	return true, "", nil
}
