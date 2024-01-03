//go:build debug
// +build debug

// Use this file if in debug mode (ie if `-tags debug` is set at in `go build`)

package startup

import "maliketh/pkg/sandbox"

func doStartup() {

	// Just print the sandbox status
	sandbox.SandboxAll()
}
