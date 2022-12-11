package listener

type TaskStartEventListener interface {
	StartListening()
}

// type DefaultTaskStartEventListener struct {
// 	fuzzingMasterRequestChannel chan string
// 	taskFinishTopic             topic.TaskStartTopic
// }

// func (l DefaultTaskStartEventListener) Init(
// 	fuzzingMasterRequestChannel chan string,
// 	taskFinishTopic topic.TaskStartTopic,
// ) DefaultTaskStartEventListener {
// 	l.fuzzingMasterRequestChannel = fuzzingMasterRequestChannel
// 	l.taskFinishTopic = taskFinishTopic

// 	return l
// }

// func (l DefaultTaskStartEventListener) StartListening() {
// 	l.taskFinishTopic.Subscribe(l.processEvent)
// }

// func (l DefaultTaskStartEventListener) processEvent() {

// }
