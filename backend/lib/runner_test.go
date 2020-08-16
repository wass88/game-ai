package lib

import (
	"testing"
)

func TestRunnerCmd(t *testing.T) {
	conf := RunnerConf{
		PlayoutCmd: "playout",
		PlayerCmd:  "player",
		Api:        "http://xxx",
	}
	task := PlayoutTask{
		Game:      "reversi",
		PlayoutID: PlayoutID{1, nil},
		Token:     "TOKEN",
		Players:   []int64{1, 2},
	}
	cmd := conf.Cmd(task)
	exp := []string{"playout", "reversi", "--send", "http://xxx!1!TOKEN",
		"player 1", "player 2"}
	for i, e := range exp {
		if cmd.Args[i] != e {
			t.Fatalf("%v[0] = %s is not %s", cmd.Args, cmd.Args[i], e)
		}
	}
}
