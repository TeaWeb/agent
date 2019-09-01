package teaagent

import (
	"fmt"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/types"
	"os"
	"runtime"
	"syscall"
)

// lookup status
func onStatus() {
	pidString, err := files.NewFile(Tea.Root + Tea.DS + "logs" + Tea.DS + "pid").ReadAllString()
	if err != nil {
		fmt.Println("Agent not started yet")
		return
	}

	pid := types.Int(pidString)
	proc, err := os.FindProcess(pid)
	if err != nil {
		fmt.Println("Agent not started yet")
		return
	}
	if proc == nil {
		fmt.Println("Agent not started yet")
		return
	}

	if runtime.GOOS == "windows" {
		fmt.Println("Agent is running, pid:" + pidString)
		return
	}

	err = proc.Signal(syscall.Signal(0))
	if err != nil {
		fmt.Println("Agent not started yet")
		return
	}
	fmt.Println("Agent is running, pid:" + pidString)
}
