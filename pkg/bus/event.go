package bus

type TaskStartEvent struct {
	TaskId string
}

type TaskFinishEvent struct {
	TaskId string
}

type TaskInputRequestEvent struct {
	TaskId string
}

type InstrumentExecutionEvent struct {
	TransactionId string
}
