package teaagent

import (
	"fmt"
	"github.com/TeaWeb/agent/teaconfigs"
	"github.com/TeaWeb/agent/teaconst"
	"github.com/TeaWeb/agent/teautils"
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

var connectConfig *teaconfigs.AgentConfig = nil
var runningAgent *agents.AgentConfig = nil
var runningTasks = map[string]*Task{} // task id => task
var runningTasksLocker = sync.Mutex{}
var runningItems = map[string]*Item{} // item id => task
var runningItemsLocker = sync.Mutex{}
var isBooting = true
var connectionIsBroken = false

// 启动
func Start() {
	// 当前ROOT
	if !Tea.IsTesting() {
		exePath := teautils.Executable()
		if strings.Contains(filepath.Base(exePath), "@") { // 是不是升级的文件
			Tea.UpdateRoot(filepath.Dir(filepath.Dir(filepath.Dir(exePath))))
		} else {
			Tea.UpdateRoot(filepath.Dir(filepath.Dir(exePath)))
		}
	}

	// 帮助
	if lists.ContainsAny(os.Args, "h", "-h", "help", "-help") {
		printHelp()
		return
	}

	// 版本号
	if lists.ContainsAny(os.Args, "version", "-v") {
		fmt.Println("v" + teaconst.AgentVersion)
		return
	}

	if len(os.Args) == 0 {
		writePid()
	}

	// 连接配置
	{
		config, err := teaconfigs.SharedAgentConfig()
		if err != nil {
			logs.Println("start failed:" + err.Error())
			return
		}
		connectConfig = config
	}

	// 检查新版本
	if shouldStartNewVersion() {
		return
	}

	// 启动
	if lists.ContainsString(os.Args, "start") {
		onStart()
		return
	}

	// 停止
	if lists.ContainsString(os.Args, "stop") {
		onStop()
		return
	}

	// 重启
	if lists.ContainsString(os.Args, "restart") {
		onStop()
		onStart()
		return
	}

	// 查看状态
	if lists.ContainsString(os.Args, "status") {
		onStatus()
		return
	}

	// 运行某个脚本
	if lists.ContainsAny(os.Args, "run") {
		runTaskOrItem()
		return
	}

	// 测试连接
	if lists.ContainsAny(os.Args, "test", "-t") {
		err := testConnection()
		if err != nil {
			logs.Println("error:", err.Error())
		} else {
			logs.Println("connection to master is ok")
		}
		return
	}

	// 日志
	if lists.ContainsAny(os.Args, "background", "-d") {
		writePid()

		logDir := files.NewFile(Tea.Root + "/logs")
		if !logDir.IsDir() {
			logDir.Mkdir()
		}

		fp, err := os.OpenFile(Tea.Root+"/logs/run.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			log.SetOutput(fp)
		} else {
			logs.Println(err)
		}
	}

	// Windows服务
	if lists.ContainsAny(os.Args, "service") && runtime.GOOS == "windows" {
		manager := teautils.NewServiceManager("TeaWeb Agent", "TeaWeb Agent Manager")
		manager.Run()
	}

	logs.Println("agent starting ...")

	// 启动监听端口
	if connectConfig.Id != "local" && runtime.GOOS != "windows" {
		go startListening()
	}

	// 下载配置
	{
		err := downloadConfig()
		if err != nil {
			logs.Println("start failed:" + err.Error())
			return
		}
	}

	// 启动
	logs.Println("agent boot tasks ...")
	bootTasks()
	isBooting = false

	// 定时
	logs.Println("agent schedule tasks ...")
	scheduleTasks()

	// 监控项数据
	logs.Println("agent schedule items ...")
	scheduleItems()

	// 检测Apps
	logs.Println("agent detect tasks ...")
	detectApps()

	// 检查更新
	checkNewVersion()

	// 推送日志
	go pushEvents()

	// 同步配置
	for {
		err := pullEvents()
		if err != nil {
			logs.Println("pull error:", err.Error())
			time.Sleep(30 * time.Second)
		}
	}
}

// 初始化连接
func initConnection() {
	detectApps()
}

// 启动监听端口欧
var statusServer *Server = nil

func startListening() {
	statusServer = NewServer()
	err := statusServer.Start()
	if err != nil {
		logs.Error(err)
	}
}

// 检测App
func detectApps() {
	// 暂时不做任何事情
}

// 查找任务
func findTaskName(taskId string) string {
	if runningAgent == nil {
		return ""
	}
	_, task := runningAgent.FindTask(taskId)
	if task == nil {
		return ""
	}
	return task.Name
}

func writePid() {
	// write pid
	pidFile := files.NewFile(Tea.Root + "/logs/pid")
	pidFile.WriteString(fmt.Sprintf("%d", os.Getpid()))
}
