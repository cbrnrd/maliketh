package builder

import (
	"testing"
)

func TestBuildOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    BuildOptions
		wantErr bool
	}{
		{
			name:    "default options are valid",
			opts:    DefaultBuildOptions,
			wantErr: false,
		},
		{
			name: "empty options get defaults",
			opts: BuildOptions{},
			wantErr: false,
		},
		{
			name: "invalid implant type",
			opts: BuildOptions{
				ImplantType: "invalid",
			},
			wantErr: true,
		},
		{
			name: "cpp implant with linux target fails",
			opts: BuildOptions{
				ImplantType: ImplantTypeCpp,
				TargetOS:    TargetLinux,
			},
			wantErr: true,
		},
		{
			name: "go implant with linux target succeeds",
			opts: BuildOptions{
				ImplantType: ImplantTypeGo,
				TargetOS:    TargetLinux,
				TargetArch:  ArchAMD64,
			},
			wantErr: false,
		},
		{
			name: "invalid target os",
			opts: BuildOptions{
				ImplantType: ImplantTypeGo,
				TargetOS:    "freebsd",
			},
			wantErr: true,
		},
		{
			name: "invalid target arch",
			opts: BuildOptions{
				ImplantType: ImplantTypeGo,
				TargetArch:  "mips",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBuildOptions_ToCppCompilerFlags(t *testing.T) {
	opts := BuildOptions{
		InitialSleepSeconds: 300,
		SchtaskPersist:      true,
		UseAntidebug:        false,
		KillParent:          true,
		UseAntiVM:           false,
		ScheduledTaskName:   "TestTask",
		RegisterMaxRetries:  10,
	}

	flags := opts.ToCppCompilerFlags()

	expected := []string{
		"-DINITIAL_SLEEP_SECONDS=300",
		"-DSCHTASK_PERSIST=TRUE",
		"-DUSE_ANTIDEBUG=FALSE",
		"-DKILL_PARENT=TRUE",
		"-DUSE_ANTIVM=FALSE",
		"-DSCHEDULED_TASK_NAME='OBFUSCATED(\\\"TestTask\\\")'",
		"-DREGISTER_MAX_RETRIES=10",
	}

	if len(flags) != len(expected) {
		t.Errorf("expected %d flags, got %d", len(expected), len(flags))
	}

	for i, flag := range flags {
		if flag != expected[i] {
			t.Errorf("flag[%d] = %q, want %q", i, flag, expected[i])
		}
	}
}

func TestBuildOptions_ToGoLDFlags(t *testing.T) {
	tests := []struct {
		name    string
		debug   bool
		want    string
	}{
		{
			name:  "release mode",
			debug: false,
			want:  "-s -w -X maliketh.config.DEBUG=false",
		},
		{
			name:  "debug mode",
			debug: true,
			want:  "-X maliketh.config.DEBUG=true -X maliketh.config.INITIAL_SLEEP=1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := BuildOptions{Debug: tt.debug}
			got := opts.ToGoLDFlags()
			if got != tt.want {
				t.Errorf("ToGoLDFlags() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCleanupString(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"TestTask", "TestTask"},
		{"Test Task", "TestTask"},
		{"Test-Task_123", "TestTask123"},
		{"Test@Task#$%^", "TestTask"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := cleanupString(tt.input)
			if got != tt.want {
				t.Errorf("cleanupString(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestFromMap(t *testing.T) {
	tests := []struct {
		name    string
		input   map[string]any
		check   func(*BuildOptions) bool
		wantErr bool
	}{
		{
			name:  "empty map uses defaults",
			input: map[string]any{},
			check: func(o *BuildOptions) bool {
				return o.ImplantType == ImplantTypeGo &&
					o.TargetOS == TargetWindows &&
					o.TargetArch == ArchAMD64
			},
			wantErr: false,
		},
		{
			name: "parses implant type",
			input: map[string]any{
				"implant_type": "cpp",
			},
			check: func(o *BuildOptions) bool {
				return o.ImplantType == ImplantTypeCpp
			},
			wantErr: false,
		},
		{
			name: "parses go implant for linux",
			input: map[string]any{
				"implant_type": "go",
				"target_os":    "linux",
				"target_arch":  "arm64",
			},
			check: func(o *BuildOptions) bool {
				return o.ImplantType == ImplantTypeGo &&
					o.TargetOS == TargetLinux &&
					o.TargetArch == ArchARM64
			},
			wantErr: false,
		},
		{
			name: "parses boolean options",
			input: map[string]any{
				"use_antidebug": false,
				"use_antivm":    false,
				"debug":         true,
			},
			check: func(o *BuildOptions) bool {
				return !o.UseAntidebug && !o.UseAntiVM && o.Debug
			},
			wantErr: false,
		},
		{
			name: "parses numeric options",
			input: map[string]any{
				"initial_sleep_seconds": float64(600),
				"register_max_retries":  float64(20),
			},
			check: func(o *BuildOptions) bool {
				return o.InitialSleepSeconds == 600 && o.RegisterMaxRetries == 20
			},
			wantErr: false,
		},
		{
			name: "invalid implant type fails",
			input: map[string]any{
				"implant_type": "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts, err := FromMap(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil && !tt.check(opts) {
				t.Errorf("FromMap() check failed for %+v", opts)
			}
		})
	}
}

func TestBoolToUpperString(t *testing.T) {
	if boolToUpperString(true) != "TRUE" {
		t.Error("expected TRUE for true")
	}
	if boolToUpperString(false) != "FALSE" {
		t.Error("expected FALSE for false")
	}
}

func TestNew(t *testing.T) {
	// Test with nil config
	b := New(nil, nil)
	if b.config == nil {
		t.Error("expected default config")
	}
	if b.config.BuildTimeout != DefaultBuildTimeout {
		t.Error("expected default build timeout")
	}

	// Test with custom config
	cfg := &Config{
		CppImplantDir: "/custom/cpp",
		GoImplantDir:  "/custom/go",
	}
	b = New(cfg, nil)
	if b.config.CppImplantDir != "/custom/cpp" {
		t.Error("expected custom cpp dir")
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.CppImplantDir != "/implant" {
		t.Errorf("expected /implant, got %s", cfg.CppImplantDir)
	}
	if cfg.GoImplantDir != "/go_implant" {
		t.Errorf("expected /go_implant, got %s", cfg.GoImplantDir)
	}
	if cfg.BuildTimeout != DefaultBuildTimeout {
		t.Errorf("expected %v, got %v", DefaultBuildTimeout, cfg.BuildTimeout)
	}
}
