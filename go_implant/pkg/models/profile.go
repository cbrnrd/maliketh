package models

import "encoding/json"

type MalleableProfile struct {
	ImplantId string        `json:"id"`
	Config    ImplantConfig `json:"config"`
}

type ImplantConfig struct {
	Cookie                string  `json:"cookie"`
	UserAgent             string  `json:"user_agent"`
	Sleep                 int     `json:"sleep_time"`
	Jitter                float32 `json:"jitter"`
	MaxRetries            int     `json:"max_retries"`
	AutoSelfDstruct       bool    `json:"auto_self_destruct"`
	RetryWait             int     `json:"retry_wait"`
	RetryJitter           float32 `json:"retry_jitter"`
	TailoringHashRounds   int     `json:"tailoring_hash_rounds"`
	Base64EncryptionKey   string
	Base64ServerPublicKey string   `json:"enc_key"`
	TailoringHashFunction string   `json:"tailoring_hash_function"`
	TailoringHashes       []string `json:"tailoring_hashes"`
}

func ParseMalleableProfile(profileJson string, base64privateKey string) (MalleableProfile, error) {
	var profile MalleableProfile
	err := json.Unmarshal([]byte(profileJson), &profile)
	profile.Config.Base64EncryptionKey = base64privateKey
	return profile, err
}

// UpdateConfig updates this profile with the given values in `config`. Unknown keys are ignored.
func (c *ImplantConfig) UpdateConfig(config map[string]interface{}) {
	for k, v := range config {
		switch k {
		case "user_agent":
			c.UserAgent = v.(string)
		case "sleep_time":
			c.Sleep = v.(int)
		case "jitter":
			c.Jitter = float32(v.(float64))
		case "max_retries":
			c.MaxRetries = v.(int)
		case "auto_self_destruct":
			c.AutoSelfDstruct = v.(bool)
		case "retry_wait":
			c.RetryWait = v.(int)
		case "retry_jitter":
			c.RetryJitter = float32(v.(float64))
		case "tailoring_hash_rounds":
			c.TailoringHashRounds = v.(int)
		case "tailoring_hash_function":
			c.TailoringHashFunction = v.(string)
		case "tailoring_hashes":
			c.TailoringHashes = v.([]string)
		}
	}
}
