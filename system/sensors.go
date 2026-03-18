package system

import (
	"os/exec"
	"strings"
)

// GetSensors runs `sensors` (lm-sensors) and returns the trimmed output.
func GetSensors() (string, error) {
	cmd := exec.Command("sensors")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
