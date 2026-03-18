package system

import (
	"fmt"
	"os/exec"
	"strings"
)

// SystemdStatus returns the output of `systemctl status <service>`.
func SystemdStatus(service string) (string, error) {
	return runSystemctl("status", service)
}

// SystemdRestart restarts the named service via systemctl.
func SystemdRestart(service string) (string, error) {
	return runSystemctl("restart", service)
}

// SystemdStart starts the named service via systemctl.
func SystemdStart(service string) (string, error) {
	return runSystemctl("start", service)
}

// SystemdStop stops the named service via systemctl.
func SystemdStop(service string) (string, error) {
	return runSystemctl("stop", service)
}

// runSystemctl executes systemctl with the given action and service name.
// arguments are passed as discrete values — no shell expansion occurs.
func runSystemctl(action, service string) (string, error) {
	cmd := exec.Command("systemctl", action, service)
	out, err := cmd.CombinedOutput()
	output := strings.TrimSpace(string(out))
	if err != nil {
		return output, fmt.Errorf("systemctl %s %s: %w", action, service, err)
	}
	return output, nil
}
