package lib

import (
	"os/exec"
	"testing"
)

func TestPlayout(t *testing.T) {
	cmd0 := exec.Command("/Users/admin/Documents/reversi-random/target/release/reversi_random")
	cmd1 := exec.Command("/Users/admin/Documents/reversi-random/target/release/reversi_random")
	r, err := Playout(cmd0, cmd1)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+#v", r)
	if r.Exception != "" {
		t.Fatal(r.Exception)
	}
	if r.Result[0].Exception != "" {
		t.Fatal(r.Result[0].Exception)
	}
	if r.Result[1].Exception != "" {
		t.Fatal(r.Result[1].Exception)
	}
	if r.Result[0].Result != -r.Result[1].Result {
		t.Fatalf("%d != - %d", r.Result[0].Result, r.Result[1].Result)
	}
}

func TestReversi(t *testing.T) {
	r := NewReversi()
	p := r.playable()
	if len(p) != 4 {
		t.Fatalf("Playable is 4 : %+#v", p)
	}
}
