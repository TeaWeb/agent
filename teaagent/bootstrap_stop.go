package teaagent

import (
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/types"
	"os"
)

// 停止
func onStop() {
	pidFile := files.NewFile(Tea.Root + "/logs/pid")
	pid, err := pidFile.ReadAllString()
	if err != nil {
		logs.Println("error:", err.Error())
	} else {
		process, err := os.FindProcess(types.Int(pid))
		if err != nil {
			logs.Println("error:", err.Error())
		} else {
			process.Kill()
			logs.Println("stopped pid", pid)
			pidFile.Delete()
		}
	}
}
