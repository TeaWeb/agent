package teaagent

import (
	"fmt"
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/TeaWeb/code/teaconfigs/widgets"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/timers"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"log"
	"math"
	"runtime"
	"strings"
	"time"
)

// 系统应用
type SystemAppsProbe struct {
	apps []*agents.AppConfig

	systemApp *agents.AppConfig

	cpuTicker     *time.Ticker
	memoryTicker  *time.Ticker
	loadTicker    *time.Ticker
	diskTicker    *time.Ticker
	networkTicker *time.Ticker
	clockTicker   *time.Ticker
}

func NewSystemAppsProbe() *SystemAppsProbe {
	systemApp := agents.NewAppConfig()
	systemApp.Name = "系统"
	systemApp.IsSystem = true
	systemApp.Id = "system"

	return &SystemAppsProbe{
		apps:      []*agents.AppConfig{systemApp},
		systemApp: systemApp,
	}
}

func (this *SystemAppsProbe) Apps() []*agents.AppConfig {
	return this.apps
}

func (this *SystemAppsProbe) AddApp(app *agents.AppConfig) {
	this.apps = append(this.apps, app)
}

func (this *SystemAppsProbe) Run() {
	log.Println("run system apps probe")
	this.runCPU()
	this.runLoad()
	this.runMemory()
	this.runClock()
	this.runNetwork()
	this.runDisk()
	this.runApps()
}

func (this *SystemAppsProbe) runCPU() {
	// item
	item := agents.NewItem()
	item.Id = "cpu.usage"
	item.Name = "CPU使用量（%）"
	item.Interval = "60s"
	this.systemApp.AddItem(item)

	// 阈值
	threshold1 := agents.NewThreshold()
	threshold1.Param = "${0}"
	threshold1.Value = "80"
	threshold1.NoticeLevel = agents.NoticeLevelWarning
	threshold1.Operator = agents.ThresholdOperatorGte
	item.AddThreshold(threshold1)

	// chart
	chart := widgets.NewChart()
	chart.Id = "cpu.chart1"
	chart.Name = "CPU使用量（%）"
	chart.Columns = 2
	chart.Type = "javascript"
	chart.Options = maps.Map{
		"Code": `
var chart = new charts.LineChart();
chart.max = 100;

var query = new values.Query();
query.limit(30)
var ones = query.desc().cache(60).findAll();
ones.reverse();

var lines = [];

{
	var line = new charts.Line();
	line.color = colors.ARRAY[0];
	line.isFilled = true;
	line.values = [];
	lines.push(line);
}

ones.$each(function (k, v) {
	lines[0].values.push(v.value);
	
	var minute = v.timeFormat.minute.substring(8);
	chart.labels.push(minute.substr(0, 2) + ":" + minute.substr(2, 2));
});

chart.addLines(lines);
chart.render();
`,
	}
	item.AddChart(chart)

	// 数值
	if this.cpuTicker != nil {
		this.cpuTicker.Stop()
	}
	this.cpuTicker = this.every(60*time.Second, func() {
		stat, err := cpu.Percent(1*time.Second, false)
		if err != nil || len(stat) == 0 {
			return
		}
		PushEvent(NewItemEvent(runningAgent.Id, this.systemApp.Id, item.Id, math.Round(stat[0])))
	})
}

func (this *SystemAppsProbe) runMemory() {
	//item
	item := agents.NewItem()
	item.Id = "memory.usage"
	item.Name = "内存使用量"
	item.Interval = "60s"
	this.systemApp.AddItem(item)

	// 阈值
	{
		threshold1 := agents.NewThreshold()
		threshold1.Param = "${virtualPercent}"
		threshold1.Value = "80"
		threshold1.NoticeLevel = agents.NoticeLevelWarning
		threshold1.Operator = agents.ThresholdOperatorGte
		item.AddThreshold(threshold1)
	}

	// chart
	{
		chart := widgets.NewChart()
		chart.Id = "memory.usage.chart1"
		chart.Name = "内存使用量（%）"
		chart.Columns = 2
		chart.Type = "javascript"
		chart.Options = maps.Map{
			"Code": `
var chart = new charts.LineChart();

var query = new values.Query();
query.limit(30)
var ones = query.desc().cache(60).findAll();
ones.reverse();

var lines = [];

{
	var line = new charts.Line();
	line.color = colors.ARRAY[0];
	line.isFilled = true;
	line.values = [];
	lines.push(line);
}

ones.$each(function (k, v) {
	lines[0].values.push(v.value.virtualPercent);

	var minute = v.timeFormat.minute.substring(8);
	chart.labels.push(minute.substr(0, 2) + ":" + minute.substr(2, 2));
});

chart.addLines(lines);
chart.max = 100;
chart.render();
`,
		}
		item.AddChart(chart)
	}

	{
		chart := widgets.NewChart()
		chart.Id = "memory.usage.chart2"
		chart.Name = "当前内存使用量"
		chart.Columns = 1
		chart.Type = "javascript"
		chart.Options = maps.Map{
			"Code": `
var chart = new charts.StackBarChart();

var latest = new values.Query().latest(1);
var hasWarning = false;
if (latest.length > 0) {
	hasWarning = (latest[0].value.swapPercent > 50) || (latest[0].value.virtualPercent > 80);
	chart.values = [ 
		[latest[0].value.swapUsed, latest[0].value.swapTotal - latest[0].value.swapUsed],
		[latest[0].value.virtualUsed, latest[0].value.virtualTotal - latest[0].value.virtualUsed]
	];
	chart.labels = [ "虚拟内存（" +  (Math.round(latest[0].value.swapUsed * 10) / 10) + "G/" + Math.round(latest[0].value.swapTotal) + "G"  + "）", "物理内存（" + (Math.round(latest[0].value.virtualUsed * 10) / 10)+ "G/" + Math.round(latest[0].value.virtualTotal)  + "G"  + "）"];
} else {
	chart.values = [ [0, 0], [0, 0] ];
	chart.labels = [ "虚拟内存", "物理内存" ];
}
if (hasWarning) {
	chart.colors = [ colors.RED, colors.GREEN ];
} else {
	chart.colors = [ colors.BROWN, colors.GREEN ];
}
chart.render();
`,
		}
		item.AddChart(chart)
	}

	// 数值
	if this.memoryTicker != nil {
		this.memoryTicker.Stop()
	}
	this.memoryTicker = this.every(60*time.Second, func() {
		stat, err := mem.VirtualMemory()
		if err != nil {
			return
		}

		swap, err := mem.SwapMemory()
		if err != nil {
			return
		}

		PushEvent(NewItemEvent(runningAgent.Id, this.systemApp.Id, item.Id, maps.Map{
			"virtualUsed":    float64(stat.Used) / 1024 / 1024 / 1024,
			"virtualPercent": stat.UsedPercent,
			"virtualTotal":   float64(stat.Total) / 1024 / 1024 / 1024,
			"swapUsed":       float64(swap.Used) / 1024 / 1024 / 1024,
			"swapPercent":    swap.UsedPercent,
			"swapTotal":      float64(swap.Total) / 1024 / 1024 / 1024,
		}))
	})
}

func (this *SystemAppsProbe) runLoad() {
	if runtime.GOOS == "windows" {
		return
	}

	// item
	item := agents.NewItem()
	item.Id = "cpu.load"
	item.Name = "负载（Load）"
	item.Interval = "60s"
	this.systemApp.AddItem(item)

	// 阈值
	{
		threshold1 := agents.NewThreshold()
		threshold1.Param = "${load5}"
		threshold1.Value = "10"
		threshold1.NoticeLevel = agents.NoticeLevelWarning
		threshold1.Operator = agents.ThresholdOperatorGte
		item.AddThreshold(threshold1)
	}

	{
		threshold2 := agents.NewThreshold()
		threshold2.Param = "${load5}"
		threshold2.Value = "20"
		threshold2.NoticeLevel = agents.NoticeLevelWarning
		threshold2.Operator = agents.ThresholdOperatorGte
		item.AddThreshold(threshold2)
	}

	// chart
	chart := widgets.NewChart()
	chart.Id = "cpu.load.chart1"
	chart.Name = "负载（Load）"
	chart.Columns = 2
	chart.Type = "javascript"
	chart.Options = maps.Map{
		"Code": `
var chart = new charts.LineChart();

var query = new values.Query();
query.limit(30)
var ones = query.desc().cache(60).findAll();
ones.reverse();

var lines = [];

{
	var line = new charts.Line();
	line.name = "1分钟";
	line.color = colors.ARRAY[0];
	line.isFilled = true;
	line.values = [];
	lines.push(line);
}

{
	var line = new charts.Line();
	line.name = "5分钟";
	line.color = colors.BROWN;
	line.isFilled = false;
	line.values = [];
	lines.push(line);
}

{
	var line = new charts.Line();
	line.name = "15分钟";
	line.color = colors.RED;
	line.isFilled = false;
	line.values = [];
	lines.push(line);
}

var maxValue = 10;

ones.$each(function (k, v) {
	lines[0].values.push(v.value.load1);
	lines[1].values.push(v.value.load5);
	lines[2].values.push(v.value.load15);

	if (v.value.load1 > maxValue) {
		maxValue = Math.ceil(v.value.load1 / 10) * 10;
	}
	if (v.value.load5 > maxValue) {
		maxValue = Math.ceil(v.value.load5 / 10) * 10;
	}
	if (v.value.load15 > maxValue) {
		maxValue = Math.ceil(v.value.load15 / 10) * 10;
	}
	
	var minute = v.timeFormat.minute.substring(8);
	chart.labels.push(minute.substr(0, 2) + ":" + minute.substr(2, 2));
});

chart.addLines(lines);
chart.max = maxValue;
chart.render();
`,
	}
	item.AddChart(chart)

	// 数值
	if this.loadTicker != nil {
		this.loadTicker.Stop()
	}
	this.loadTicker = this.every(60*time.Second, func() {
		stat, err := load.Avg()
		if err != nil || stat == nil {
			return
		}
		PushEvent(NewItemEvent(runningAgent.Id, this.systemApp.Id, item.Id, maps.Map{
			"load1":  stat.Load1,
			"load5":  stat.Load5,
			"load15": stat.Load15,
		}))
	})
}

func (this *SystemAppsProbe) runNetwork() {
	duration := 60

	// item
	item := agents.NewItem()
	item.Id = "network.usage"
	item.Name = "网络相关"
	item.Interval = fmt.Sprintf("%ds", duration)
	this.systemApp.AddItem(item)

	// 图表
	{
		chart := widgets.NewChart()
		chart.Id = "network.usage.received"
		chart.Name = "出口带宽（M/s）"
		chart.Columns = 2
		chart.Type = "javascript"
		chart.Options = maps.Map{
			"Code": `
var chart = new charts.LineChart();

var line = new charts.Line();
line.isFilled = true;

var ones = new values.Query().cache(60).latest(60);
ones.reverse();
ones.$each(function (k, v) {
	line.values.push(Math.round(v.value.avgSent / 1024 / 1024 * 100) / 100);
	
	var minute = v.timeFormat.minute.substring(8);
	chart.labels.push(minute.substr(0, 2) + ":" + minute.substr(2, 2));
});
var maxValue = line.values.$max();
if (maxValue < 1) {
	chart.max = 1;
} else if (maxValue < 5) {
	chart.max = 5;
} else if (maxValue < 10) {
	chart.max = 10;
}

chart.addLine(line);
chart.render();
`,
		}
		item.AddChart(chart)
	}

	{
		chart := widgets.NewChart()
		chart.Id = "network.usage.sent"
		chart.Name = "入口带宽（M/s）"
		chart.Columns = 2
		chart.Type = "javascript"
		chart.Options = maps.Map{
			"Code": `
var chart = new charts.LineChart();

var line = new charts.Line();
line.isFilled = true;

var ones = new values.Query().cache(60).latest(60);
ones.reverse();
ones.$each(function (k, v) {
	line.values.push(Math.round(v.value.avgReceived / 1024 / 1024 * 100) / 100);
	
	var minute = v.timeFormat.minute.substring(8);
	chart.labels.push(minute.substr(0, 2) + ":" + minute.substr(2, 2));
});
var maxValue = line.values.$max();
if (maxValue < 1) {
	chart.max = 1;
} else if (maxValue < 5) {
	chart.max = 5;
} else if (maxValue < 10) {
	chart.max = 10;
}

chart.addLine(line);
chart.render();
`,
		}
		item.AddChart(chart)
	}

	// 数值
	if this.networkTicker != nil {
		this.networkTicker.Stop()
	}

	countSent := uint64(0)
	countReceived := uint64(0)

	lastTotalSent := uint64(0)
	lastTotalReceived := uint64(0)

	this.networkTicker = this.every(time.Duration(duration)*time.Second, func() {
		stats, err := net.IOCounters(false)
		if err != nil {
			return
		}
		totalSent := uint64(0)
		totalReceived := uint64(0)
		for _, stat := range stats {
			totalSent += stat.BytesSent
			totalReceived += stat.BytesRecv
		}

		if lastTotalSent == 0 {
			lastTotalSent = totalSent
			lastTotalReceived = totalReceived
			return
		}

		if totalSent > lastTotalSent {
			countSent = totalSent - lastTotalSent
		} else {
			countSent = 0
		}
		if totalReceived > lastTotalReceived {
			countReceived = totalReceived - lastTotalReceived
		} else {
			countReceived = 0
		}

		PushEvent(NewItemEvent(runningAgent.Id, this.systemApp.Id, item.Id, maps.Map{
			"avgSent":     countSent / uint64(duration),
			"avgReceived": countReceived / uint64(duration),
		}))

		lastTotalSent = totalSent
		lastTotalReceived = totalReceived
	})
}

func (this *SystemAppsProbe) runDisk() {
	// item
	item := agents.NewItem()
	item.Id = "disk.usage"
	item.Name = "文件系统"
	item.Interval = "120s"
	this.systemApp.AddItem(item)

	// 图表
	{
		chart := widgets.NewChart()
		chart.Id = "disk.usage.chart1"
		chart.Name = "文件系统"
		chart.Columns = 2
		chart.Type = "javascript"
		chart.Options = maps.Map{
			"Code": `
var chart = new charts.StackBarChart();
chart.values = [];
chart.labels = [];

var latest = new values.Query().cache(60).latest(1);
if (latest.length > 0) {
	var partitions = latest[0].value;
	partitions.$each(function (k, v) {
		chart.values.push([v.used, v.total - v.used]);
		chart.labels.push(v.name + "（" + (Math.round(v.used / 1024 / 1024 / 1024 * 100) / 100)+ "G/" + (Math.round(v.total / 1024 / 1024 / 1024 * 100) / 100) +"G）");
	});
}

chart.colors = [ colors.BROWN, colors.GREEN ];
chart.render();
`,
		}
		item.AddChart(chart)
	}

	// 数值
	if this.diskTicker != nil {
		this.diskTicker.Stop()
	}

	this.diskTicker = this.every(120*time.Second, func() {
		partitions, err := disk.Partitions(false)
		if err != nil {
			logs.Error(err)
			return
		}
		lists.Sort(partitions, func(i int, j int) bool {
			p1 := partitions[i]
			p2 := partitions[j]
			return p1.Mountpoint > p2.Mountpoint
		})

		result := []maps.Map{}
		for _, partition := range partitions {
			if runtime.GOOS != "windows" && !strings.Contains(partition.Device, "/") && !strings.Contains(partition.Device, "\\") {
				continue
			}

			usage, err := disk.Usage(partition.Mountpoint)
			if err != nil {
				continue
			}
			result = append(result, maps.Map{
				"name":    partition.Mountpoint,
				"used":    usage.Used,
				"total":   usage.Total,
				"percent": usage.UsedPercent,
			})
		}

		PushEvent(NewItemEvent(runningAgent.Id, this.systemApp.Id, item.Id, result))
	})
}

func (this *SystemAppsProbe) runClock() {
	// item
	item := agents.NewItem()
	item.Id = "clock"
	item.Name = "时钟"
	item.Interval = "60s"
	this.systemApp.AddItem(item)

	// 时钟
	{
		chart := widgets.NewChart()
		chart.Id = "clock"
		chart.Name = "时钟"
		chart.Columns = 1
		chart.Type = "javascript"
		chart.Options = maps.Map{
			"Code": `
var chart = new charts.Clock();
var latest = new values.Query().latest(1);
if (latest.length > 0) {
	chart.timestamp = latest[0].value.timestamp;
	console.log(chart.timestamp);
}
chart.render();
`,
		}
		item.AddChart(chart)
	}

	// 数值
	if this.clockTicker != nil {
		this.clockTicker.Stop()
	}
	this.clockTicker = this.every(60*time.Second, func() {
		PushEvent(NewItemEvent(runningAgent.Id, this.systemApp.Id, item.Id, maps.Map{
			"timestamp": time.Now().Unix(),
		}))
	})
}

func (this *SystemAppsProbe) runApps() {

}

func (this *SystemAppsProbe) every(duration time.Duration, f func()) *time.Ticker {
	go f()
	return timers.Every(duration, func(ticker *time.Ticker) {
		f()
	})
}
