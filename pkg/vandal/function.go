package vandal

import "github.com/dogefuzz/dogefuzz/pkg/common"

func NewFunctionFromLines(lines []string) common.Function {
	function := common.Function{}
	function.Signature = readStringPropertyLine("Public function signature", lines[0])
	function.EntryBlock = readStringPropertyLine("Entry block", lines[1])
	function.ExitBlock = readStringPropertyLine("Exit block", lines[2])
	function.Body = readSlicePropertyLine("Body", lines[3])
	return function
}
