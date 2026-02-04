package model

import "strings"

type Opcode int

const (
	OpCmd            Opcode = 0x01
	OpSelfDestruct   Opcode = 0x02
	OpSysInfo        Opcode = 0x03
	OpSleep          Opcode = 0x04
	OpUpdateConfig   Opcode = 0x05
	OpDownload       Opcode = 0x06
	OpUpload         Opcode = 0x07
	OpInject         Opcode = 0x08
	OpChdir          Opcode = 0x09
	OpPwd            Opcode = 0x0A
	OpGetEnv         Opcode = 0x0B
	OpLs             Opcode = 0x0C
	OpPs             Opcode = 0x0D
	OpWhoami         Opcode = 0x0E
	OpDisableDefender Opcode = 0x0F
	OpClipboardGet   Opcode = 0x10
	OpClipboardSet   Opcode = 0x11
)

var opcodeNames = map[Opcode]string{
	OpCmd:            "CMD",
	OpSelfDestruct:   "SELFDESTRUCT",
	OpSysInfo:        "SYSINFO",
	OpSleep:          "SLEEP",
	OpUpdateConfig:   "UPDATE_CONFIG",
	OpDownload:       "DOWNLOAD",
	OpUpload:         "UPLOAD",
	OpInject:         "INJECT",
	OpChdir:          "CHDIR",
	OpPwd:            "PWD",
	OpGetEnv:         "GETENV",
	OpLs:             "LS",
	OpPs:             "PS",
	OpWhoami:         "WHOAMI",
	OpDisableDefender: "DISABLE_DEFENDER",
	OpClipboardGet:   "CLIPBOARD_GET",
	OpClipboardSet:   "CLIPBOARD_SET",
}

var nameToOpcode = map[string]Opcode{
	"CMD":             OpCmd,
	"SELFDESTRUCT":    OpSelfDestruct,
	"SYSINFO":         OpSysInfo,
	"SLEEP":           OpSleep,
	"UPDATE_CONFIG":   OpUpdateConfig,
	"DOWNLOAD":        OpDownload,
	"UPLOAD":          OpUpload,
	"INJECT":          OpInject,
	"CHDIR":           OpChdir,
	"PWD":             OpPwd,
	"GETENV":          OpGetEnv,
	"LS":              OpLs,
	"PS":              OpPs,
	"WHOAMI":          OpWhoami,
	"DISABLE_DEFENDER": OpDisableDefender,
	"CLIPBOARD_GET":   OpClipboardGet,
	"CLIPBOARD_SET":   OpClipboardSet,
}

func (o Opcode) String() string {
	if name, ok := opcodeNames[o]; ok {
		return name
	}
	return "UNKNOWN"
}

func OpcodeFromName(name string) (Opcode, bool) {
	op, ok := nameToOpcode[strings.ToUpper(name)]
	return op, ok
}
