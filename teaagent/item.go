package teaagent

import (
	"errors"
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
	source := this.config.Source()
	if source == nil {
		errMsg := "item " + this.config.Name + " source '" + this.config.SourceCode + "' does not exist, please update this agent to latest version"
		PushEvent(NewItemEvent(runningAgent.Id, this.appId, this.config.Id, "", errors.New(errMsg)))
		return "", errors.New(errMsg)
	}
	return source.Execute(nil)
}

// 定时运行
func (this *Item) Schedule() {
	if this.lastTimer != nil {
		this.lastTimer.Stop()
	}
	this.lastTimer = timers.Every(this.config.IntervalDuration(), func(ticker *time.Ticker) {
		source := this.config.Source()
		if source == nil {
			errMsg := "item " + this.config.Name + " source '" + this.config.SourceCode + "' does not exist, please update this agent to latest version"
			logs.Println(errMsg)
			PushEvent(NewItemEvent(runningAgent.Id, this.appId, this.config.Id, "", errors.New(errMsg)))
			return
		}

		value, err := source.Execute(nil)
		if err != nil {
			logs.Println("error:" + err.Error())
		}

		// 执行动作
		for _, threshold := range this.config.Thresholds {
			if len(threshold.Actions) == 0 {
				continue
			}
			if threshold.Test(value) {
				logs.Println("run " + this.config.Name + " [" + threshold.Param + " " + threshold.Operator + " " + threshold.Value + "] actions")
				err1 := threshold.RunActions(map[string]string{})
				if err1 != nil {
					logs.Println("error:" + err1.Error())

					if err == nil {
						err = err1
					}
				}
			}
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
