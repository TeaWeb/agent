package teaagent

import (
	"fmt"
	"github.com/TeaWeb/agent/teaconst"
	"github.com/TeaWeb/agent/teautils"
	teautils2 "github.com/TeaWeb/code/teautils"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/logs"
	"os/exec"
	"time"
)

// 启动
func onStart() {
	// 检查是否已经在运行
	proc := teautils2.CheckPid(Tea.Root + "/logs/pid")
	if proc != nil {
		fmt.Println(teaconst.AgentProductName+" already started, pid:", proc.Pid)
		return
	}

	cmdFile := teautils.Executable()
	cmd := exec.Command(cmdFile, "background")
	cmd.Dir = Tea.Root
	err := cmd.Start()
	if err != nil {
		logs.Error(err)
		return
	}

	failed := false
	go func() {
		err = cmd.Wait()
		if err != nil {
			logs.Error(err)
		}

		failed = true
	}()

	time.Sleep(1 * time.Second)
	if failed {
		fmt.Println("error: process terminated, lookup 'logs/run.log' for more details")
	} else {
		fmt.Println(teaconst.AgentProductName+" started ok pid:", cmd.Process.Pid)
	}
}
