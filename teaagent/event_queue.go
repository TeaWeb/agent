package teaagent

var eventQueue = make(chan *ProcessEvent, 1024)

func PushEvent(event *ProcessEvent) {
	eventQueue <- event
}
