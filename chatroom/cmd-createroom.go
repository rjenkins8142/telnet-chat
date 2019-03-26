package chatroom

import "fmt"

func createRoomHandler(cmd *ParsedCommand) (string, error) {
	ok, emsg, err := checkArgs(cmd, createRoomHelper, 1, 1)
	if !ok {
		return emsg, err
	}
	_, err = CreateRoom(cmd.args[0])
	if err != nil {
		emsg = fmt.Sprintf("Error creating room %s: %s", cmd.args[0], err)
	} else {
		emsg = fmt.Sprintf("Created room [%s]", cmd.args[0])
	}
	cmd.user.DirectMessage(emsg)
	return "", nil
}

func createRoomHelper(cmd *ParsedCommand) (string, error) {
	msg := cmd.command + " <room> - creates a new chat room that can later be joined."
	cmd.user.DirectMessage(msg)
	return "", nil
}
