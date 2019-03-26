package chatroom

import "fmt"

func changeNickHandler(cmd *ParsedCommand) (string, error) {
	ok, emsg, err := checkArgs(cmd, changeNickHelper, 1, 1)
	if !ok {
		return emsg, err
	}
	if existingNick(cmd.args[0]) {
		emsg = fmt.Sprintf("Nickname %s already taken", cmd.args[0])
	} else {
		removeNick(cmd.user.nick)
		addNick(cmd.args[0])
		cmd.user.nick = cmd.args[0]
		emsg = fmt.Sprintf("Changed nickname to [%s]", cmd.args[0])
	}
	cmd.user.DirectMessage(emsg)
	return "", nil
}

func changeNickHelper(cmd *ParsedCommand) (string, error) {
	msg := cmd.command + " <new name> - allows you to change your nickname."
	cmd.user.DirectMessage(msg)
	return "", nil
}
