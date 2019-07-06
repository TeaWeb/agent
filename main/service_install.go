package main

import (
	"github.com/TeaWeb/agent/teautils"
	"github.com/iwind/TeaGo/Tea"
	"log"
	"runtime"
)

// 安装服务
func main() {
	log.Println("installing ...")
	manager := teautils.NewServiceManager("TeaWeb Agent", "TeaWeb Agent Manager")

	var exePath = Tea.Root + Tea.DS + "bin" + Tea.DS + "teaweb-agent"
	if runtime.GOOS == "windows" {
		exePath += ".exe"
	}
	err := manager.Install(exePath, []string{"service"})
	if err != nil {
		log.Println("ERROR: " + err.Error())
		manager.PauseWindow()
		return
	}

	log.Println("install service successfully")

	// start
	log.Println("starting ...")
	err = manager.Start()
	if err != nil {
		log.Println("ERROR: " + err.Error())
	}

	log.Println("started successfully")
	log.Println("done.")

	manager.PauseWindow()
}
