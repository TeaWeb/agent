package teaagent

import (
	"bytes"
	"github.com/TeaWeb/code/teaconfigs/agents"
	"os"
	"os/exec"
)

type Process struct {
	Env  []*agents.EnvVariable
	Cwd  string
	File string

	proc    *os.Process
	onStart func()
	onStop  func()
}

// 获取新进程
func NewProcess() *Process {
	return &Process{}
}

// on start
func (this *Process) OnStart(f func()) {
	this.onStart = f
}

// on stop
func (this *Process) OnStop(f func()) {
	this.onStop = f
}

// 立即运行
func (this *Process) Run() (stdout string, stderr string, err error) {
	// execute
	cmd := exec.Command(this.File)

	outputWriter := bytes.NewBuffer([]byte{})
	cmd.Stdout = outputWriter

	errorWriter := bytes.NewBuffer([]byte{})
	cmd.Stderr = errorWriter

	// cwd
	if len(this.Cwd) > 0 {
		cmd.Dir = this.Cwd
	}

	// env
	for _, envVar := range this.Env {
		cmd.Env = append(cmd.Env, envVar.Name+"="+envVar.Value)
	}

	err = cmd.Start()
	if this.onStart != nil {
		this.onStart()
	}
	this.proc = cmd.Process
	err = cmd.Wait()
	defer func() {
		stdout = outputWriter.String()
		stderr = errorWriter.String()
		if this.onStop != nil {
			this.onStop()
		}
	}()
	if err != nil {
		return
	}

	return
}

// Kill进程
func (this *Process) Kill() error {
	if this.proc == nil {
		return nil
	}
	return this.proc.Kill()
}
