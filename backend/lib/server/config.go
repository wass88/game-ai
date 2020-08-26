package server

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	AIRunner AIRunnerConf `json:"ai_runner"`
	Playout  PlayoutConf  `json:"playout"`
	Session  SessionConf  `json:"session"`
	DBName   string       `json:"db_name"`
}

type AIRunnerConf struct {
	Cmd   string `json:"cmd"`
	API   string `json:"api"`
	Dir   string `json:"dir"`
	CPU   int    `json:"cpu"`
	MemMB int    `json:"mem_mb"`
}
type PlayoutConf struct {
	Cmd string `json:"cmd"`
	API string `json:"api"`
}

type SessionConf struct {
	SessionSecret      string `json:"session_secret"`
	GithubClientID     string `json:"github_client_id"`
	GithubClientSecret string `json:"github_client_secret"`
	GithubRedirectURL  string `json:"github_redirect_url"`
}

func LoadConfig(filename string) (*Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	var res *Config
	err = json.Unmarshal(bytes, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
