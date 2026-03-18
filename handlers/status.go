package handlers

import (
	"fmt"

	"vox-imperialis/system"
)

// Status reports system uptime, load, memory, and root disk usage.
func Status(cmd Command) string {
	s, err := system.GetSystemStatus()
	if err != nil {
		return fmt.Sprintf("[error]\nfailed to get system status: %v", err)
	}
	return fmt.Sprintf(
		"[status]\nhost:   %s\nuptime: %s\nload:   %s\nmemory: %s\ndisk /:  %s",
		s.Hostname, s.Uptime, s.Load, s.Memory, s.Disk,
	)
}
