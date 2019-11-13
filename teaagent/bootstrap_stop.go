package teaagent

import (
	"fmt"
	"github.com/TeaWeb/agent/teaconst"
	"github.com/TeaWeb/code/teautils"
	"github.com/iwind/TeaGo/Tea"
)

// 停止
func onStop() {
	pidFile := Tea.Root + "/logs/pid"
	proc := teautils.CheckPid(pidFile)
	if proc == nil {
		fmt.Println(teaconst.AgentProductName + " agent is not running")
		return
	}

	_ = proc.Kill()
	fmt.Println(teaconst.AgentProductName+" stopped pid:", proc.Pid)

	_ = teautils.DeletePid(pidFile)
}
