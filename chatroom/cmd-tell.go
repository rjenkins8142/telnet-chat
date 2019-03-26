package chatroom

import (
	"fmt"
	"strings"
)

func tellHandler(cmd *ParsedCommand) (string, error) {
	ok, emsg, err := checkArgs(cmd, tellHelper, 2, -1)
	if !ok {
		return emsg, err
	}
	if !existingNick(cmd.args[0]) {
		emsg = fmt.Sprintf("User %s not found", cmd.args[0])
	} else {
		nick := strings.ToUpper(cmd.args[0])
		fullMsg := strings.Join(cmd.args[1:], " ")
	RoomLoop:
		// TODO: Change "nickNames" mapping to point to the user instead of a bool to make this lookup more sensical.
		for _, r := range ListRooms {
			for _, u := range r.users {
				userNick := strings.ToUpper(u.nick)
				if strings.EqualFold(userNick, nick) {
					if cmd.user == u {
						emsg = "Talking to yourself?!?"
					} else {
						u.Message("%s whispers to you: %s", cmd.user.nick, fullMsg)
						emsg = fmt.Sprintf("Sent \"%s\" to %s\n", fullMsg, u.nick)
					}
					break RoomLoop
				}
			}
		}
	}
	cmd.user.DirectMessage(emsg)
	return "", nil
}

func tellHelper(cmd *ParsedCommand) (string, error) {
	msg := cmd.command + " <user> <message> - sends a private message to the specified user."
	cmd.user.DirectMessage(msg)
	return "", nil
}
