package reporter

import (
	"context"
	"encoding/json"
	"os"

	"github.com/dogefuzz/dogefuzz/pkg/common"
)

type fileReporter struct {
	outputFile string
}

func NewFileReporter(outputFile string) *fileReporter {
	return &fileReporter{outputFile: outputFile}
}

func (r *fileReporter) SendOutput(ctx context.Context, report common.TaskReport) error {
	if fileExists(r.outputFile) {
		err := os.Remove(r.outputFile)
		if err != nil {
			return err
		}
	}

	jsonContent, err := json.MarshalIndent(report, "", "\t")
	if err != nil {
		return err
	}

	return os.WriteFile(r.outputFile, jsonContent, 0644)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
