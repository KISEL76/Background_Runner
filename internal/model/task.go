package model

type Status string

const (
	StatusQueued  Status = "queued"
	StatusRunning Status = "running"
	StatusDone    Status = "done"
	StatusFailed  Status = "failed"
)

type Task struct {
	ID         string `json:"id"`
	Payload    string `json:"payload"`
	MaxRetries int    `json:"max_retries"`

	Attempts int    `json:"-"`
	Status   Status `json:"-"`
	Err      string `json:"-"`
}
