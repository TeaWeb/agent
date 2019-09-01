package teaagent

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/TeaWeb/agent/teaconst"
	"github.com/TeaWeb/agent/teautils"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"io/ioutil"
	"net/http"
	"runtime"
	"time"
)

// 向Master同步事件
var db *teautils.FileBuffer = nil

func pushEvents() {
	db = teautils.NewFileBuffer(Tea.Root + "/logs/agent")

	// 读取本地数据库日志并发送到Master
	go func() {
		for db.Next() {
			db.Read(func(value []byte) {
				// Push到Master服务器
				req, err := http.NewRequest(http.MethodPost, connectConfig.Master+"/api/agent/push", bytes.NewReader(value))
				if err != nil {
					logs.Println("error:", err.Error())
				} else {
					err = func() error {
						req.Header.Set("Content-Type", "application/json")
						req.Header.Set("User-Agent", "TeaWeb Agent")
						req.Header.Set("Tea-Agent-Id", connectConfig.Id)
						req.Header.Set("Tea-Agent-Key", connectConfig.Key)
						req.Header.Set("Tea-Agent-Version", teaconst.AgentVersion)
						req.Header.Set("Tea-Agent-Os", runtime.GOOS)
						req.Header.Set("Tea-Agent-Arch", runtime.GOARCH)
						resp, err := HTTPClient.Do(req)

						if err != nil {
							logs.Println("error:", err.Error())
							return err
						}
						defer resp.Body.Close()
						if resp.StatusCode != 200 {
							return errors.New("") // 保持空字符串，方便其他地方识别错误
						}

						respBody, err := ioutil.ReadAll(resp.Body)
						if err != nil {
							logs.Println("error:", err.Error())
							return err
						}

						respJSON := maps.Map{}
						err = json.Unmarshal(respBody, &respJSON)
						if err != nil {
							logs.Println("error:", err.Error())
							return err
						}

						if respJSON.GetInt("code") != 200 {
							logs.Println("[/api/agent/push]error response from master:", string(respBody))
							time.Sleep(60 * time.Second)
							return nil
						}
						return nil
					}()
					if err != nil {
						time.Sleep(60 * time.Second)
					}
				}
			})
			time.Sleep(1 * time.Second)
		}
	}()

	// 读取日志并写入到本地数据库
	logId := time.Now().UnixNano()
	for {
		event := <-eventQueue

		if runningAgent.Id != "local" {
			// 进程事件
			if event, found := event.(*ProcessEvent); found {
				if event.EventType == ProcessEventStdout || event.EventType == ProcessEventStderr {
					logs.Println("[" + findTaskName(event.TaskId) + "]" + string(event.Data))
				} else if event.EventType == ProcessEventStart {
					logs.Println("[" + findTaskName(event.TaskId) + "]start")
				} else if event.EventType == ProcessEventStop {
					logs.Println("[" + findTaskName(event.TaskId) + "]stop")
				}
			}
		}

		jsonData, err := event.AsJSON()
		if err != nil {
			logs.Println("error:", err.Error())
			continue
		}

		if db != nil {
			logId++
			err = db.Write(append(jsonData, '\n'))
			if err != nil {
				logs.Println("[ERROR]" + err.Error())
			}
		}
	}
}
