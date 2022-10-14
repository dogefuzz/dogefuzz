package model

type InstrumentExecutionRequest struct {
	Name         string   `json:"name"`
	Input        string   `json:"input"`
	Instructions []uint64 `json:"instructions"`
	TxHash       string   `json:"txHash"`
}

type InstrumentWeaknessRequest struct {
	OracleEvents []OracleEvent `json:"oracleEvents"`
	Execution    Execution     `json:"execution"`
	TxHash       string        `json:"txHash"`
}

type OracleEvent string

// type Profile string

type Execution struct {
	Metadata    ExecutionMetadata `json:"metadata"`
	CallsLength int               `json:"callsLength"`
	Trace       ExecutionTrace    `json:"trace"`
}

type ExecutionMetadata struct {
	Caller string `json:"caller"`
	Callee string `json:"callee"`
	Value  string `json:"value"`
	Gas    string `json:"gas"`
	Input  string `json:"input"`
}

type ExecutionTrace []string
