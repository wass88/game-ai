package lib

import (
	"fmt"
	"os/exec"
)

type RunnerConf struct {
	PlayoutCmd string
	PlayerCmd  string
	API        string
}

func (r *RunnerConf) Cmd(t PlayoutTask) *exec.Cmd {
	send := fmt.Sprintf("%s!%d!%s", r.API, t.PlayoutID.ID, t.Token)
	args := []string{t.Game, "--send=" + send}
	for _, p := range t.Players {
		args = append(args, fmt.Sprintf("%s %d", r.PlayerCmd, p))
	}
	cmd := exec.Command(r.PlayoutCmd, args...)
	return cmd
}

func (r *RunnerConf) Run(t PlayoutTask) error {
	return nil
}
