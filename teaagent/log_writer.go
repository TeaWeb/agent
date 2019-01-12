package teaagent

// 日志写入器
type StdoutLogWriter struct {
	TaskId   string
	UniqueId string
	Pid      int
}

func (this *StdoutLogWriter) Write(p []byte) (n int, err error) {
	event := NewProcessEvent(ProcessEventStdout, this.TaskId, this.UniqueId, this.Pid, p)
	PushEvent(event)

	n = len(p)
	return
}

type StderrLogWriter struct {
	TaskId   string
	UniqueId string
	Pid      int
}

func (this *StderrLogWriter) Write(p []byte) (n int, err error) {
	event := NewProcessEvent(ProcessEventStderr, this.TaskId, this.UniqueId, this.Pid, p)
	PushEvent(event)

	n = len(p)
	return
}
