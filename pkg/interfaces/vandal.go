package interfaces

import (
	"context"

	"github.com/dogefuzz/dogefuzz/pkg/common"
)

type VandalClient interface {
	Decompile(ctx context.Context, source string, name string) ([]common.Block, []common.Function, error)
}
