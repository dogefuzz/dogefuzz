package generators

import (
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/dogefuzz/dogefuzz/pkg/common"
)

func TaskReportGen() common.TaskReport {

	weaknessesCount := gofakeit.Number(1, 10)
	weaknesses := make([]string, weaknessesCount)
	for idx := 0; idx < int(weaknessesCount); idx++ {
		weaknesses[idx] = gofakeit.LetterN(255)
	}

	return common.TaskReport{
		TaskId:             gofakeit.UUID(),
		TimeElapsed:        time.Duration(gofakeit.Number(1, 60)) * time.Minute,
		ContractName:       gofakeit.Name(),
		TotalInstructions:  gofakeit.Uint64(),
		Coverage:           uint64(gofakeit.Uint32()),
		CoverageByTime:     TimeSeriesDataGen(),
		MinDistance:        gofakeit.Uint64(),
		MinDistanceByTime:  TimeSeriesDataGen(),
		DetectedWeaknesses: weaknesses,
		Instructions: make(map[string]string),
		InstructionHitsHeatMap: make(map[string]uint64),
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
