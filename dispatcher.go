package main

import (
	"log"

	"vox-imperialis/handlers"
)

// Dispatcher routes incoming commands to registered handlers.
type Dispatcher struct {
	registry map[string]handlers.HandlerFunc
}

// NewDispatcher creates an empty Dispatcher.
func NewDispatcher() *Dispatcher {
	return &Dispatcher{registry: make(map[string]handlers.HandlerFunc)}
}

// Register adds a handler for the given command name.
func (d *Dispatcher) Register(name string, h handlers.HandlerFunc) {
	d.registry[name] = h
}

// Dispatch looks up the handler for cmd.Name and calls it.
// Returns a formatted error string for unknown or empty commands.
func (d *Dispatcher) Dispatch(cmd handlers.Command) string {
	log.Printf("dispatch: from=%s cmd=%q args=%v", cmd.Sender, cmd.Name, cmd.Args)
	h, ok := d.registry[cmd.Name]
	if !ok {
		return formatUnknown(cmd.Name)
	}
	return h(cmd)
}
