package reporter

type ReporterType string

const (
	CONSOLE_REPORTER ReporterType = "console"
	WEBHOOK_REPOTER  ReporterType = "webhook"
	FILE_REPORTER    ReporterType = "file"
)
