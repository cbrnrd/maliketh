package shellcode

// Injects a bytearray into current process and executes it
func RunShellcode(sc []byte, bg bool) error {
	return runShellcode(sc, bg)
}