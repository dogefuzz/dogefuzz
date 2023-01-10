package reporter

import (
	"context"
	"fmt"

	"github.com/dogefuzz/dogefuzz/pkg/common"
)

type consoleReporter struct {
}

func NewConsoleReporter() *consoleReporter {
	return &consoleReporter{}
}

func (r *consoleReporter) SendOutput(ctx context.Context, report common.TaskReport) error {
	fmt.Println("********** FUZZING EXECUTION RESULT **********")
	fmt.Printf("Time Elapsed: %v\n", report.TimeElapsed)
	fmt.Printf("Contract Name: %s\n", report.ContractName)
	fmt.Printf("Coverage: %d\n", report.Coverage)
	fmt.Printf("Min Distance: %d\n", report.MinDistance)
	fmt.Printf("Transactions: %d\n", len(report.Transactions))
	fmt.Printf("Weakneses Found:\n")
	for _, weakness := range report.DetectedWeaknesses {
		fmt.Printf("\t- %s\n", weakness)
	}
	return nil
}
