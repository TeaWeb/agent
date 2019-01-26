package teaagent

import (
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/timers"
	"time"
)

// 监控项定义
type Item struct {
	appId     string
	config    *agents.Item
	lastTimer *time.Ticker
}

// 获取新任务
func NewItem(appId string, config *agents.Item) *Item {
	return &Item{
		appId:  appId,
		config: config,
	}
}

// 运行一次
func (this *Item) Run() (value interface{}, err error) {
	return this.config.Source().Execute(nil)
}

// 定时运行
func (this *Item) Schedule() {
	if this.lastTimer != nil {
		this.lastTimer.Stop()
	}
	this.lastTimer = timers.Every(this.config.IntervalDuration(), func(ticker *time.Ticker) {
		value, err := this.config.Source().Execute(nil)
		if err != nil {
			logs.Println("error:" + err.Error())
		}

		if value != nil {
			PushEvent(NewItemEvent(runningAgent.Id, this.appId, this.config.Id, value, err))
		} else {
			PushEvent(NewItemEvent(runningAgent.Id, this.appId, this.config.Id, "", err))
		}
	})
}

// 停止运行
func (this *Item) Stop() {
	if this.lastTimer != nil {
		this.lastTimer.Stop()
	}
}
