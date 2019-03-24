package chatroom

import (
	"fmt"
	"strings"
)

func helpHandler(cmd *ParsedCommand) (string, error) {
	// See if they want alias information also
	includeAliases := false
	for _, arg := range cmd.args {
		switch arg {
		case "-a":
			includeAliases = true
		case "-h", "-?":
			return helpHelper(cmd)
		}
	}
	msg := "Available commands:\n"
	// Create a map for the unique list of commands, since commands includes all the aliases.
	m := make(map[string]bool)
	for _, com := range commands {
		if _, ok := m[com.name]; !ok {
			m[com.name] = true
			msg += fmt.Sprintf("%s - %s", com.name, com.description)
			if includeAliases {
				msg += " - aliases: (" + strings.Join(com.aliases, ", ") + ")"
			}
			msg += "\n"
		}
	}
	msg += "Most commands also accept a -? (eg. /help -?) to list any additional options they may have."
	cmd.user.DirectMessage(msg)
	return "", nil
}

func helpHelper(cmd *ParsedCommand) (string, error) {
	msg := "/help lists all available commands. Include -a to list all the aliases for each command."
	cmd.user.DirectMessage(msg)
	return "", nil
}
