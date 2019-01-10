package teaconfigs

import (
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
)

type AgentConfig struct {
	Master string `yaml:"master" json:"master"`
	Id     string `yaml:"id" json:"id"`
	Key    string `yaml:"key" json:"key"`
}

func SharedAgentConfig() (*AgentConfig, error) {
	reader, err := files.NewReader(Tea.ConfigFile("agent.conf"))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	a := &AgentConfig{}
	err = reader.ReadYAML(a)
	if err != nil {
		return nil, err
	}
	return a, nil
}
