package teaagent

import (
	"encoding/json"
	"time"
)

// 日志类型
type ProcessEventType = string

const (
	ProcessEventStart  = "start"
	ProcessEventStdout = "stdout"
	ProcessEventStderr = "stderr"
	ProcessEventStop   = "stop"
)

// 进程日志
type ProcessEvent struct {
	AgentId   string           `json:"agentId"`
	TaskId    string           `json:"taskId"`
	UniqueId  string           `json:"uniqueId"`
	Pid       int              `json:"pid"`
	EventType ProcessEventType `json:"eventType"`
	Data      string           `json:"data"`
	Timestamp int64            `json:"timestamp"`
}

// 新日志对象
func NewProcessEvent(eventType ProcessEventType, taskId string, uniqueId string, pid int, data []byte) *ProcessEvent {
	return &ProcessEvent{
		AgentId:   runningAgent.Id,
		TaskId:    taskId,
		UniqueId:  uniqueId,
		Pid:       pid,
		EventType: eventType,
		Data:      string(data),
		Timestamp: time.Now().Unix(),
	}
}

// 转换为JSON
func (this *ProcessEvent) AsJSON() ([]byte, error) {
	return json.Marshal(this)
}
