package sandbox

import "os"

func sandboxFilepath() bool {
	return false
}

func sandboxDisk(size int) bool {
	return false
}

func sandboxTmp(entries int) bool {
	tmp_dir := "/tmp"
	files, err := os.ReadDir(tmp_dir)
	if err != nil {
		return true
	}

	return len(files) < entries
}
