package main

import (
	"github.com/TeaWeb/agent/teaagent"
	"github.com/TeaWeb/plugin/loader"
	"github.com/TeaWeb/plugin/plugins"
	"github.com/iwind/TeaGo/timers"
	"time"
)

func main() {
	p := plugins.NewPlugin()
	p.Name = "Agent"
	p.Code = "agent.teaweb"
	p.Developer = "TeaWeb"
	p.Version = "v0.0.3"
	p.Date = "2019-08-28"
	p.Site = "https://github.com/TeaWeb/agent"
	p.Description = "本地Agent插件"
	p.OnStart(func() {
		timers.Delay(5*time.Second, func(timer *time.Timer) {
			teaagent.Start()
		})
	})
	loader.Start(p)
}
