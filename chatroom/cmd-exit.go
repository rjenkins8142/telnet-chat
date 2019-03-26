package chatroom

func exitHandler(cmd *ParsedCommand) (string, error) {
	ok, emsg, err := checkArgs(cmd, exitHelper, 0, 0)
	if !ok {
		return emsg, err
	}
	cmd.user.DirectMessage("Ok, bye!")
	return Cleanup, nil
}

func exitHelper(cmd *ParsedCommand) (string, error) {
	msg := cmd.command + " - exits the telnet-chat session."
	cmd.user.DirectMessage(msg)
	return "", nil
}
