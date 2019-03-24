package chatroom

func exitHandler(cmd *ParsedCommand) (string, error) {
	for _, arg := range cmd.args {
		switch arg {
		case "-h", "-?":
			return exitHelper(cmd)
		}
	}
	cmd.user.DirectMessage("Ok, bye!")
	return Cleanup, nil
}

func exitHelper(cmd *ParsedCommand) (string, error) {
	msg := "/exit exits the telnet-chat session."
	cmd.user.DirectMessage(msg)
	return "", nil
}
