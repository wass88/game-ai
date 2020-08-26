package server

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func (p *PlayerCmd) Cmd(c AIRunnerConf) string {
	//Warning!!!: Check Cmd Injection
	cmd := []string{c.Cmd, "-api", c.API, "-dir", c.Dir, "-cpu", strconv.Itoa(c.CPU), "-mem", strconv.Itoa(c.MemMB), "-github", p.Github, "-branch", p.Branch, "-commit", p.Commit, "run"}
	return strings.Join(cmd, " ")
}

func (t PlayoutTask) Cmd(r *Config) *exec.Cmd {
	send := fmt.Sprintf("%s@%d@%s", r.Playout.API, t.PlayoutID.ID, t.Token)
	args := []string{"--send=" + send, t.Game}
	for _, p := range t.Players {
		args = append(args, fmt.Sprintf("%s", p.Cmd(r.AIRunner)))
	}
	cmd := exec.Command(r.Playout.Cmd, args...)
	return cmd
}

func (t PlayoutTask) Run(r *Config) error {
	return nil
}
