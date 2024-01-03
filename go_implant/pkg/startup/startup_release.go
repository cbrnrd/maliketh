//go:build !debug
// +build !debug

// Use this file if not in debug mode

package startup

import (
	"maliketh/pkg/config"
	"maliketh/pkg/sandbox"
	"os"
)

func doStartup() {

	// Do our initial sleep, exiting if we're being sleep skipped
	sandbox.SleepNExitIfSandboxed(config.INITIAL_SLEEP)

	if sandbox.SandboxAll() {
		os.Exit(0)
	}

}
