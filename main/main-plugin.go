package main

import (
	"github.com/TeaWeb/agent/teaagent"
	"github.com/TeaWeb/agent/teaconst"
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
	p.Version = "v" + teaconst.AgentVersion
	p.Date = "2019-11-01"
	p.Site = "https://github.com/TeaWeb/agent"
	p.Description = "本地Agent插件"
	p.OnStart(func() {
		timers.Delay(2*time.Second, func(timer *time.Timer) {
			teaagent.Start()
		})
	})
	loader.Start(p)
}
