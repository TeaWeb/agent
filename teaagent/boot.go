package teaagent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TeaWeb/agent/teaconfigs"
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/go-yaml/yaml"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"io/ioutil"
	"net"
	"net/http"
	"reflect"
	"regexp"
	"sync"
	"time"
)

var connectConfig *teaconfigs.AgentConfig = nil
var runningAgent *agents.AgentConfig = nil
var runningTasks = map[string]*Task{}
var runningTasksLocker = sync.Mutex{}
var isBooting = true

// 启动
func Start() {
	logs.Println("starting ...")

	// 连接配置
	{
		config, err := teaconfigs.SharedAgentConfig()
		if err != nil {
			logs.Println("start failed:" + err.Error())
			return
		}
		connectConfig = config
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
	bootTasks()
	isBooting = false

	// 定时
	scheduleTasks()

	// 同步数据
	for {
		err := pullTasks()
		if err != nil {
			logs.Println("pull error:", err.Error())
			time.Sleep(5 * time.Second)
		}
	}
}

// 连接主服务器
func downloadConfig() error {
	master := connectConfig.Master
	if len(master) == 0 {
		return errors.New("'master' should not be empty")
	}
	req, err := http.NewRequest(http.MethodGet, master+"/api/agent", nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "TeaWeb Agent")
	req.Header.Set("Tea-Agent-Id", connectConfig.Id)
	req.Header.Set("Tea-Agent-Key", connectConfig.Key)
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("invalid status response from master '" + fmt.Sprintf("%d", resp.StatusCode) + "'")
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	respMap := maps.Map{}
	err = json.Unmarshal(data, &respMap)
	if err != nil {
		return err
	}

	if respMap == nil {
		return errors.New("response data should not be nil")
	}

	if respMap.GetInt("code") != 200 {
		return errors.New("response json code should be 200")
	}

	jsonData := respMap.Get("data")
	if jsonData == nil || reflect.TypeOf(jsonData).Kind() != reflect.Map {
		return errors.New("response json data should be a map")
	}

	dataMap := maps.NewMap(jsonData)
	config := dataMap.GetString("config")

	agent := &agents.AgentConfig{}
	err = yaml.Unmarshal([]byte(config), agent)
	if err != nil {
		return err
	}

	if len(agent.Id) == 0 {
		return errors.New("invalid agent id")
	}

	err = agent.Validate()
	if err != nil {
		return err
	}

	// 保存
	agentsDir := files.NewFile(Tea.ConfigFile("agents/"))
	if !agentsDir.IsDir() {
		err = agentsDir.Mkdir()
		if err != nil {
			return err
		}
	}
	agentFile := files.NewFile(Tea.ConfigFile("agents/agent." + agent.Id + ".conf"))
	err = agentFile.WriteString(config)
	if err != nil {
		return err
	}

	runningAgent = agent
	connectConfig.Id = agent.Id
	connectConfig.Key = agent.Key

	if !isBooting {
		scheduleTasks()
	}

	return nil
}

// 启动任务
func bootTasks() {
	logs.Println("booting ...")
	if !runningAgent.On {
		return
	}
	for _, app := range runningAgent.Apps {
		if !app.On {
			continue
		}
		for _, taskConfig := range app.Tasks {
			if !taskConfig.On {
				continue
			}
			task := NewTask(taskConfig)
			if task.ShouldBoot() {
				_, stdout, stderr, err := task.Run()
				if len(stdout) > 0 {
					logs.Println(stdout)
				}
				if len(stderr) > 0 {
					logs.Println(stderr)
				}
				if err != nil {
					logs.Println(err.Error())
				}
			}
		}
	}
}

// 定时任务
func scheduleTasks() error {
	// 生成脚本
	taskIds := []string{}

	for _, app := range runningAgent.Apps {
		for _, taskConfig := range app.Tasks {
			taskIds = append(taskIds, taskConfig.Id)

			// 是否正在运行
			runningTask, found := runningTasks[taskConfig.Id]
			isChanged := true
			if found {
				// 如果有修改，则需要重启
				if runningTask.config.Version != taskConfig.Version {
					logs.Println("stop schedule task", taskConfig.Id, taskConfig.Name)
					runningTask.Stop()

					if taskConfig.On && len(taskConfig.Schedule) > 0 {
						logs.Println("restart schedule task", taskConfig.Id, taskConfig.Name)
						runningTask.config = taskConfig
						runningTask.Schedule()
					}
				} else {
					isChanged = false
				}
			} else if taskConfig.On && len(taskConfig.Schedule) > 0 { // 新任务，则启动
				logs.Println("schedule task", taskConfig.Id, taskConfig.Name)
				task := NewTask(taskConfig)
				task.Schedule()

				runningTasksLocker.Lock()
				runningTasks[taskConfig.Id] = task
				runningTasksLocker.Unlock()
			}

			// 生成脚本
			if isChanged {
				_, err := taskConfig.Generate()
				if err != nil {
					return err
				}
			}
		}
	}

	// 停止运行
	for taskId, runningTask := range runningTasks {
		if !lists.Contains(taskIds, taskId) {
			err := runningTask.Stop()
			if err != nil {
				logs.Error(err)
			}
		}
	}

	// 删除不存在的任务脚本
	files.NewFile(Tea.ConfigFile("agents/")).Range(func(file *files.File) {
		filename := file.Name()
		if !regexp.MustCompile("^task\\.\\w+\\.script$").MatchString(filename) {
			return
		}
		taskId := filename[len("task:") : len(filename)-len(".script")]
		if !lists.Contains(taskIds, taskId) {
			err := file.Delete()
			if err != nil {
				logs.Error(err)
			}
		}
	})

	return nil
}

// 从主服务器同步数据
func pullTasks() error {
	logs.Println("pull tasks")

	master := connectConfig.Master
	if len(master) == 0 {
		return errors.New("'master' should not be empty")
	}
	req, err := http.NewRequest(http.MethodGet, master+"/api/agent/pull", nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "TeaWeb Agent")
	req.Header.Set("Tea-Agent-Id", connectConfig.Id)
	req.Header.Set("Tea-Agent-Key", connectConfig.Key)
	connectingFailed := false
	client := http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				// 握手配置
				conn, err := (&net.Dialer{
					Timeout:   5 * time.Second,
					KeepAlive: 0,
					DualStack: true,
				}).DialContext(ctx, network, addr)
				if err != nil {
					connectingFailed = true
				}
				return conn, err
			},
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		if connectingFailed {
			return err
		}

		// 如果是超时的则不提示，因为长连接依赖超时设置
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("invalid status response from master '" + fmt.Sprintf("%d", resp.StatusCode) + "'")
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	respMap := maps.Map{}
	err = json.Unmarshal(data, &respMap)
	if err != nil {
		return err
	}

	if respMap == nil {
		return errors.New("response data should not be nil")
	}

	if respMap.GetInt("code") != 200 {
		return errors.New("response json code should be 200")
	}

	jsonData := respMap.Get("data")
	if jsonData == nil || reflect.TypeOf(jsonData).Kind() != reflect.Map {
		return errors.New("response json data should be a map")
	}

	dataMap := maps.NewMap(jsonData)
	events := dataMap.Get("events")
	if events == nil || reflect.TypeOf(events).Kind() != reflect.Slice {
		return nil
	}

	eventsValue := reflect.ValueOf(events)
	count := eventsValue.Len()
	for i := 0; i < count; i ++ {
		event := eventsValue.Index(i).Interface()
		if event == nil || reflect.TypeOf(event).Kind() != reflect.Map {
			continue
		}
		eventMap := maps.NewMap(event)
		name := eventMap.GetString("name")
		switch name {
		case "UPDATE_AGENT":
			downloadConfig()
		case "ADD_APP":
			downloadConfig()
		case "UPDATE_APP":
			downloadConfig()
		case "REMOVE_APP":
			downloadConfig()
		case "ADD_TASK":
			downloadConfig()
		case "UPDATE_TASK":
			downloadConfig()
		case "REMOVE_TASK":
			downloadConfig()
		}
	}

	return nil
}
