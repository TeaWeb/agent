package teaagent

import (
	"github.com/TeaWeb/code/teaconst"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/utils/string"
	"net/http"
)

// Agent提供的Web Server，可以用来读取状态、进行控制等
type Server struct {
	Addr string
}

// 获取新服务对象
func NewServer() *Server {
	return &Server{
		Addr: "127.0.0.1:7778",
	}
}

// 启动
func (this *Server) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/status", this.HandleStatus)
	err := http.ListenAndServe(this.Addr, mux)
	return err
}

// 处理 /status
func (this *Server) HandleStatus(resp http.ResponseWriter, req *http.Request) {
	resp.Write([]byte(stringutil.JSONEncode(maps.Map{
		"version": teaconst.TeaVersion,
	})))
}
