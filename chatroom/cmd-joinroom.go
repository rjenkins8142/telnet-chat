package chatroom

import "fmt"

func joinRoomHandler(cmd *ParsedCommand) (string, error) {
	ok, emsg, err := checkArgs(cmd, joinRoomHelper, 1, 1)
	if !ok {
		return emsg, err
	}
	err = cmd.user.JoinRoomName(cmd.args[0])
	if err != nil {
		emsg = fmt.Sprintf("Error joining room %s: %s", cmd.args[0], err)
	} else {
		emsg = fmt.Sprintf("Joined room [%s]", cmd.args[0])
	}
	cmd.user.DirectMessage(emsg)
	return "", nil
}

func joinRoomHelper(cmd *ParsedCommand) (string, error) {
	msg := cmd.command + " <room> - joins an existing chat room, leaving your current chat room."
	cmd.user.DirectMessage(msg)
	return "", nil
}
