package model

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"time"
)

type TaskStatus string

const (
	TaskStatusCreated  TaskStatus = "CREATED"
	TaskStatusTasked   TaskStatus = "TASKED"
	TaskStatusComplete TaskStatus = "COMPLETE"
	TaskStatusError    TaskStatus = "ERROR"
)

type OperatorRole string

const (
	RoleAdmin    OperatorRole = "admin"
	RoleOperator OperatorRole = "operator"
)

type Operator struct {
	ID              int64        `json:"id"`
	Username        string       `json:"username"`
	PublicKey       string       `json:"public_key"`
	VerifyKey       string       `json:"verify_key"`
	LoginSecret     string       `json:"login_secret"`
	AuthToken       *string      `json:"auth_token,omitempty"`
	AuthTokenExpiry *time.Time   `json:"auth_token_expiry,omitempty"`
	CreatedAt       time.Time    `json:"created_at"`
	LastLogin       *time.Time   `json:"last_login,omitempty"`
	RMQQueue        string       `json:"rmq_queue"`
	Role            OperatorRole `json:"role"`
	Revoked         bool         `json:"revoked"`
}

type Implant struct {
	ID        int64     `json:"id"`
	ImplantID string    `json:"implant_id"`
	Hostname  string    `json:"hostname"`
	IP        string    `json:"ip"`
	OS        string    `json:"os"`
	Arch      string    `json:"arch"`
	Username  string    `json:"user"`
	ServerSK  string    `json:"-"`
	ImplantPK string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	LastSeen  time.Time `json:"last_seen"`
}

type ImplantAlias struct {
	ID        int64  `json:"id"`
	ImplantID string `json:"implant_id"`
	Alias     string `json:"alias"`
}

type ImplantConfig struct {
	ID                   int64     `json:"id"`
	ImplantID            string    `json:"implant_id"`
	Cookie               string    `json:"cookie"`
	KillDate             *time.Time `json:"kill_date,omitempty"`
	UserAgent            string    `json:"user_agent"`
	AutoSelfDestruct     bool      `json:"auto_self_destruct"`
	SleepTime            int       `json:"sleep_time"`
	Jitter               float64   `json:"jitter"`
	MaxRetries           int       `json:"max_retries"`
	RetryWait            int       `json:"retry_wait"`
	RetryJitter          float64   `json:"retry_jitter"`
	EncKey               string    `json:"enc_key"`
	TailoringHashFunc    string    `json:"tailoring_hash_function"`
	TailoringHashRounds  int       `json:"tailoring_hash_rounds"`
	TailoringHashes      []string  `json:"tailoring_hashes"`
}

func (c *ImplantConfig) ToJSON() map[string]any {
	return map[string]any{
		"cookie":                  c.Cookie,
		"kill_date":               c.KillDate,
		"user_agent":              c.UserAgent,
		"auto_self_destruct":      c.AutoSelfDestruct,
		"sleep_time":              c.SleepTime,
		"jitter":                  c.Jitter,
		"max_retries":             c.MaxRetries,
		"retry_wait":              c.RetryWait,
		"retry_jitter":            c.RetryJitter,
		"enc_key":                 c.EncKey,
		"tailoring_hash_function": c.TailoringHashFunc,
		"tailoring_hash_rounds":   c.TailoringHashRounds,
		"tailoring_hashes":        c.TailoringHashes,
	}
}

type Task struct {
	ID           int64           `json:"id"`
	TaskID       string          `json:"task_id"`
	OperatorName string          `json:"operator_name"`
	ImplantID    string          `json:"implant_id"`
	Opcode       Opcode          `json:"opcode"`
	Args         json.RawMessage `json:"args"`
	Status       TaskStatus      `json:"status"`
	Output       *string         `json:"output,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
	ReadAt       *time.Time      `json:"read_at,omitempty"`
	ExecutedAt   *time.Time      `json:"executed_at,omitempty"`
}

func (t *Task) ToFilteredJSON() map[string]any {
	return map[string]any{
		"task_id": t.TaskID,
		"opcode":  t.Opcode,
		"args":    t.Args,
	}
}

func RandomID(n int) string {
	bytes := make([]byte, n)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)[:n]
}
