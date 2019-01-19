package teaagent

import (
	"encoding/json"
	"time"
)

// 监控项事件
type ItemEvent struct {
	Event     string      `json:"event"`
	AgentId   string      `json:"agentId"`
	AppId     string      `json:"appId"`
	ItemId    string      `json:"itemId"`
	Value     interface{} `json:"value"`
	Timestamp int64       `json:"timestamp"`
}

// 获取新监控项事件
func NewItemEvent(agentId string, appId string, itemId string, value interface{}) *ItemEvent {
	return &ItemEvent{
		Event:     "ItemEvent",
		AgentId:   agentId,
		AppId:     appId,
		ItemId:    itemId,
		Value:     value,
		Timestamp: time.Now().Unix(),
	}
}

// 转换为JSON
func (this *ItemEvent) AsJSON() ([]byte, error) {
	return json.Marshal(this)
}
