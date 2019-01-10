package teaagent

import (
	"github.com/TeaWeb/code/teaconfigs/agents"
	"testing"
)

func TestApp_Run(t *testing.T) {
	config := agents.NewAppConfig()

	app := NewApp(config)
	t.Log(app)
}
