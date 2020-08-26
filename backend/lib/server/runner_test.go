package server

import (
	"testing"
)

func TestPlayoutCmd(t *testing.T) {
	cmd := PlayerCmd{1, "w/e", "master", "cccccc"}
	ai := AIRunnerConf{
		Cmd:   "cmd",
		API:   "http://",
		Dir:   "dir",
		CPU:   100,
		MemMB: 256,
	}
	r := cmd.Cmd(ai)
	exp := "cmd -api http:// -dir dir -cpu 100 -mem 256 -github w/e -branch master -commit cccccc run"
	if r != exp {
		t.Fatalf("Expected %s, got %s", exp, r)
	}

}

func TestRunnerCmd(t *testing.T) {
	playout := PlayoutConf{
		Cmd: "playout",
		API: "http://xxx",
	}
	ai := AIRunnerConf{
		Cmd:   "cmd",
		API:   "http://",
		Dir:   "dir",
		CPU:   100,
		MemMB: 256,
	}
	conf := Config{Playout: playout, AIRunner: ai}
	task := PlayoutTask{
		Game:      "reversi",
		PlayoutID: PlayoutID{1, nil},
		Token:     "TOKEN",
		Players: []PlayerCmd{
			{1, "w/e", "master", "cccccc"},
			{1, "w/e", "master", "dddddd"},
		},
	}
	cmd := task.Cmd(&conf)
	exp1 := "cmd -api http:// -dir dir -cpu 100 -mem 256 -github w/e -branch master -commit cccccc run"
	exp2 := "cmd -api http:// -dir dir -cpu 100 -mem 256 -github w/e -branch master -commit dddddd run"
	exp := []string{"playout", "--send=http://xxx@1@TOKEN", "reversi",
		exp1, exp2}
	for i, e := range exp {
		if cmd.Args[i] != e {
			t.Fatalf("%v[0]\n%s != %s", cmd.Args, cmd.Args[i], e)
		}
	}
}
