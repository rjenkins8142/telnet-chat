package chatroom

import "fmt"

func listRoomsHandler(cmd *ParsedCommand) (string, error) {
	ok, emsg, err := checkArgs(cmd, listRoomsHelper, 0, 0)
	if !ok {
		return emsg, err
	}
	emsg = "Current rooms:\n\n"
	for _, r := range ListRooms {
		cur := ""
		if cmd.user.room == r {
			cur = "* "
		}
		emsg += fmt.Sprintf("%s%s - (%d users)\n", cur, r.name, len(r.users))
	}
	emsg += "\n* - current room\n"
	cmd.user.DirectMessage(emsg)
	return "", nil
}

func listRoomsHelper(cmd *ParsedCommand) (string, error) {
	msg := cmd.command + " - shows a list of all available chat rooms."
	cmd.user.DirectMessage(msg)
	return "", nil
}
