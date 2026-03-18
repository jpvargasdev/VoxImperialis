package handlers

import (
	"fmt"

	"vox-imperialis/system"
)

// Sensors returns hardware sensor readings via the `sensors` command.
func Sensors(cmd Command) string {
	out, err := system.GetSensors()
	if err != nil {
		return fmt.Sprintf("[error]\nsensors unavailable: %v", err)
	}
	if out == "" {
		return "[sensors]\nno data available"
	}
	return "[sensors]\n" + out
}
