package teaagent

import (
	"encoding/json"
	"github.com/TeaWeb/code/teaconfigs/agents"
	"time"
)

type SystemAppsEvent struct {
	Event     string              `json:"event"`
	Timestamp int64               `json:"timestamp"`
	Apps      []*agents.AppConfig `json:"apps"`
}

func NewSystemAppsEvent() *SystemAppsEvent {
	return &SystemAppsEvent{
		Event:     "SystemAppsEvent",
		Timestamp: time.Now().Unix(),
	}
}

func (this *SystemAppsEvent) AsJSON() (data []byte, err error) {
	return json.Marshal(this)
}
