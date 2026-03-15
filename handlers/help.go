package handlers

const helpText = `[vox imperialis]
commands:
  help                    show this help
  status                  uptime, load, memory, disk
  sensors                 hardware sensor readings
  service status  <name>  show service status
  service start   <name>  start a service
  service stop    <name>  stop a service
  service restart <name>  restart a service`

// Help returns the list of available commands.
func Help(cmd Command) string {
	return helpText
}
