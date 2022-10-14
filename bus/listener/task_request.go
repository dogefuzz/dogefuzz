package listener

import "github.com/dogefuzz/dogefuzz/bus/topic"

type TaskRequestEventListener interface {
	StartListening()
}

type DefaultTaskRequestEventListener struct {
	fuzzingMasterRequestChannel chan string
	taskFinishTopic             topic.TaskRequestTopic
}

func (l DefaultTaskRequestEventListener) Init(
	fuzzingMasterRequestChannel chan string,
	taskFinishTopic topic.TaskRequestTopic,
) DefaultTaskRequestEventListener {
	l.fuzzingMasterRequestChannel = fuzzingMasterRequestChannel
	l.taskFinishTopic = taskFinishTopic

	return l
}

func (l DefaultTaskRequestEventListener) StartListening() {
	l.taskFinishTopic.Subscribe(l.processEvent)
}

func (l DefaultTaskRequestEventListener) processEvent() {

}
