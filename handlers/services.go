package handlers

import (
	"fmt"
	"log"
	"strings"

	"vox-imperialis/system"
)

// NewServiceHandler returns a HandlerFunc that manages systemd services.
// Only service names present in allowedServices are permitted.
func NewServiceHandler(allowedServices []string) HandlerFunc {
	allowed := make(map[string]struct{}, len(allowedServices))
	for _, s := range allowedServices {
		allowed[strings.ToLower(strings.TrimSpace(s))] = struct{}{}
	}

	return func(cmd Command) string {
		if len(cmd.Args) < 2 {
			return "[error]\nusage: service <action> <name>\nactions: status, start, stop, restart"
		}

		action := strings.ToLower(cmd.Args[0])
		service := strings.ToLower(cmd.Args[1])

		if _, ok := allowed[service]; !ok {
			return fmt.Sprintf("[error]\nservice '%s' is not in the allowed list", service)
		}

		log.Printf("service action=%s name=%s sender=%s", action, service, cmd.Sender)

		var (
			output string
			err    error
		)

		switch action {
		case "status":
			output, err = system.SystemdStatus(service)
		case "start":
			output, err = system.SystemdStart(service)
		case "stop":
			output, err = system.SystemdStop(service)
		case "restart":
			output, err = system.SystemdRestart(service)
		default:
			return fmt.Sprintf(
				"[error]\nunknown action '%s'\nvalid actions: status, start, stop, restart",
				action,
			)
		}

		result := "success"
		if err != nil {
			result = "failed: " + err.Error()
		}

		resp := fmt.Sprintf("[service]\nname:   %s\naction: %s\nresult: %s", service, action, result)
		if output != "" {
			resp += "\n\n" + output
		}
		return resp
	}
}
