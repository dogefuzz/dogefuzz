package interfaces

import "github.com/dogefuzz/dogefuzz/pkg/common"

type SolidityCompiler interface {
	CompileSource(source string) (*common.Contract, error)
}
