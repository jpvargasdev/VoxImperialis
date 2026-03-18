package handlers

const helpText = `[vox imperialis]
commands:
  help                                show this help
  status                              uptime, load, memory, disk
  sensors                             hardware sensor readings
  service status  <name>              show service status
  service start   <name>              start a service
  service stop    <name>              stop a service
  service restart <name>              restart a service
  watchface create  <name> <prompt>   create new watchface project
  watchface list                      list all watchface projects
  watchface status  <name>            check agent progress
  watchface iterate <name> <feedback> send feedback to agent
  watchface stop    <name>            stop agent work
  watchface delete  <name>            remove watchface project`

// Help returns the list of available commands.
func Help(cmd Command) string {
	return helpText
}
