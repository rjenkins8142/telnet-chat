package chatroom

import "fmt"

func whoHandler(cmd *ParsedCommand) (string, error) {
	ok, emsg, err := checkArgs(cmd, whoHelper, 0, 0)
	if !ok {
		return emsg, err
	}
	emsg = "Current users:\n\n"
	cnt := 0
	for _, r := range ListRooms {
		for _, u := range r.users {
			emsg += fmt.Sprintf("%s - (%s)\n", u.nick, r.name)
			cnt++
		}
	}
	emsg += fmt.Sprintf("\n%d total users", cnt)
	cmd.user.DirectMessage(emsg)
	return "", nil
}

func whoHelper(cmd *ParsedCommand) (string, error) {
	msg := cmd.command + " - shows a list of all current users."
	cmd.user.DirectMessage(msg)
	return "", nil
}
