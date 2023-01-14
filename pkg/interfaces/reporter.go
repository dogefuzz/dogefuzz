package interfaces

import (
	"context"

	"github.com/dogefuzz/dogefuzz/pkg/common"
)

type Reporter interface {
	SendOutput(ctx context.Context, report common.TaskReport) error
}
