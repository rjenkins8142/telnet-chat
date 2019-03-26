package chatroom

import "fmt"

func deleteRoomHandler(cmd *ParsedCommand) (string, error) {
	ok, emsg, err := checkArgs(cmd, deleteRoomHelper, 1, 1)
	if !ok {
		return emsg, err
	}
	err = RemoveRoom(cmd.args[0])
	if err != nil {
		emsg = fmt.Sprintf("Error removing room %s: %s", cmd.args[0], err)
	} else {
		emsg = fmt.Sprintf("Removed room [%s]", cmd.args[0])
	}
	cmd.user.DirectMessage(emsg)
	return "", nil
}

func deleteRoomHelper(cmd *ParsedCommand) (string, error) {
	msg := cmd.command + " <room> - deletes a chat room. Any users in that chat room will be kicked to the lobby."
	cmd.user.DirectMessage(msg)
	return "", nil
}
