package teaagent

import (
	"fmt"
	"github.com/TeaWeb/agent/teautils"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/logs"
	"os/exec"
	"time"
)

// 启动
func onStart() {
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
		fmt.Println("started ok")
	}
}
