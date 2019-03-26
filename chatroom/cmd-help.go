package chatroom

import (
	"fmt"
	"sort"
	"strings"
)

func helpHandler(cmd *ParsedCommand) (string, error) {
	// See if they want alias information also
	argCnt := 0
	includeAliases := false
	usage := helpHelper
	for _, arg := range cmd.args {
		switch arg {
		case "-a":
			includeAliases = true
			argCnt++
		case "-h", "-?":
			return usage(cmd)
		}
	}
	msg := "Available commands:\n"
	// Create a map for the unique list of commands, since commands includes all the aliases.
	m := make(map[string]bool)
	commandList := []string{}
	for _, com := range commands {
		if _, ok := m[com.name]; !ok {
			m[com.name] = true
			// Store commandList so we can sort it later.
			commandList = append(commandList, com.name)
		}
	}
	sort.Strings(commandList)
	for _, comName := range commandList {
		com := commands[comName]
		msg += fmt.Sprintf("%s - %s", com.name, com.description)
		if includeAliases && len(com.aliases) > 0 {
			msg += " - aliases: (" + strings.Join(com.aliases, ", ") + ")"
		}
		msg += "\n"
	}
	msg += "Most commands also accept a -? (eg. /help -?) to list any additional options they may have."
	cmd.user.DirectMessage(msg)
	return "", nil
}

func helpHelper(cmd *ParsedCommand) (string, error) {
	msg := cmd.command + " - lists all available commands. Include -a to list all the aliases for each command."
	cmd.user.DirectMessage(msg)
	return "", nil
}
