package generators

import (
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/dogefuzz/dogefuzz/pkg/common"
)

func TaskReportGen() common.TaskReport {
	transactionsCount := gofakeit.Number(1, 10)
	transactions := make([]common.TransactionReport, transactionsCount)
	for idx := 0; idx < int(transactionsCount); idx++ {
		transactions[idx] = TransactionReportGen()
	}

	weaknessesCount := gofakeit.Number(1, 10)
	weaknesses := make([]string, weaknessesCount)
	for idx := 0; idx < int(weaknessesCount); idx++ {
		weaknesses[idx] = gofakeit.LetterN(255)
	}

	return common.TaskReport{
		TimeElapsed:        time.Duration(gofakeit.Number(1, 60)) * time.Minute,
		ContractName:       gofakeit.Name(),
		Coverage:           gofakeit.Uint64(),
		CoverageByTime:     TimeSeriesDataGen(),
		MinDistance:        gofakeit.Uint64(),
		MinDistanceByTime:  TimeSeriesDataGen(),
		Transactions:       transactions,
		DetectedWeaknesses: weaknesses,
	}
}

func TimeSeriesDataGen() common.TimeSeriesData {
	seriesLength := gofakeit.Number(1, 10)
	start := gofakeit.Date()

	xs := make([]time.Time, seriesLength)
	ys := make([]uint64, seriesLength)
	for idx := 0; idx < int(seriesLength); idx++ {
		xs[idx] = start.Add(time.Duration(idx) * time.Minute)
		ys[idx] = gofakeit.Uint64()
	}
	return common.TimeSeriesData{X: xs, Y: ys}
}

func TransactionReportGen() common.TransactionReport {
	inputsCount := gofakeit.Number(1, 10)
	inputs := make([]string, inputsCount)
	for idx := 0; idx < int(inputsCount); idx++ {
		inputs[idx] = gofakeit.LetterN(255)
	}

	weaknessesCount := gofakeit.Number(1, 10)
	weaknesses := make([]string, weaknessesCount)
	for idx := 0; idx < int(weaknessesCount); idx++ {
		weaknesses[idx] = gofakeit.LetterN(255)
	}

	instructionsCount := gofakeit.Number(1, 10)
	instructions := make([]string, instructionsCount)
	for idx := 0; idx < int(instructionsCount); idx++ {
		instructions[idx] = gofakeit.LetterN(255)
	}

	return common.TransactionReport{
		Timestamp:            gofakeit.Date(),
		BlockchainHash:       gofakeit.LetterN(255),
		Inputs:               inputs,
		DetectedWeaknesses:   weaknesses,
		ExecutedInstructions: instructions,
		DeltaCoverage:        gofakeit.Uint64(),
		DeltaMinDistance:     gofakeit.Uint64(),
	}
}
