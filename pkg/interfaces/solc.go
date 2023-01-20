package interfaces

import "github.com/dogefuzz/dogefuzz/pkg/common"

type SolidityCompiler interface {
	CompileSource(contractName string, source string) (*common.Contract, error)
}
