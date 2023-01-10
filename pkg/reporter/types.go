package reporter

type ReporterType string

const (
	CONSOLE_REPORTER ReporterType = "console"
	CALLBACK_REPOTER ReporterType = "callback"
)
