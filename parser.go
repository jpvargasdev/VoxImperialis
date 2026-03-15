package main

import (
	"strings"

	"vox-imperialis/handlers"
)

// ParseCommand parses a raw message text into a structured Command.
// The first word is the command name; remaining words are arguments.
func ParseCommand(text, sender string) handlers.Command {
	text = strings.TrimSpace(text)
	fields := strings.Fields(text)
	if len(fields) == 0 {
		return handlers.Command{Sender: sender}
	}
	return handlers.Command{
		Name:   strings.ToLower(fields[0]),
		Args:   fields[1:],
		Sender: sender,
	}
}
