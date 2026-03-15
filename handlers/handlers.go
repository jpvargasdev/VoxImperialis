package handlers

// Command represents a parsed user command.
type Command struct {
	Name   string
	Args   []string
	Sender string
}

// HandlerFunc processes a command and returns a plain-text response.
type HandlerFunc func(cmd Command) string
