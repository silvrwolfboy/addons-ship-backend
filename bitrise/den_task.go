package bitrise

import "time"

// TaskParams ...
type TaskParams struct {
	Workflow    string      `json:"workflow_id"`
	StackID     string      `json:"stack_id"`
	BuildConfig interface{} `json:"build_config"`
	Secrets     interface{} `json:"secrets"`
	InlineEnvs  interface{} `json:"inline_envs"`
	WebhookURL  string      `json:"webhook_url"`
}

// TriggerResponse ...
type TriggerResponse struct {
	Aborted                bool       `json:"aborted"`
	Config                 string     `json:"config"`
	ConfigType             string     `json:"config_type"`
	CreatedAt              time.Time  `json:"created_at"`
	ExitCode               *int       `json:"exit_code"`
	FinishedAt             *time.Time `json:"finished_at"`
	GeneratedLogChunkCount *int       `json:"generated_log_chunk_count"`
	TaskIdentifier         string     `json:"id"`
	StartedAt              *time.Time `json:"started_at"`
	Tags                   string     `json:"tags"`
	TimedOut               bool       `json:"timed_out"`
	TimeoutSeconds         int        `json:"timeout_seconds"`
	UpdatedAt              time.Time  `json:"updated_at"`
	WebhookURL             string     `json:"webhook_url"`
}