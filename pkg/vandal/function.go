package vandal

type Function struct {
	Signature  string   `json:"signature"`
	EntryBlock string   `json:"entryBlock"`
	ExitBlock  string   `json:"exitBlock"`
	Body       []string `json:"body"`
}

func NewFunctionFromLines(lines []string) Function {
	function := Function{}
	function.Signature = readStringPropertyLine("Public function signature", lines[0])
	function.EntryBlock = readStringPropertyLine("Entry block", lines[1])
	function.ExitBlock = readStringPropertyLine("Exit block", lines[2])
	function.Body = readSlicePropertyLine("Body", lines[3])
	return function
}
