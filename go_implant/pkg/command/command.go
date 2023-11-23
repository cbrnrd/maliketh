package command

// CmdOut executes a given command and returns its output.
func CmdOut(command string) (string, error) {
	return cmdOut(command)
}
