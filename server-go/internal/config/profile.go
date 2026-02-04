package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type MaleableProfile struct {
	Name          string         `yaml:"name"`
	Routes        Routes         `yaml:"routes"`
	Globals       GlobalOptions  `yaml:"globals,inline"`
	ImplantConfig ImplantProfile `yaml:"client"`
	ServerConfig  ServerProfile  `yaml:"server,inline"`
}

type Routes struct {
	BasePath    string `yaml:"base_path"`
	Register    Route  `yaml:"register"`
	Checkin     Route  `yaml:"checkin"`
	TaskResults Route  `yaml:"task_results"`
}

type Route struct {
	Path    string   `yaml:"path"`
	Methods []string `yaml:"methods"`
}

type GlobalOptions struct {
	ImplantIDCookie      string `yaml:"implant_id_cookie"`
	RegistrationPassword string `yaml:"registration_password"`
}

type ImplantProfile struct {
	KillDate             string   `yaml:"kill_date"`
	UserAgent            string   `yaml:"user_agent"`
	Sleep                int      `yaml:"sleep"`
	Jitter               float64  `yaml:"jitter"`
	MaxRetries           int      `yaml:"max_retries"`
	AutoSelfDestruct     bool     `yaml:"auto_self_destruct"`
	RetryWait            int      `yaml:"retry_wait"`
	RetryJitter          float64  `yaml:"retry_jitter"`
	TailoringHashFunc    string   `yaml:"tailoring_hash_function"`
	TailoringHashRounds  int      `yaml:"tailoring_hash_rounds"`
	TailoringHashes      []string `yaml:"tailoring_hashes"`
}

type ServerProfile struct {
	Headers     map[string]string `yaml:"headers"`
	RedirectURL string            `yaml:"redirect_url"`
}

func LoadProfile(path string) (*MaleableProfile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read profile file: %w", err)
	}

	var profile MaleableProfile
	if err := yaml.Unmarshal(data, &profile); err != nil {
		return nil, fmt.Errorf("failed to parse profile YAML: %w", err)
	}

	if err := profile.Validate(); err != nil {
		return nil, fmt.Errorf("invalid profile: %w", err)
	}

	return &profile, nil
}

func (p *MaleableProfile) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("profile must have a name")
	}
	if p.Routes.BasePath == "" {
		return fmt.Errorf("routes must have a base_path")
	}
	if p.Routes.Register.Path == "" || len(p.Routes.Register.Methods) == 0 {
		return fmt.Errorf("register route must have path and methods")
	}
	if p.Routes.Checkin.Path == "" || len(p.Routes.Checkin.Methods) == 0 {
		return fmt.Errorf("checkin route must have path and methods")
	}
	if p.Routes.TaskResults.Path == "" || len(p.Routes.TaskResults.Methods) == 0 {
		return fmt.Errorf("task_results route must have path and methods")
	}
	if p.Globals.ImplantIDCookie == "" {
		p.Globals.ImplantIDCookie = "SESSID"
	}
	if p.ServerConfig.RedirectURL == "" {
		p.ServerConfig.RedirectURL = "https://www.google.com"
	}
	return nil
}

func DefaultProfile() *MaleableProfile {
	return &MaleableProfile{
		Name: "default",
		Routes: Routes{
			BasePath:    "/api",
			Register:    Route{Path: "/register", Methods: []string{"POST"}},
			Checkin:     Route{Path: "/checkin", Methods: []string{"GET"}},
			TaskResults: Route{Path: "/results", Methods: []string{"POST"}},
		},
		Globals: GlobalOptions{
			ImplantIDCookie: "SESSID",
		},
		ImplantConfig: ImplantProfile{
			Sleep:               60,
			Jitter:              0.1,
			MaxRetries:          -1,
			TailoringHashFunc:   "sha256",
			TailoringHashRounds: 1,
		},
		ServerConfig: ServerProfile{
			RedirectURL: "https://www.google.com",
			Headers:     make(map[string]string),
		},
	}
}
