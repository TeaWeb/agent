package teaagent

import "github.com/TeaWeb/code/teaconfigs/agents"

type App struct {
	config *agents.AppConfig
}

func NewApp(config *agents.AppConfig) *App {
	return &App{
		config: config,
	}
}
