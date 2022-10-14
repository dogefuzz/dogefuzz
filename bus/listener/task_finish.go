package listener

import (
	"github.com/gongbell/contractfuzzer/bus/topic"
)

type TaskFinishEventListener interface {
	StartListening()
}

type DefaultTaskFinishEventListener struct {
	fuzzingMasterFinishChannel chan string
	taskFinishTopic            topic.TaskFinishTopic
}

func (l DefaultTaskFinishEventListener) Init(
	fuzzingMasterFinishChannel chan string,
	taskFinishTopic topic.TaskFinishTopic,
) DefaultTaskFinishEventListener {
	l.fuzzingMasterFinishChannel = fuzzingMasterFinishChannel
	l.taskFinishTopic = taskFinishTopic

	return l
}

func (l DefaultTaskFinishEventListener) StartListening() {
	l.taskFinishTopic.Subscribe(l.processEvent)
}

func (l DefaultTaskFinishEventListener) processEvent() {

}
