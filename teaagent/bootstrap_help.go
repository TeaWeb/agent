package teaagent

import (
	"github.com/TeaWeb/agent/teaconst"
	"github.com/TeaWeb/agent/teautils"
)

// 打印帮助
func printHelp() {
	teautils.NewCommandHelp().
		Product(teaconst.AgentProductName).
		Version(teaconst.AgentVersion).
		Usage("./bin/"+teaconst.AgentProcessName+" [options]").
		Option("-h", "show this help").
		Option("-v|version", "show agent version").
		Option("start", "start agent in background").
		Option("stop", "stop running agent").
		Option("restart", "restart the agent").
		Option("status", "lookup agent status").
		Option("run [TASK ID]", "run task").
		Option("run [ITEM ID]", "run app item").
		Option("init -master=[MASTER SERVER] -group=[GROUP KEY]", "register agent to master server and specified group").
		Append("To run agent in foreground\n   bin/teaweb-agent").
		Print()
}
