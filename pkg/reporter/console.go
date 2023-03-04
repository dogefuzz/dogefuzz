package reporter

import (
	"context"
	"fmt"
	"io"

	"github.com/dogefuzz/dogefuzz/pkg/common"
)

type consoleReporter struct {
	writer io.Writer
}

func NewConsoleReporter(writer io.Writer) *consoleReporter {
	return &consoleReporter{writer: writer}
}

func (r *consoleReporter) SendOutput(ctx context.Context, report common.TaskReport) error {
	fmt.Fprintln(r.writer, "********** FUZZING EXECUTION RESULT **********")
	fmt.Fprintf(r.writer, "Time Elapsed: %v\n", report.TimeElapsed)
	fmt.Fprintf(r.writer, "Contract Name: %s\n", report.ContractName)
	fmt.Fprintf(r.writer, "Coverage: %.2f%%\n", 100*(float64(report.Coverage)/float64(report.TotalInstructions)))
	fmt.Fprintf(r.writer, "Min Distance: %d\n", report.MinDistance)
	fmt.Fprintf(r.writer, "Transactions: %d\n", len(report.Transactions))
	fmt.Fprintln(r.writer, "Weakneses Found:")
	if len(report.DetectedWeaknesses) == 0 {
		fmt.Fprintln(r.writer, "None")
	}
	for _, weakness := range report.DetectedWeaknesses {
		fmt.Fprintf(r.writer, "\t- %s\n", weakness)
	}
	return nil
}
