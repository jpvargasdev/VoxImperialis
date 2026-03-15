package main

import "fmt"

// formatUnknown returns a user-facing error for an unrecognised command.
func formatUnknown(name string) string {
	if name == "" {
		return "[error]\nempty command — type 'help' for available commands"
	}
	return fmt.Sprintf("[error]\nunknown command '%s' — type 'help' for available commands", name)
}
